package finder

import (
	"concert-manager/domain"
	"concert-manager/file"
	"concert-manager/log"
	"concert-manager/ranker"
	"fmt"
	"strings"
	"time"
)

type finder interface {
	FindAllEvents(string, string) ([]domain.EventDetails, error)
}

type eventRanker interface {
	Rank(domain.EventDetails) domain.RankInfo
}

type upcomingEventsData struct {
	Events     []domain.EventDetails `json:"events"`
	LastLoaded time.Time             `json:"last_loaded"`
}

type savedDataCache interface {
	GetArtists() []domain.Artist
	GetVenues() []domain.Venue
	GetSavedEvents() []domain.Event
}

var upcomingEventTTL, _ = time.ParseDuration("24h")

type EventCacheFile struct {
	Timestamp      time.Time                     `json:"timestamp"`
	Version        string                        `json:"version"`
	Location       Location                      `json:"location"`
	UpcomingEvents map[string]upcomingEventsData `json:"upcoming_events"`
}

const (
	eventCacheVersion = "1.0"
	eventCacheFile    = "events.json"
)

type Cache struct {
	Location       Location
	Finder         finder
	Ranker         eventRanker
	SavedDataCache savedDataCache
	MetadataFinder MetadataFinder
	upcomingEvents map[string]upcomingEventsData
}

const (
	defaultCity      = "Atlanta"
	defaultStateCode = "GA"
)

func NewUpcomingEventCache() *Cache {
	cache := Cache{}
	cache.Location = Location{City: defaultCity, StateCode: defaultStateCode}
	cache.upcomingEvents = map[string]upcomingEventsData{}
	return &cache
}

func (c *Cache) GetUpcomingEvents() []domain.EventDetails {
	return c.GetRecommendedEvents(ranker.NoMinRec)
}

func (c *Cache) GetRecommendedEvents(level ranker.RecLevel) []domain.EventDetails {
	key := c.Location.key()
	if d, ok := c.upcomingEvents[key]; !ok {
		c.doRefresh()
	} else if isExpired(d.LastLoaded, upcomingEventTTL) {
		go c.doRefresh()
	}

	threshold, _ := ranker.ToThreshold(level)
	var events []domain.EventDetails
	for _, event := range c.upcomingEvents[key].Events {
		if event.Ranks.Rank >= threshold {
			events = append(events, domain.CloneEventDetail(event))
		}
	}
	return events
}

func (c *Cache) doRefresh() {
	err := c.RefreshUpcomingEvents()
	if err != nil {
		log.Alert("Failed to refresh upcoming events", err)
	}
}

func (c *Cache) RefreshUpcomingEvents() error {
	loc := c.Location
	key := c.Location.key()
	log.Info("Refreshing upcoming events for", key)
	events, err := c.Finder.FindAllEvents(loc.City, loc.StateCode)
	if err != nil {
		if _, ok := c.upcomingEvents[key]; !ok {
			eventData := upcomingEventsData{Events: []domain.EventDetails{}, LastLoaded: time.Time{}}
			c.upcomingEvents[key] = eventData
		}
		return err
	}

	log.Debugf("Cache found %d upcoming events", len(events))
	for i, event := range events {
		events[i] = c.enrichSavedData(event)
	}
	events = c.MetadataFinder.PopulateMetadata(events)

	for i, event := range events {
		rank := c.Ranker.Rank(event)
		events[i].Ranks = &rank
	}

	log.Debugf("Finished upcoming event refresh, found %d events for key %s", len(events), key)
	eventData := upcomingEventsData{Events: events, LastLoaded: time.Now().Round(0)}
	c.upcomingEvents[key] = eventData
	c.saveEventsToFile()
	return nil
}

func (c *Cache) GetLocation() Location {
	return c.Location
}

func (c *Cache) ChangeLocation(city, stateCode string) {
	loc := Location{City: city, StateCode: stateCode}
	log.Debugf("Updating upcomingEventCache location from %s to %s", c.Location, loc)
	c.Location = loc
}

type Location struct {
	City      string
	StateCode string
}

func (c Location) key() string {
	return fmt.Sprintf("%s#%s", strings.ToLower(c.City), strings.ToLower(c.StateCode))
}

func (c Location) String() string {
	return fmt.Sprintf("%s, %s", c.City, c.StateCode)
}

func (c *Cache) Invalidate() {
	c.upcomingEvents = map[string]upcomingEventsData{}
}

func isExpired(lastLoad time.Time, ttl time.Duration) bool {
	elapsedTime := time.Since(lastLoad)
	log.Debugf("Upcoming lastLoaded: %v, now: %v, elapsed: %v, ttl: %v", lastLoad, time.Now(), elapsedTime, ttl)
	return elapsedTime > ttl
}

func (c *Cache) InitializeFromFile() error {
	filePath, err := file.GetCacheFilePath(eventCacheFile)
	if err != nil {
		return fmt.Errorf("failed to get event cache file path: %w", err)
	}

	if !file.FileExists(filePath) {
		log.Debug("Event cache file does not exist")
		c.doRefresh()
		return nil
	}

	log.Info("Loading upcoming events from cache file")
	err = c.loadEventsFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to load upcoming events from file: %v", err)
	}

	key := c.Location.key()
	if _, ok := c.upcomingEvents[key]; ok && file.IsFileStale(filePath, upcomingEventTTL) {
		log.Info("Event cache file is stale, starting background refresh")
		go c.doRefresh()
	} else if !ok {
		log.Debug("No events for current location, starting background refresh")
		go c.doRefresh()
	}

	return nil
}

func (c *Cache) initializeEmpty() {
	c.upcomingEvents = make(map[string]upcomingEventsData)
}

func (c *Cache) loadEventsFromFile(filePath string) error {
	var cacheFile EventCacheFile
	err := file.ReadJSONFile(filePath, &cacheFile)
	if err != nil {
		return err
	}

	if cacheFile.Version != eventCacheVersion {
		go c.doRefresh()
		return nil
	}

	c.upcomingEvents = cacheFile.UpcomingEvents
	if c.upcomingEvents == nil {
		c.upcomingEvents = make(map[string]upcomingEventsData)
	}

	log.Infof("Loaded upcoming events for %d locations from cache file", len(c.upcomingEvents))
	return nil
}

func (c *Cache) saveEventsToFile() {
	filePath, err := file.GetCacheFilePath(eventCacheFile)
	if err != nil {
		log.Errorf("Failed to get event cache file path for saving: %v", err)
		return
	}

	cacheFile := EventCacheFile{
		Timestamp:      time.Now().Round(0),
		Version:        eventCacheVersion,
		Location:       c.Location,
		UpcomingEvents: c.upcomingEvents,
	}

	err = file.WriteJSONFile(filePath, cacheFile)
	if err != nil {
		log.Errorf("Failed to save events to file: %v", err)
		return
	}

	totalEvents := 0
	for _, data := range c.upcomingEvents {
		totalEvents += len(data.Events)
	}
	log.Infof("Successfully saved %d upcoming events across %d locations to cache file", totalEvents, len(c.upcomingEvents))
}
