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
		Update(context.Context, string, data.Venue) error
		Delete(context.Context, string) error
		Exists(context.Context, data.Venue) (bool, error)
		FindAll(context.Context) ([]data.Venue, error)
	}
	ArtistRepo interface {
		Add(context.Context, data.Artist) (string, error)
		Update(context.Context, string, data.Artist) error
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

func (r *DatabaseRepository) AddVenue(ctx context.Context, venue data.Venue) (string, error) {
	log.Debug("Request to add venue", venue)
	if !venue.Populated() {
		log.Debug("Skipping adding venue because required fields are missing", venue)
		return "", errors.New("failed to create venue due to empty fields")
	}

	id, err := r.VenueRepo.Add(ctx, venue)
	if err != nil {
		log.Errorf("Error while adding venue %v, %v\n", venue, err)
		return "", err
	}
	return id, nil
}

func (r *DatabaseRepository) UpdateVenue(ctx context.Context, id string, venue data.Venue) error {
	log.Debug("Request to update venue", id, venue)
	err := r.VenueRepo.Update(ctx, id, venue)
	if err != nil {
		log.Errorf("Error while updating venue %v, %v\n", id, err)
		return err
	}
	return nil
}

func (r *DatabaseRepository) DeleteVenue(ctx context.Context, id string) error {
	log.Debug("Request to delete venue", id)
	err := r.VenueRepo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting venue %v, %v\n", id, err)
		return err
	}
	return nil
}

func (r *DatabaseRepository) ListVenues(ctx context.Context) ([]data.Venue, error) {
	log.Debug("Request to list all venues")
    venues, err := r.VenueRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all venues,", err)
		return nil, err
	}
	return venues, nil
}

func (r *DatabaseRepository) AddArtist(ctx context.Context, artist data.Artist) (string, error) {
	log.Debug("Request to add artist", artist)
	if !artist.Populated() {
		log.Debug("Skipping adding artist because required fields are missing", artist)
		return "", errors.New("failed to create artist due to empty fields")
	}

	id, err := r.ArtistRepo.Add(ctx, artist)
	if err != nil {
		log.Errorf("Error while adding artist %v, %v\n", artist, err)
		return "", err
	}
	return id, nil
}

func (r *DatabaseRepository) UpdateArtist(ctx context.Context, id string, artist data.Artist) error {
	log.Debug("Request to update artist", id, artist)
	err := r.ArtistRepo.Update(ctx, id, artist)
	if err != nil {
		log.Errorf("Error while updating artist %v, %v\n", id, err)
		return err
	}
	return nil
}

func (r *DatabaseRepository) DeleteArtist(ctx context.Context, id string) error {
	log.Debug("Request to delete artist", id)
	err := r.ArtistRepo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting artist %v, %v\n", id, err)
		return err
	}
	return nil
}

func (r *DatabaseRepository) ListArtists(ctx context.Context) ([]data.Artist, error) {
	log.Debug("Request to list all artists")
    artists, err := r.ArtistRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all artists", err)
		return nil, err
	}
	return artists, nil
}

// Requires that all the artists and the venue already exist
func (r *DatabaseRepository) AddEvent(ctx context.Context, event data.Event) (string, error) {
	log.Debug("Request to add event", event)
	if !event.Populated() {
		log.Debug("Skipping adding event because required fields are missing", event)
		return "", errors.New("failed to create event due to empty fields")
	}

	id, err := r.EventRepo.Add(ctx, event)
	if err != nil {
		log.Errorf("Error while adding event %v, %v\n", event, err)
		return "", err
	}
	return id, nil
}

func (r *DatabaseRepository) DeleteEvent(ctx context.Context, id string) error {
	log.Debug("Request to delete event", id)
    err := r.EventRepo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting event %v, %v\n", id, err)
		return err
	}
	return nil
}

func (r *DatabaseRepository) ListEvents(ctx context.Context) ([]data.Event, error) {
	log.Debug("Request to list all events")
    events, err := r.EventRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all events", err)
		return nil, err
	}
	return events, nil
}
