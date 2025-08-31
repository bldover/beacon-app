package db

import (
	"concert-manager/domain"
	"concert-manager/log"

	"context"
	"errors"
)

type (
	VenueDatabase interface {
		Add(context.Context, domain.Venue) (string, error)
		Update(context.Context, domain.Venue) error
		Delete(context.Context, string) error
		FindAll(context.Context) ([]domain.Venue, error)
	}
	ArtistDatabase interface {
		Add(context.Context, domain.Artist) (string, error)
		Update(context.Context, domain.Artist) error
		Delete(context.Context, string) error
		FindAll(context.Context) ([]domain.Artist, error)
	}
	EventDatabase interface {
		Add(context.Context, domain.Event) (string, error)
		Delete(context.Context, string) error
		FindAll(context.Context) ([]domain.Event, error)
	}
	EventRepository struct {
		VenueRepo  VenueDatabase
		ArtistRepo ArtistDatabase
		EventRepo  EventDatabase
	}
)

func (r *EventRepository) AddVenue(ctx context.Context, venue domain.Venue) (domain.Venue, error) {
	log.Debug("Request to add venue", venue)
	if !venue.Populated() {
		log.Debug("Skipping adding venue because required fields are missing", venue)
		return venue, errors.New("failed to create venue due to empty fields")
	}
	newVenue := domain.CloneVenue(venue)
	id, err := r.VenueRepo.Add(ctx, newVenue)
	if err != nil {
		log.Errorf("Error while adding venue %v, %v\n", venue, err)
		return venue, err
	}
	newVenue.ID.Primary = id
	log.Debug("Added venue to database", newVenue)
	return newVenue, nil
}

func (r *EventRepository) UpdateVenue(ctx context.Context, venue domain.Venue) (domain.Venue, error) {
	log.Debug("Request to update venue", venue)
	updateVenue := domain.CloneVenue(venue)
	err := r.VenueRepo.Update(ctx, updateVenue)
	if err != nil {
		log.Errorf("Error while updating venue %v, %v\n", venue, err)
		return venue, err
	}
	log.Debug("Updated  venue in database", updateVenue)
	return updateVenue, nil
}

func (r *EventRepository) DeleteVenue(ctx context.Context, id string) error {
	log.Debug("Request to delete venue", id)
	err := r.VenueRepo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting venue %v, %v\n", id, err)
		return err
	}
	log.Debug("Deleted venue from database", id)
	return nil
}

func (r *EventRepository) ListVenues(ctx context.Context) ([]domain.Venue, error) {
	log.Debug("Request to list all venues")
	venues, err := r.VenueRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all venues,", err)
		return nil, err
	}
	return venues, nil
}

func (r *EventRepository) AddArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error) {
	log.Debug("Request to add artist", artist)
	if !artist.Populated() {
		log.Debug("Skipping adding artist because required fields are missing", artist)
		return artist, errors.New("failed to create artist due to empty fields")
	}
	newArtist := domain.CloneArtist(artist)
	if newArtist.Genres.Spotify == nil {
		newArtist.Genres.Spotify = []string{}
	}
	if newArtist.Genres.Ticketmaster == nil {
		newArtist.Genres.Ticketmaster = []string{}
	}
	if newArtist.Genres.LastFm == nil {
		newArtist.Genres.LastFm = []string{}
	}
	if newArtist.Genres.User == nil {
		newArtist.Genres.User = []string{}
	}
	id, err := r.ArtistRepo.Add(ctx, newArtist)
	if err != nil {
		log.Errorf("Error while adding artist %v, %v\n", artist, err)
		return artist, err
	}

	newArtist.ID.Primary = id
	log.Debug("Added artist to database", newArtist)
	return newArtist, nil
}

func (r *EventRepository) UpdateArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error) {
	log.Debug("Request to update artist", artist)
	updateArtist := domain.CloneArtist(artist)
	err := r.ArtistRepo.Update(ctx, updateArtist)
	if err != nil {
		log.Errorf("Error while updating artist %v, %v\n", artist, err)
		return artist, err
	}
	log.Debug("Updated artist in database", updateArtist)
	return updateArtist, nil
}

func (r *EventRepository) DeleteArtist(ctx context.Context, id string) error {
	log.Debug("Request to delete artist", id)
	err := r.ArtistRepo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting artist %v, %v\n", id, err)
		return err
	}
	log.Debug("Deleted artist from database", id)
	return nil
}

func (r *EventRepository) ListArtists(ctx context.Context) ([]domain.Artist, error) {
	log.Debug("Request to list all artists")
	artists, err := r.ArtistRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all artists", err)
		return nil, err
	}
	return artists, nil
}

// Requires that all the artists and the venue already exist
func (r *EventRepository) AddEvent(ctx context.Context, event domain.Event) (domain.Event, error) {
	log.Debug("Request to add event", event)
	if !event.Populated() {
		log.Debug("Skipping adding event because required fields are missing", event)
		return event, errors.New("failed to create event due to empty fields")
	}
	newEvent := domain.CloneEvent(event)
	id, err := r.EventRepo.Add(ctx, newEvent)
	if err != nil {
		log.Errorf("Error while adding event %v, %v\n", event, err)
		return event, err
	}
	newEvent.ID.Primary = id
	log.Debug("Added event to database", newEvent)
	return newEvent, nil
}

func (r *EventRepository) DeleteEvent(ctx context.Context, id string) error {
	log.Debug("Request to delete event", id)
	err := r.EventRepo.Delete(ctx, id)
	if err != nil {
		log.Errorf("Error while deleting event %v, %v\n", id, err)
		return err
	}
	log.Debug("Deleted event from database", id)
	return nil
}

func (r *EventRepository) ListEvents(ctx context.Context) ([]domain.Event, error) {
	log.Debug("Request to list all events")
	events, err := r.EventRepo.FindAll(ctx)
	if err != nil {
		log.Error("Error while listing all events", err)
		return nil, err
	}
	return events, nil
}
