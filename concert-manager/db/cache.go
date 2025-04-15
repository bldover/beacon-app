package db

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
	AddEvent(context.Context, data.Event) (data.Event, error)
	DeleteEvent(context.Context, string) error
	ListArtists(context.Context) ([]data.Artist, error)
	AddArtist(context.Context, data.Artist) (data.Artist, error)
	UpdateArtist(context.Context, data.Artist) (data.Artist, error)
	DeleteArtist(context.Context, string) error
	ListVenues(context.Context) ([]data.Venue, error)
	AddVenue(context.Context, data.Venue) (data.Venue, error)
	UpdateVenue(context.Context, data.Venue) (data.Venue, error)
	DeleteVenue(context.Context, string) error
}

type Cache struct {
	Database       Database
	savedEvents    []data.Event
	artists        []data.Artist
	venues         []data.Venue
}

func (c *Cache) LoadCaches() {
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

func (c *Cache) RefreshSavedEvents() error {
    log.Info("Refreshing saved event cache")
	savedEvents, err := c.Database.ListEvents(context.Background())
	if err != nil {
		return err
	}
	c.savedEvents = savedEvents
	log.Info("Successfully refreshed saved events")
	return nil
}

func (c *Cache) RefreshArtists() error {
	log.Info("Refreshing artists cache")
	artists, err := c.Database.ListArtists(context.Background())
	if err != nil {
		return err
	}
	c.artists = artists
	log.Info("Successfully refreshed artists")
	return nil
}

func (c *Cache) RefreshVenues() error {
	log.Info("Refreshing venues cache")
	venues, err := c.Database.ListVenues(context.Background())
	if err != nil {
		return err
	}
	c.venues = venues
	log.Info("Successfully refreshed venues cache")
	return nil
}

func (c Cache) GetSavedEvents() []data.Event {
	log.Debug("Retrieving saved events from cache")
	if c.savedEvents == nil {
		return []data.Event{}
	}
	return util.CloneEvents(c.savedEvents)
}

func (c Cache) GetPassedSavedEvents() []data.Event {
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

func (c *Cache) AddSavedEvent(event data.Event) (*data.Event, error) {
	log.Debug("Adding saved event to cache", event)
	existingIdx := slices.IndexFunc(c.savedEvents, event.Equals)
	if existingIdx >= 0 {
		log.Debugf("Skipping adding event %v because it already existed in the cache", event)
		existing := util.CloneEvent(c.savedEvents[existingIdx])
		return &existing, nil
	}
	if event.MainAct != nil && event.MainAct.Populated() {
		artist, err := c.AddArtist(*event.MainAct)
		if err != nil {
			return nil, err
		}
		event.MainAct = artist
	}
	for i, opener := range event.Openers {
		artist, err := c.AddArtist(opener)
		if err != nil {
			return nil, err
		}
		event.Openers[i] = *artist
	}
	venue, err := c.AddVenue(event.Venue)
	if err != nil {
		return nil, err
	}
	event.Venue = *venue

	newEvent, err := c.Database.AddEvent(context.Background(), event)
	if err != nil {
		return nil, err
	}

	c.savedEvents = append(c.savedEvents, newEvent)
	log.Debug("Added saved event to cache", newEvent)
	return &newEvent, nil
}

func (c *Cache) DeleteSavedEvent(id string) error {
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

func (c Cache) GetArtists() []data.Artist {
	log.Debug("Retrieving artists from cache")
	return slices.Clone(c.artists)
}

func (c *Cache) AddArtist(artist data.Artist) (*data.Artist, error) {
	log.Debug("Adding artist to cache", artist)
	existingIdx := slices.IndexFunc(c.artists, artist.Equals)
	if existingIdx >= 0 {
		existing := util.CloneArtist(c.artists[existingIdx])
		log.Debugf("Skipping adding artist %v because it already existed in the cache", artist)
		return &existing, nil
	}

	newArtist, err := c.Database.AddArtist(context.Background(), artist)
	if err != nil {
		return nil, err
	}

	c.artists = append(c.artists, newArtist)
	log.Debug("Added artist to cache", newArtist)
	return &newArtist, nil
}

func (c *Cache) UpdateArtist(id string, artist data.Artist) error {
	log.Debugf("Updating artist in cache, id=%v, %v", id, artist)
	artistIdx := slices.IndexFunc(c.artists, func(a data.Artist) bool {
		return a.Id == id
	})
	if artistIdx == -1 {
		log.Errorf("Unable to find artist %v when updating cache", id)
		return errors.New("artist is not cached")
	}

	artist.Id = id
	updatedArtist, err := c.Database.UpdateArtist(context.Background(), artist)
	if err != nil {
		return err
	}

	c.artists = slices.Replace(c.artists, artistIdx, artistIdx+1, updatedArtist)

	for i, event := range c.savedEvents {
		if event.MainAct != nil && event.MainAct.Id == id {
			c.savedEvents[i].MainAct = &updatedArtist
		}
		for j, opener := range event.Openers {
			if opener.Id == id {
				c.savedEvents[i].Openers[j] = updatedArtist
			}
		}
	}

	log.Debug("Updated artist in cache", updatedArtist)
	return nil
}

func (c *Cache) DeleteArtist(id string) error {
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

func (c Cache) GetVenues() []data.Venue {
	log.Debug("Retrieving venues from cache")
	return util.CloneVenues(c.venues)
}

func (c *Cache) AddVenue(venue data.Venue) (*data.Venue, error) {
	log.Debug("Adding venue to cache", venue)
	existingIdx := slices.IndexFunc(c.venues, venue.Equals)
	if existingIdx >= 0 {
		existing := util.CloneVenue(c.venues[existingIdx])
		log.Debugf("Skipping adding venue %v because it already existed in the cache", venue)
		return &existing, nil
	}

	newVenue, err := c.Database.AddVenue(context.Background(), venue)
	if err != nil {
		return nil, err
	}

	c.venues = append(c.venues, newVenue)
	log.Debug("Added venue to cache", newVenue)
	return &newVenue, nil
}

func (c *Cache) UpdateVenue(id string, venue data.Venue) error {
	log.Debugf("Updating venue in cache, id=%v, %v", id, venue)
	venueIdx := slices.IndexFunc(c.venues, func(a data.Venue) bool {
		return a.Id == id
	})
	if venueIdx == -1 {
		log.Errorf("Unable to find venue %v when updating cache", id)
		return errors.New("venue is not cached")
	}

	venue.Id = id
	updatedVenue, err := c.Database.UpdateVenue(context.Background(), venue)
	if err != nil {
		return err
	}

	c.venues = slices.Replace(c.venues, venueIdx, venueIdx+1, updatedVenue)

	for i, event := range c.savedEvents {
		if event.Venue.Id == id {
			c.savedEvents[i].Venue = updatedVenue
		}
	}

	log.Debug("Updated venue in cache", updatedVenue)
	return nil
}

func (c *Cache) DeleteVenue(id string) error {
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
