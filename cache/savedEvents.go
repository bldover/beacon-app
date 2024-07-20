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
	AddEvent(context.Context, data.Event) (string, error)
	DeleteEvent(context.Context, string) error
	ListArtists(context.Context) ([]data.Artist, error)
	AddArtist(context.Context, data.Artist) (string, error)
	DeleteArtist(context.Context, string) error
	ListVenues(context.Context) ([]data.Venue, error)
	AddVenue(context.Context, data.Venue) (string, error)
	DeleteVenue(context.Context, string) error
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

func (c *SavedEventCache) AddSavedEvent(event data.Event) (*data.Event, error) {
	existingIdx := slices.IndexFunc(c.savedEvents, event.Equals)
	if existingIdx >= 0 {
		log.Debugf("Skipping adding event %v because it already existed in the cache", event)
		existing := util.CloneEvent(c.savedEvents[existingIdx])
		return &existing, nil
	}
	if event.MainAct.Populated() {
		artist, err := c.AddArtist(event.MainAct)
		if err != nil {
			return nil, err
		}
		event.MainAct.Id = artist.Id
	}
	for i, opener := range event.Openers {
		artist, err := c.AddArtist(opener)
		if err != nil {
			return nil, err
		}
		event.Openers[i].Id = artist.Id
	}
	venue, err := c.AddVenue(event.Venue)
	if err != nil {
		return nil, err
	}
	event.Venue.Id = venue.Id

	id, err := c.Database.AddEvent(context.Background(), event)
	if err != nil {
		return nil, err
	}

	event.Id = id
	c.savedEvents = append(c.savedEvents, util.CloneEvent(event))
	return &event, nil
}

func (c *SavedEventCache) DeleteSavedEvent(id string) error {
	eventIdx := slices.IndexFunc(c.savedEvents, func(e data.Event) bool {
		return e.Id == id
	})
	if eventIdx == -1 {
		log.Errorf("Unable to find event %v when deleting from cache", id)
		return errors.New("event is not cached")
	}

	if err := c.Database.DeleteEvent(context.Background(), id); err != nil {
		return err
	}

	c.savedEvents = slices.Delete(c.savedEvents, eventIdx, eventIdx+1)
	return nil
}

func (c SavedEventCache) GetArtists() []data.Artist {
	return slices.Clone(c.artists)
}

func (c *SavedEventCache) AddArtist(artist data.Artist) (*data.Artist, error) {
	existingIdx := slices.IndexFunc(c.artists, artist.Equals)
	if existingIdx >= 0 {
		existing := util.CloneArtist(c.artists[existingIdx])
		log.Debugf("Skipping adding artist %v because it already existed in the cache", artist)
		return &existing, nil
	}

	id, err := c.Database.AddArtist(context.Background(), artist)
	if err != nil {
		return nil, err
	}

	artist.Id = id
	c.artists = append(c.artists, util.CloneArtist(artist))
	return &artist, nil
}

func (c *SavedEventCache) DeleteArtist(id string) error {
	artistIdx := slices.IndexFunc(c.artists, func(a data.Artist) bool {
		return a.Id == id
	})
	if artistIdx == -1 {
		log.Errorf("Unable to find artist %v when deleting from cache", id)
		return errors.New("artist is not cached")
	}

	if err := c.Database.DeleteArtist(context.Background(), id); err != nil {
		return err
	}

	c.artists = slices.Delete(c.artists, artistIdx, artistIdx+1)
	return nil
}

func (c SavedEventCache) GetVenues() []data.Venue {
	return util.CloneVenues(c.venues)
}

func (c *SavedEventCache) AddVenue(venue data.Venue) (*data.Venue, error) {
	existingIdx := slices.IndexFunc(c.venues, venue.Equals)
	if existingIdx >= 0 {
		existing := util.CloneVenue(c.venues[existingIdx])
		log.Debugf("Skipping adding venue %v because it already existed in the cache", venue)
		return &existing, nil
	}

	id, err := c.Database.AddVenue(context.Background(), venue)
	if err != nil {
		return nil, err
	}

	venue.Id = id
	c.venues = append(c.venues, util.CloneVenue(venue))
	return &venue, nil
}

func (c *SavedEventCache) DeleteVenue(id string) error {
	venueIdx := slices.IndexFunc(c.venues, func(v data.Venue) bool {
		return v.Id == id
	})
	if venueIdx == -1 {
		log.Errorf("Unable to find venue %v when deleting from cache", id)
		return errors.New("venue is not cached")
	}

	if err := c.Database.DeleteVenue(context.Background(), id); err != nil {
		return err
	}

	c.venues = slices.Delete(c.venues, venueIdx, venueIdx+1)
	return nil
}
