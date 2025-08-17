package finder

import (
	"concert-manager/domain"
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
	events     []domain.EventDetails
	lastLoaded time.Time
}

type savedDataCache interface {
	GetArtists() []domain.Artist
	GetVenues() []domain.Venue
	GetSavedEvents() []domain.Event
}

var upcomingEventTTL, _ = time.ParseDuration("24h")

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
	} else if isExpired(d.lastLoaded, upcomingEventTTL) {
		go c.doRefresh()
	}

	threshold, _ := ranker.ToThreshold(level)
	var events []domain.EventDetails
	for _, event := range c.upcomingEvents[key].events {
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
			eventData := upcomingEventsData{events: []domain.EventDetails{}, lastLoaded: time.Time{}}
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
	eventData := upcomingEventsData{events: events, lastLoaded: time.Now().Round(0)}
	c.upcomingEvents[key] = eventData
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
