package db

import (
	"concert-manager/data"
	"concert-manager/log"
	"context"
	"errors"
)

type Venue = data.Venue
type Artist = data.Artist
type Event = data.Event

type Context = context.Context

type (
	VenueRepo interface {
		Add(Context, Venue) (string, error)
		Delete(Context, Venue) error
		Exists(Context, Venue) (bool, error)
		FindAll(Context) (*[]Venue, error)
	}
	ArtistRepo interface {
		Add(Context, Artist) (string, error)
		Delete(Context, Artist) error
		Exists(Context, Artist) (bool, error)
		FindAll(Context) (*[]Artist, error)
	}
	EventRepo interface {
		Add(Context, Event) (string, error)
		Delete(Context, Event) error
		Exists(Context, Event) (bool, error)
		FindAll(Context) (*[]Event, error)
	}
	DatabaseRepository struct {
		VenueRepo  VenueRepo
		ArtistRepo ArtistRepo
		EventRepo  EventRepo
	}
)

func (interactor *DatabaseRepository) AddVenue(ctx Context, venue Venue) error {
	log.Debug("Request to add venue", venue)
	if !venue.Populated() {
		log.Debug("Skipping adding venue because required fields are missing", venue)
		return errors.New("failed to create venue due to empty fields")
	}

	_, err := interactor.VenueRepo.Add(ctx, venue)
	if err != nil {
		log.Errorf("Error while adding venue %v, %v\n", venue, err)
		return err
	}
	return nil
}

func (interactor *DatabaseRepository) DeleteVenue(ctx Context, venue Venue) error {
	log.Debug("Request to delete venue", venue)
	err := interactor.VenueRepo.Delete(ctx, venue)
	if err != nil {
		log.Errorf("Error while deleting venue %v, %v\n", venue, err)
		return err
	}
	return nil
}

func (interactor *DatabaseRepository) ListVenues(ctx Context) (*[]Venue, error) {
	log.Debug("Request to list all venues")
    venues, err := interactor.VenueRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all venues,", err)
		return nil, err
	}
	return venues, nil
}

func (interactor *DatabaseRepository) AddArtist(ctx Context, artist Artist) error {
	log.Debug("Request to add artist", artist)
	if !artist.Populated() {
		log.Debug("Skipping adding artist because required fields are missing", artist)
		return errors.New("failed to create artist due to empty fields")
	}

	_, err := interactor.ArtistRepo.Add(ctx, artist)
	if err != nil {
		log.Errorf("Error while adding artist %v, %v\n", artist, err)
		return err
	}
	return nil
}

func (interactor *DatabaseRepository) DeleteArtist(ctx Context, artist Artist) error {
	log.Debug("Request to delete artist", artist)
	err := interactor.ArtistRepo.Delete(ctx, artist)
	if err != nil {
		log.Errorf("Error while deleting artist %v, %v\n", artist, err)
		return err
	}
	return nil
}

func (interactor *DatabaseRepository) ListArtists(ctx Context) (*[]Artist, error) {
	log.Debug("Request to list all artists")
    artists, err := interactor.ArtistRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all artists", err)
		return nil, err
	}
	return artists, nil
}

// Requires that all the artists and the venue already exist
func (interactor *DatabaseRepository) AddEvent(ctx Context, event Event) error {
	log.Debug("Request to add event", event)
	if !event.Populated() {
		log.Debug("Skipping adding event because required fields are missing", event)
		return errors.New("failed to create event due to empty fields")
	}

	_, err := interactor.EventRepo.Add(ctx, event)
	if err != nil {
		log.Errorf("Error while adding event %v, %v\n", event, err)
		return err
	}
	return nil
}

// Creates event and also venue and artists if needed
func (interactor *DatabaseRepository) AddEventRecursive(ctx Context, event Event) error {
	log.Debug("Request to recursively add event", event)
	if !event.Populated() {
		log.Debug("Skipping adding event because required fields are missing", event)
		return errors.New("all fields are required")
	}

	if err := interactor.AddVenue(ctx, event.Venue); err != nil {
		log.Errorf("Failed to add venue %v while recursively adding event, %v\n", event.Venue, err)
		return err
	}
	if event.MainAct.Populated() {
		if err := interactor.AddArtist(ctx, event.MainAct); err != nil {
			log.Errorf("Failed to add artist %v while recursively adding event, %v\n", event.MainAct, err)
			return err
		}
	}
	for _, opener := range event.Openers {
		if err := interactor.AddArtist(ctx, opener); err != nil {
			log.Errorf("Failed to add artist %v while recursively adding event, %v\n", event.MainAct, err)
			return err
		}
	}

	_, err := interactor.EventRepo.Add(ctx, event)
	if err != nil {
		log.Errorf("Error while recursively adding event %v, %v\n", event, err)
		return err
	}
	return nil
}

func (interactor *DatabaseRepository) DeleteEvent(ctx Context, event Event) error {
	log.Debug("Request to delete event", event)
    err := interactor.EventRepo.Delete(ctx, event)
	if err != nil {
		log.Errorf("Error while deleting event %v, %v\n", event, err)
		return err
	}
	return nil
}

func (interactor *DatabaseRepository) ListEvents(ctx Context) (*[]Event, error) {
	log.Debug("Request to list all events")
    events, err := interactor.EventRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all events", err)
		return nil, err
	}
	return events, nil
}
