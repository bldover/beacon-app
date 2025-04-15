package finder

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/util"
	"fmt"
	"slices"
	"strings"
	"time"
)

type finder interface {
	FindAllEvents(string, string) ([]data.EventDetails, error)
}

type ranker interface {
	Rank(data.EventDetails) data.RankInfo
}

type savedEventCache interface {
    GetArtists() []data.Artist
	GetVenues() []data.Venue
	GetSavedEvents() []data.Event
}

type upcomingEventsData struct {
	events     []data.EventDetails
	lastLoaded time.Time
}

var	upcomingEventTTL, _ = time.ParseDuration("24h")

type Cache struct {
	Location       Location
	Finder         finder
	Ranker         ranker
	upcomingEvents map[string]upcomingEventsData
	savedEventCache savedEventCache
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

func (c *Cache) GetUpcomingEvents() []data.EventDetails {
	return c.GetRecommendedEvents(NoMinRec)
}

func (c *Cache) GetRecommendedEvents(level RecLevel) []data.EventDetails {
	key := c.Location.key()
	if d, ok := c.upcomingEvents[key]; !ok {
		c.doRefresh()
	} else if isExpired(d.lastLoaded, upcomingEventTTL) {
		go c.doRefresh()
	}

	threshold, _ := ToThreshold(level)
	var events []data.EventDetails
	for _, event := range c.upcomingEvents[key].events {
		if event.Ranks.Rank >= threshold {
			events = append(events, util.CloneEventDetail(event))
		}
	}
	return events
}

func (c *Cache) doRefresh() {
    err := c.RefreshUpcomingEvents()
	if err != nil {
		log.Error("Failed to refresh upcoming events", err)
	}
}

func (c *Cache) RefreshUpcomingEvents() error {
	loc := c.Location
	key := c.Location.key()
	log.Info("Refreshing upcoming events for", key)
	events, err := c.Finder.FindAllEvents(loc.City, loc.StateCode)
	if err != nil {
		if _, ok := c.upcomingEvents[key]; !ok {
			eventData := upcomingEventsData{events: []data.EventDetails{}, lastLoaded: time.Time{}}
			c.upcomingEvents[key] = eventData
		}
		return err
	}

	log.Debugf("Cache found %d upcoming events. Starting rank population", len(events))
	for i, event := range events {
		rank := c.Ranker.Rank(event)
		events[i].Ranks = &rank
	}

	log.Debugf("Finished upcoming event refresh, found %d events for key %s", len(events), key)
	eventData := upcomingEventsData{events: events, lastLoaded: time.Now().Round(0)}
	c.upcomingEvents[key] = eventData
	return nil
}

func (c Cache) enrichSavedData(events []data.EventDetails) []data.EventDetails {
	savedArtists := c.savedEventCache.GetArtists()
	savedVenues := c.savedEventCache.GetVenues()
	savedEvents := c.savedEventCache.GetSavedEvents()

	enriched := make([]data.EventDetails, len(events))
	for _, event := range events {
		eventIdx := slices.IndexFunc(savedEvents, event.Event.EqualsFields)
		if eventIdx != -1 {
			event.Event = savedEvents[eventIdx]
			enriched = append(enriched, event)
			continue
		}

		venueIdx := slices.IndexFunc(savedVenues, event.Event.Venue.EqualsFields)
		if venueIdx != -1 {
			event.Event.Venue = savedVenues[venueIdx]
		}

		if event.Event.MainAct != nil {
			mainActIdx := slices.IndexFunc(savedArtists, (*event.Event.MainAct).EqualsFields)
			if mainActIdx != -1 {
				event.Event.MainAct = &savedArtists[mainActIdx]
			}
		}

		for i, opener := range event.Event.Openers {
			openerIdx := slices.IndexFunc(savedArtists, opener.EqualsFields)
			if openerIdx != -1 {
				event.Event.Openers[i] = savedArtists[openerIdx]
			}
		}
		enriched = append(enriched, event)
	}

	return enriched
}

func (c Cache) GetLocation() Location {
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
