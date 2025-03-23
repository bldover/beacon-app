package cache

import (
	"concert-manager/data"
	"concert-manager/finder"
	"concert-manager/log"
	"concert-manager/util"
	"fmt"
	"strings"
	"time"
)

type Finder interface {
	FindAllEvents(request finder.FindEventRequest) ([]data.EventDetails, error)
}

type Ranker interface {
	Rank(data.EventDetails) data.EventRank
}

type upcomingEventsData struct {
	events     []data.EventDetails
	lastLoaded time.Time
}

type eventRanksData struct {
	ranks      []data.EventRank
	lastLoaded time.Time
}

type UpcomingEventCache struct {
	Location       Location
	Finder         Finder
	Ranker         Ranker
	upcomingEvents map[string]upcomingEventsData
	eventRanks     map[string]eventRanksData
}

const (
	defaultCity      = "Atlanta"
	defaultStateCode = "GA"
)

func NewUpcomingEventCache() *UpcomingEventCache {
	cache := UpcomingEventCache{}
	cache.Location = Location{City: defaultCity, StateCode: defaultStateCode}
	cache.upcomingEvents = map[string]upcomingEventsData{}
	cache.eventRanks = map[string]eventRanksData{}
	return &cache
}

var	upcomingEventTTL, _ = time.ParseDuration("24h")

func (c *UpcomingEventCache) GetUpcomingEvents() []data.EventDetails {
	key := c.Location.key()
	if d, ok := c.upcomingEvents[key]; !ok {
		c.doRefresh()
	} else if isExpired(d.lastLoaded, upcomingEventTTL) {
		go c.doRefresh()
	}
	return util.CloneEventDetails(c.upcomingEvents[key].events)
}

func (c *UpcomingEventCache) doRefresh() {
    err := c.RefreshUpcomingEvents()
	if err != nil {
		log.Error("Failed to refresh upcoming events", err)
	}
}

func (c *UpcomingEventCache) RefreshUpcomingEvents() error {
	log.Info("Refreshing upcoming events")
	loc := c.Location
	key := c.Location.key()
	request := finder.FindEventRequest{City: loc.City, State: loc.StateCode}
	events, err := c.Finder.FindAllEvents(request)
	if err != nil {
		if _, ok := c.upcomingEvents[key]; !ok {
			eventData := upcomingEventsData{events: []data.EventDetails{}, lastLoaded: time.Time{}}
			c.upcomingEvents[key] = eventData
		}
		return err
	}

	eventData := upcomingEventsData{events: events, lastLoaded: time.Now().Round(0)}
	c.upcomingEvents[key] = eventData
	return nil
}

func (c *UpcomingEventCache) Invalidate() {
	c.upcomingEvents = map[string]upcomingEventsData{}
	c.eventRanks = map[string]eventRanksData{}
}

type Location struct {
	City      string
	StateCode string
}

func (c UpcomingEventCache) GetLocation() Location {
	return c.Location
}

func (c *UpcomingEventCache) ChangeLocation(city, stateCode string) {
	loc := Location{City: city, StateCode: stateCode}
	log.Debugf("Updating upcomingEventCache location from %s to %s", c.Location, loc)
	c.Location = loc
}

func (c Location) key() string {
	return fmt.Sprintf("%s#%s", strings.ToLower(c.City), strings.ToLower(c.StateCode))
}

func (c Location) String() string {
    return fmt.Sprintf("%s, %s", c.City, c.StateCode)
}

type Threshold float64

const (
	HighThreshold   = Threshold(0.15)
	MediumThreshold = Threshold(0.08)
	LowThreshold    = Threshold(0.03)
	NoThreshold     = Threshold(0)
)

func (t Threshold) Level() string {
	switch t {
	case HighThreshold:
		return "High"
	case MediumThreshold:
		return "Medium"
	case LowThreshold:
		return "Low"
	case NoThreshold:
		return "None"
	default:
		return "Invalid"
	}
}

func (c *UpcomingEventCache) GetRecommendedEvents(threshold Threshold) []data.EventRank {
	key := c.Location.key()
	if d, ok := c.eventRanks[key]; !ok {
		c.LoadRecommendations()
	} else if isExpired(d.lastLoaded, upcomingEventTTL) {
		go c.LoadRecommendations()
	}

	var events []data.EventRank
	for _, event := range c.eventRanks[key].ranks {
		if event.Rank >= float64(threshold) {
			events = append(events, util.CloneEventRank(event))
		}
	}
	return events
}

func (c *UpcomingEventCache) LoadRecommendations() {
	rankedEvents := []data.EventRank{}
	events := c.GetUpcomingEvents()
	for _, event := range events {
		rankedEvent := c.Ranker.Rank(event)
		rankedEvents = append(rankedEvents, rankedEvent)
	}

	key := c.Location.key()
	ranksData := eventRanksData{ranks: rankedEvents, lastLoaded: time.Now()}
	c.eventRanks[key] = ranksData
}

func isExpired(lastLoad time.Time, ttl time.Duration) bool {
	elapsedTime := time.Since(lastLoad)
	log.Debugf("Upcoming lastLoaded: %v, now: %v, elapsed: %v, ttl: %v", lastLoad, time.Now(), elapsedTime, ttl)
	return elapsedTime > ttl
}
