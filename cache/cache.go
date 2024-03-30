package cache

import (
	"concert-manager/data"
	"concert-manager/finder"
	"concert-manager/log"
	"concert-manager/util"
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
	Database       Database
	Finder         Finder
	savedEvents    []data.Event
	upcomingEvents map[string][]data.EventDetails
	artists        []data.Artist
	venues         []data.Venue
}

type eventType int

const (
	past = iota
	future
)

const (
	savedEventsInitCapacity    = 300
	upcomingEventsInitCapacity = 500
)

func NewLocalCache() *LocalCache {
	cache := LocalCache{}
	cache.savedEvents = make([]data.Event, savedEventsInitCapacity)
	cache.upcomingEvents = make(map[string][]data.EventDetails, upcomingEventsInitCapacity)
	return &cache
}

func (c *LocalCache) Sync() {
	savedEvents, err := c.Database.ListEvents(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize events:", err)
	}
	c.savedEvents = savedEvents
	log.Info("Successfully initialized saved events")

	artists, err := c.Database.ListArtists(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize artists:", err)
	}
	c.artists = artists
	log.Info("Successfully initialized artists")

	venues, err := c.Database.ListVenues(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize venues:", err)
	}
	c.venues = venues
	log.Info("Successfully initialized venues")
}

func (c LocalCache) GetSavedEvents() []data.Event {
	savedEvents := slices.Clone(c.savedEvents)
	for i := range savedEvents {
		savedEvents[i].Openers = slices.Clone(savedEvents[i].Openers)
	}
	return savedEvents
}

func (c LocalCache) GetPassedSavedEvents() []data.Event {
	passedEvents := []data.Event{}
	for _, event := range c.GetSavedEvents() {
		if util.PastDate(event.Date) && !event.Purchased {
			passedEvents = append(passedEvents, event)
		}
	}
	for i := range passedEvents {
		passedEvents[i].Openers = slices.Clone(passedEvents[i].Openers)
	}
	return passedEvents
}

func (c *LocalCache) AddSavedEvent(event data.Event) error {
	if slices.ContainsFunc(c.savedEvents, event.Equals) {
		log.Debugf("Skipping adding event %v because it already existed in the cache", event)
		return nil
	}
	if event.MainAct.Populated() {
		if err := c.AddArtist(event.MainAct); err != nil {
			return err
		}
	}
	for _, opener := range event.Openers {
		if err := c.AddArtist(opener); err != nil {
			return err
		}
	}
	if err := c.AddVenue(event.Venue); err != nil {
		return err
	}

	if err := c.Database.AddEventRecursive(context.Background(), event); err != nil {
		return err
	}

	events := c.savedEvents
	updatedEvents := append(events, event)
	c.savedEvents = updatedEvents
	return nil
}

func (c *LocalCache) DeleteSavedEvent(event data.Event) error {
	eventIdx := slices.IndexFunc(c.savedEvents, event.Equals)
	if eventIdx == -1 {
		log.Errorf("Unable to find event %v when deleting from cache", event)
		return errors.New("event is not cached")
	}

	if err := c.Database.DeleteEvent(context.Background(), event); err != nil {
		return err
	}

	c.savedEvents = slices.Delete(c.savedEvents, eventIdx, eventIdx+1)
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

	if err := c.Database.AddArtist(context.Background(), artist); err != nil {
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

	if err := c.Database.DeleteArtist(context.Background(), artist); err != nil {
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

	if err := c.Database.AddVenue(context.Background(), venue); err != nil {
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

	if err := c.Database.DeleteVenue(context.Background(), venue); err != nil {
		return err
	}

	c.venues = slices.Delete(c.venues, venueIdx, venueIdx+1)
	return nil
}

func (c LocalCache) GetUpcomingEvents(city, stateCode string) []data.EventDetails {
	key := buildUpcomingEventKey(city, stateCode)
	var events []data.EventDetails
	if upcomingEvents, ok := c.upcomingEvents[key]; ok {
		events = slices.Clone(upcomingEvents)
	} else {
		c.ReloadUpcomingEvents(city, stateCode)
		events = slices.Clone(c.upcomingEvents[key])
	}

	for i := range events {
		events[i].Event.Openers = slices.Clone(events[i].Event.Openers)
	}
	return events
}

func (c *LocalCache) ReloadUpcomingEvents(city, stateCode string) error {
	request := finder.FindEventRequest{City: city, State: stateCode}
	events, err := c.Finder.FindAllEvents(request)
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
