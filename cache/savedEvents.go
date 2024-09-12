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
	UpdateArtist(context.Context, string, data.Artist) error
	DeleteArtist(context.Context, string) error
	ListVenues(context.Context) ([]data.Venue, error)
	AddVenue(context.Context, data.Venue) (string, error)
	UpdateVenue(context.Context, string, data.Venue) error
	DeleteVenue(context.Context, string) error
}

type SavedEventCache struct {
	Database       Database
	savedEvents    []data.Event
	artists        []data.Artist
	venues         []data.Venue
}

func (c *SavedEventCache) LoadCaches() {
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

func (c *SavedEventCache) RefreshSavedEvents() error {
    log.Info("Refreshing saved event cache")
	savedEvents, err := c.Database.ListEvents(context.Background())
	if err != nil {
		return err
	}
	c.savedEvents = savedEvents
	log.Info("Successfully refreshed saved events")
	return nil
}

func (c *SavedEventCache) RefreshArtists() error {
	log.Info("Refreshing artists cache")
	artists, err := c.Database.ListArtists(context.Background())
	if err != nil {
		return err
	}
	c.artists = artists
	log.Info("Successfully refreshed artists")
	return nil
}

func (c *SavedEventCache) RefreshVenues() error {
	log.Info("Refreshing venues cache")
	venues, err := c.Database.ListVenues(context.Background())
	if err != nil {
		return err
	}
	c.venues = venues
	log.Info("Successfully refreshed venues cache")
	return nil
}

func (c SavedEventCache) GetSavedEvents() []data.Event {
	log.Debug("Retrieving saved events from cache")
	if c.savedEvents == nil {
		return []data.Event{}
	}
	return util.CloneEvents(c.savedEvents)
}

func (c SavedEventCache) GetPassedSavedEvents() []data.Event {
	log.Debug("Retrieving passed saved events from cache")
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
	log.Debug("Adding saved event to cache", event)
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
	log.Debug("Added saved event to cache", event)
	return &event, nil
}

func (c *SavedEventCache) DeleteSavedEvent(id string) error {
	log.Debug("Deleting saved event from cache", id)
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
	log.Debug("Deleted saved event from cache", id)
	return nil
}

func (c SavedEventCache) GetArtists() []data.Artist {
	log.Debug("Retrieving artists from cache")
	return slices.Clone(c.artists)
}

func (c *SavedEventCache) AddArtist(artist data.Artist) (*data.Artist, error) {
	log.Debug("Adding artist to cache", artist)
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
	log.Debug("Added artist to cache", artist)
	return &artist, nil
}

func (c *SavedEventCache) UpdateArtist(id string, artist data.Artist) error {
	log.Debugf("Updating artist in cache, id=%v, %v", id, artist)
	artistIdx := slices.IndexFunc(c.artists, func(a data.Artist) bool {
		return a.Id == id
	})
	if artistIdx == -1 {
		log.Errorf("Unable to find artist %v when updating cache", id)
		return errors.New("artist is not cached")
	}

	err := c.Database.UpdateArtist(context.Background(), id, artist)
	if err != nil {
		return err
	}

	artist.Id = id
	c.artists = slices.Replace(c.artists, artistIdx, artistIdx+1, artist)
	log.Debug("Updated artist in cache", artist)
	return nil
}

func (c *SavedEventCache) DeleteArtist(id string) error {
	log.Debug("Deleting artist from cache", id)
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
	log.Debug("Deleted artist from cache", id)
	return nil
}

func (c SavedEventCache) GetVenues() []data.Venue {
	log.Debug("Retrieving venues from cache")
	return util.CloneVenues(c.venues)
}

func (c *SavedEventCache) AddVenue(venue data.Venue) (*data.Venue, error) {
	log.Debug("Adding venue to cache", venue)
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
	log.Debug("Added venue to cache", venue)
	return &venue, nil
}

func (c *SavedEventCache) UpdateVenue(id string, venue data.Venue) error {
	log.Debugf("Updating venue in cache, id=%v, %v", id, venue)
	venueIdx := slices.IndexFunc(c.venues, func(a data.Venue) bool {
		return a.Id == id
	})
	if venueIdx == -1 {
		log.Errorf("Unable to find venue %v when updating cache", id)
		return errors.New("venue is not cached")
	}

	err := c.Database.UpdateVenue(context.Background(), id, venue)
	if err != nil {
		return err
	}

	venue.Id = id
	c.venues = slices.Replace(c.venues, venueIdx, venueIdx+1, venue)
	log.Debug("Updated venue in cache", venue)
	return nil
}

func (c *SavedEventCache) DeleteVenue(id string) error {
	log.Debug("Deleting venue from cache", id)
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
	log.Debug("Deleted venue from cache", id)
	return nil
}
