package cache

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/util"
	"context"
	"errors"
	"slices"
)

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

type SavedEventCache struct {
	Database       Database
	savedEvents    []data.Event
	artists        []data.Artist
	venues         []data.Venue
}

func (c *SavedEventCache) Sync() {
	log.Info("Initializing saved event cache")
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

	log.Info("Finished initializing saved event cache")
}

func (c SavedEventCache) GetSavedEvents() []data.Event {
	if c.savedEvents == nil {
		return []data.Event{}
	}
	return util.CloneEvents(c.savedEvents)
}

func (c SavedEventCache) GetPassedSavedEvents() []data.Event {
	if c.savedEvents == nil {
		return []data.Event{}
	}

	passedEvents := []data.Event{}
	for _, event := range c.GetSavedEvents() {
		if util.PastDate(event.Date) && !event.Purchased {
			passedEvents = append(passedEvents, util.CloneEvent(event))
		}
	}
	return passedEvents
}

func (c *SavedEventCache) AddSavedEvent(event data.Event) error {
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

	c.savedEvents = append(c.savedEvents, util.CloneEvent(event))
	return nil
}

func (c *SavedEventCache) DeleteSavedEvent(event data.Event) error {
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

func (c SavedEventCache) GetArtists() []data.Artist {
	return slices.Clone(c.artists)
}

func (c *SavedEventCache) AddArtist(artist data.Artist) error {
	if slices.Contains(c.artists, artist) {
		log.Debugf("Skipping adding artist %v because it already existed in the cache", artist)
		return nil
	}

	if err := c.Database.AddArtist(context.Background(), artist); err != nil {
		return err
	}

	c.artists = append(c.artists, util.CloneArtist(artist))
	return nil
}

func (c *SavedEventCache) DeleteArtist(artist data.Artist) error {
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

func (c SavedEventCache) GetVenues() []data.Venue {
	return util.CloneVenues(c.venues)
}

func (c *SavedEventCache) AddVenue(venue data.Venue) error {
	if slices.Contains(c.venues, venue) {
		log.Debugf("Skipping adding venue %v because it already existed in the cache", venue)
		return nil
	}

	if err := c.Database.AddVenue(context.Background(), venue); err != nil {
		return err
	}

	c.venues = append(c.venues, util.CloneVenue(venue))
	return nil
}

func (c *SavedEventCache) DeleteVenue(venue data.Venue) error {
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
