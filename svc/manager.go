package svc

import (
	"concert-manager/data"
	"concert-manager/out"
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
	EventInteractor struct {
		VenueRepo  VenueRepo
		ArtistRepo ArtistRepo
		EventRepo  EventRepo
	}
)

func (interactor *EventInteractor) AddVenue(ctx Context, venue Venue) error {
	out.Debugf("Request to add venue %v", venue)
	if !venue.Populated() {
		out.Debugf("Skipping adding venue because required fields are missing %v", venue)
		return errors.New("failed to create venue due to empty fields")
	}

	_, err := interactor.VenueRepo.Add(ctx, venue)
	if err != nil {
		out.Errorf("Error while adding venue %v, %v", venue, err)
		return err
	}
	return nil
}

func (interactor *EventInteractor) DeleteVenue(ctx Context, venue Venue) error {
	out.Debugf("Request to delete venue %v", venue)
	err := interactor.VenueRepo.Delete(ctx, venue)
	if err != nil {
		out.Errorf("Error while deleting venue %v, %v", venue, err)
		return err
	}
	return nil
}

func (interactor *EventInteractor) ListVenues(ctx Context) (*[]Venue, error) {
	out.Debugln("Request to list all venues")
    venues, err := interactor.VenueRepo.FindAll(ctx)
	if err != nil {
		out.Errorf("Error while listing all venues, %v", err)
		return nil, err
	}
	return venues, nil
}

func (interactor *EventInteractor) AddArtist(ctx Context, artist Artist) error {
	out.Debugf("Request to add artist %v", artist)
	if !artist.Populated() {
		out.Debugf("Skipping adding artist because required fields are missing %v", artist)
		return errors.New("failed to create artist due to empty fields")
	}

	_, err := interactor.ArtistRepo.Add(ctx, artist)
	if err != nil {
		out.Errorf("Error while adding artist %v, %v", artist, err)
		return err
	}
	return nil
}

func (interactor *EventInteractor) DeleteArtist(ctx Context, artist Artist) error {
	out.Debugf("Request to delete artist %v", artist)
	err := interactor.ArtistRepo.Delete(ctx, artist)
	if err != nil {
		out.Errorf("Error while deleting artist %v, %v", artist, err)
		return err
	}
	return nil
}

func (interactor *EventInteractor) ListArtists(ctx Context) (*[]Artist, error) {
	out.Debugln("Request to list all artists")
    artists, err := interactor.ArtistRepo.FindAll(ctx)
	if err != nil {
		out.Errorf("Error while listing all artists, %v", err)
		return nil, err
	}
	return artists, nil
}

// Requires that all the artists and the venue already exist
func (interactor *EventInteractor) AddEvent(ctx Context, event Event) error {
	out.Debugf("Request to add event %v", event)
	if !event.Populated() {
		out.Debugf("Skipping adding event because required fields are missing %v", event)
		return errors.New("failed to create event due to empty fields")
	}

	_, err := interactor.EventRepo.Add(ctx, event)
	if err != nil {
		out.Errorf("Error while adding event %v, %v", event, err)
		return err
	}
	return nil
}

// Creates event and also venue and artists if needed
func (interactor *EventInteractor) AddEventRecursive(ctx Context, event Event) error {
	out.Debugf("Request to recursively add event %v", event)
	if !event.Populated() {
		out.Debugf("Skipping adding event because required fields are missing %v", event)
		return errors.New("failed to create event due to empty fields")
	}

	if err := interactor.AddVenue(ctx, event.Venue); err != nil {
		out.Errorf("Failed to add venue %v while recursively adding event, %v", event.Venue, err)
		return err
	}
	if event.MainAct.Populated() {
		if err := interactor.AddArtist(ctx, event.MainAct); err != nil {
			out.Errorf("Failed to add artist %v while recursively adding event, %v", event.MainAct, err)
			return err
		}
	}
	for _, opener := range event.Openers {
		if err := interactor.AddArtist(ctx, opener); err != nil {
			out.Errorf("Failed to add artist %v while recursively adding event, %v", event.MainAct, err)
			return err
		}
	}

	_, err := interactor.EventRepo.Add(ctx, event)
	if err != nil {
		out.Errorf("Error while recursively adding event %v, %v", event, err)
		return err
	}
	return nil
}

func (interactor *EventInteractor) DeleteEvent(ctx Context, event Event) error {
	out.Debugf("Request to delete event %v", event)
    err := interactor.EventRepo.Delete(ctx, event)
	if err != nil {
		out.Errorf("Error while deleting event %v, %v", event, err)
		return err
	}
	return nil
}

func (interactor *EventInteractor) ListEvents(ctx Context) (*[]Event, error) {
	out.Debugln("Request to list all events")
    events, err := interactor.EventRepo.FindAll(ctx)
	if err != nil {
		out.Errorf("Error while listing all events, %v", err)
		return nil, err
	}
	return events, nil
}
