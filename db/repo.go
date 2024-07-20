package db

import (
	"concert-manager/data"
	"concert-manager/log"
	"context"
	"errors"
)

type (
	VenueRepo interface {
		Add(context.Context, data.Venue) (string, error)
		Delete(context.Context, string) error
		Exists(context.Context, data.Venue) (bool, error)
		FindAll(context.Context) ([]data.Venue, error)
	}
	ArtistRepo interface {
		Add(context.Context, data.Artist) (string, error)
		Delete(context.Context, string) error
		Exists(context.Context, data.Artist) (bool, error)
		FindAll(context.Context) ([]data.Artist, error)
	}
	EventRepo interface {
		Add(context.Context, data.Event) (string, error)
		Delete(context.Context, string) error
		Exists(context.Context, data.Event) (bool, error)
		FindAll(context.Context) ([]data.Event, error)
	}
	DatabaseRepository struct {
		VenueRepo  VenueRepo
		ArtistRepo ArtistRepo
		EventRepo  EventRepo
	}
)

func (interactor *DatabaseRepository) AddVenue(ctx context.Context, venue data.Venue) (string, error) {
	log.Debug("Request to add venue", venue)
	if !venue.Populated() {
		log.Debug("Skipping adding venue because required fields are missing", venue)
		return "", errors.New("failed to create venue due to empty fields")
	}

	id, err := interactor.VenueRepo.Add(ctx, venue)
	if err != nil {
		log.Errorf("Error while adding venue %v, %v\n", venue, err)
		return "", err
	}
	return id, nil
}

func (interactor *DatabaseRepository) DeleteVenue(ctx context.Context, id string) error {
	log.Debug("Request to delete venue", id)
	err := interactor.VenueRepo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting venue %v, %v\n", id, err)
		return err
	}
	return nil
}

func (interactor *DatabaseRepository) ListVenues(ctx context.Context) ([]data.Venue, error) {
	log.Debug("Request to list all venues")
    venues, err := interactor.VenueRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all venues,", err)
		return nil, err
	}
	return venues, nil
}

func (interactor *DatabaseRepository) AddArtist(ctx context.Context, artist data.Artist) (string, error) {
	log.Debug("Request to add artist", artist)
	if !artist.Populated() {
		log.Debug("Skipping adding artist because required fields are missing", artist)
		return "", errors.New("failed to create artist due to empty fields")
	}

	id, err := interactor.ArtistRepo.Add(ctx, artist)
	if err != nil {
		log.Errorf("Error while adding artist %v, %v\n", artist, err)
		return "", err
	}
	return id, nil
}

func (interactor *DatabaseRepository) DeleteArtist(ctx context.Context, id string) error {
	log.Debug("Request to delete artist", id)
	err := interactor.ArtistRepo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting artist %v, %v\n", id, err)
		return err
	}
	return nil
}

func (interactor *DatabaseRepository) ListArtists(ctx context.Context) ([]data.Artist, error) {
	log.Debug("Request to list all artists")
    artists, err := interactor.ArtistRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all artists", err)
		return nil, err
	}
	return artists, nil
}

// Requires that all the artists and the venue already exist
func (interactor *DatabaseRepository) AddEvent(ctx context.Context, event data.Event) (string, error) {
	log.Debug("Request to add event", event)
	if !event.Populated() {
		log.Debug("Skipping adding event because required fields are missing", event)
		return "", errors.New("failed to create event due to empty fields")
	}

	id, err := interactor.EventRepo.Add(ctx, event)
	if err != nil {
		log.Errorf("Error while adding event %v, %v\n", event, err)
		return "", err
	}
	return id, nil
}

func (interactor *DatabaseRepository) DeleteEvent(ctx context.Context, id string) error {
	log.Debug("Request to delete event", id)
    err := interactor.EventRepo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting event %v, %v\n", id, err)
		return err
	}
	return nil
}

func (interactor *DatabaseRepository) ListEvents(ctx context.Context) ([]data.Event, error) {
	log.Debug("Request to list all events")
    events, err := interactor.EventRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all events", err)
		return nil, err
	}
	return events, nil
}
