package cache

import (
	"concert-manager/data"
	"concert-manager/db"
	"concert-manager/finder"
	"concert-manager/log"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Finder interface {
	FindAllEvents(request finder.FindEventRequest) ([]data.EventDetails, error)
}

type Database interface {
	ListEvents(context.Context) ([]data.Event, error)
	AddEventRecursive(context.Context, data.Event) error
	DeleteEvent(context.Context, data.Event) error
	ListArtists(context.Context) ([]data.Artist, error)
	AddArtist(context.Context, data.Artist) error
	DeleteArtist(context.Context, data.Artist) error
	ListVenues(context.Context) ([]data.Venue, error)
	AddVenue(context.Context, data.Venue) error
	DeleteVenue(context.Context, data.Venue) error
}

type LocalCache struct {
	db             Database
	finder         Finder
	savedEvents    map[int][]data.Event
	upcomingEvents map[string][]data.EventDetails
	artists        []data.Artist
	venues         []data.Venue
}

const (
	savedEventsInitCapacity = 300
	upcomingEventsInitCapacity = 500
)

func NewLocalCache(db *db.DatabaseRepository, finder *finder.EventFinder) *LocalCache {
	cache := LocalCache{}
	cache.db = db
	cache.finder = finder
	cache.savedEvents = make(map[int][]data.Event, savedEventsInitCapacity)
	cache.upcomingEvents = make(map[string][]data.EventDetails, upcomingEventsInitCapacity)
	cache.Sync()
	return &cache
}

func (c *LocalCache) Sync() {
	events, err := c.db.ListEvents(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize events:", err)
	}

	pastEvents := []data.Event{}
	futureEvents := []data.Event{}
	for _, e := range events {
		if data.ValidFutureDate(e.Date) {
			futureEvents = append(futureEvents, e)
		} else {
			pastEvents = append(pastEvents, e)
		}
	}
	c.savedEvents[data.Past] = pastEvents
	c.savedEvents[data.Future] = futureEvents
	log.Info("Successfully initialized events")

	artists, err := c.db.ListArtists(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize artists:", err)
	}
	c.artists = artists
	log.Info("Successfully initialized artists")

	venues, err := c.db.ListVenues(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize venues:", err)
	}
	c.venues = venues
	log.Info("Successfully initialized venues")
}

func (c LocalCache) GetFutureEvents() []data.Event {
	return c.savedEvents[data.Future]
}

func (c LocalCache) GetPastEvents() []data.Event {
	return slices.Clone(c.savedEvents[data.Past])
}

func (c *LocalCache) AddEvent(event data.Event) error {
	key := data.Past
	if data.ValidFutureDate(event.Date) {
		key = data.Future
	}

	if slices.ContainsFunc(c.savedEvents[key], event.Equals) {
		log.Debugf("Skipping adding event %v because it already existed in the cache", event)
		return nil
	}

	if err := c.AddArtist(event.MainAct); err != nil {
		return err
	}
	for _, opener := range event.Openers {
		if err := c.AddArtist(opener); err != nil {
			return err
		}
	}
	if err := c.AddVenue(event.Venue); err != nil {
		return err
	}

	if err := c.db.AddEventRecursive(context.Background(), event); err != nil {
		return err
	}

	events := c.savedEvents[key]
	updatedEvents := append(events, event)
	c.savedEvents[key] = updatedEvents
	return nil
}

func (c *LocalCache) DeleteEvent(event data.Event) error {
	key := data.Past
	if data.ValidFutureDate(event.Date) {
		key = data.Future
	}

	eventIdx := slices.IndexFunc(c.savedEvents[key], event.Equals)
	if eventIdx == -1 {
		log.Errorf("Unable to find event %v when deleting from cache", event)
		return errors.New("event is not cached")
	}

	if err := c.db.DeleteEvent(context.Background(), event); err != nil {
		return err
	}

	c.savedEvents[key] = slices.Delete(c.savedEvents[key], eventIdx, eventIdx+1)
	return nil
}

func (c LocalCache) GetArtists() []data.Artist {
	return slices.Clone(c.artists)
}

func (c *LocalCache) AddArtist(artist data.Artist) error {
	if slices.Contains(c.artists, artist) {
		log.Debugf("Skipping adding artist %v because it already existed in the cache", artist)
		return nil
	}

	if err := c.db.AddArtist(context.Background(), artist); err != nil {
		return err
	}

	c.artists = append(c.artists, artist)
	return nil
}

func (c *LocalCache) DeleteArtist(artist data.Artist) error {
	artistIdx := slices.Index(c.artists, artist)
	if artistIdx == -1 {
		log.Errorf("Unable to find artist %v when deleting from cache", artist)
		return errors.New("artist is not cached")
	}

	if err := c.db.DeleteArtist(context.Background(), artist); err != nil {
		return err
	}

	c.artists = slices.Delete(c.artists, artistIdx, artistIdx+1)
	return nil
}

func (c LocalCache) GetVenues() []data.Venue {
	return slices.Clone(c.venues)
}

func (c *LocalCache) AddVenue(venue data.Venue) error {
	if slices.Contains(c.venues, venue) {
		log.Debugf("Skipping adding venue %v because it already existed in the cache", venue)
		return nil
	}

	if err := c.db.AddVenue(context.Background(), venue); err != nil {
		return err
	}

	c.venues = append(c.venues, venue)
	return nil
}

func (c *LocalCache) DeleteVenue(venue data.Venue) error {
	venueIdx := slices.Index(c.venues, venue)
	if venueIdx == -1 {
		log.Errorf("Unable to find venue %v when deleting from cache", venue)
		return errors.New("venue is not cached")
	}

	if err := c.db.DeleteVenue(context.Background(), venue); err != nil {
		return err
	}

	c.venues = slices.Delete(c.venues, venueIdx, venueIdx+1)
	return nil
}

func (c LocalCache) GetUpcomingEvents(city, stateCode string) []data.EventDetails {
	key := buildUpcomingEventKey(city, stateCode)
	if events, ok := c.upcomingEvents[key]; ok {
		return events
	}
	c.ReloadUpcomingEvents(city, stateCode)
	return slices.Clone(c.upcomingEvents[key])
}

func (c *LocalCache) ReloadUpcomingEvents(city, stateCode string) error {
	request := finder.FindEventRequest{City: city, State: stateCode}
	events, err := c.finder.FindAllEvents(request)
	if err != nil {
		log.Errorf("Error while loading events for %s, %s: %v", city, stateCode, err)
	}

	key := buildUpcomingEventKey(city, stateCode)
	c.upcomingEvents[key] = events
	return err
}

func buildUpcomingEventKey(city, stateCode string) string {
	return fmt.Sprintf("%s#%s", strings.ToLower(city), strings.ToLower(stateCode))
}
