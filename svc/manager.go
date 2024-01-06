package svc

import (
	"concert-manager/data"
	"context"
	"errors"
	"log"
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
	if !venue.Populated() {
		return errors.New("failed to create venue due to empty fields")
	}

	v, err := interactor.VenueRepo.Add(ctx, venue)
	if err != nil {
		return err
	}
	log.Printf("Created or found venue %v", v)
	return nil
}

func (interactor *EventInteractor) DeleteVenue(ctx Context, venue Venue) error {
	return interactor.VenueRepo.Delete(ctx, venue)
}

func (interactor *EventInteractor) ListVenues(ctx Context) (*[]Venue, error) {
    return interactor.VenueRepo.FindAll(ctx)
}

func (interactor *EventInteractor) AddArtist(ctx Context, artist Artist) error {
	if !artist.Populated() {
		return errors.New("failed to create artist due to empty fields")
	}

	a, err := interactor.ArtistRepo.Add(ctx, artist)
	if err != nil {
		return err
	}
	log.Printf("Created or found artist %v", a)
	return nil
}

func (interactor *EventInteractor) DeleteArtist(ctx Context, artist Artist) error {
	return interactor.ArtistRepo.Delete(ctx, artist)
}

func (interactor *EventInteractor) ListArtists(ctx Context) (*[]Artist, error) {
    return interactor.ArtistRepo.FindAll(ctx)
}

// Requires that all the artists and the venue already exist
func (interactor *EventInteractor) AddEvent(ctx Context, event Event) error {
	if !event.Populated() {
		return errors.New("failed to create event due to empty fields")
	}

	e, err := interactor.EventRepo.Add(ctx, event)
	if err != nil {
		return err
	}
	log.Printf("Created or found event %v", e)
	return nil
}

// Creates event and also venue and artists if needed
func (interactor *EventInteractor) AddEventRecursive(ctx Context, event Event) error {
	if !event.Populated() {
		return errors.New("failed to create event due to empty fields")
	}

	if err := interactor.AddVenue(ctx, event.Venue); err != nil {
		return err
	}
	if event.MainAct.Populated() {
		if err := interactor.AddArtist(ctx, event.MainAct); err != nil {
			return err
		}
	}
	for _, opener := range event.Openers {
		if err := interactor.AddArtist(ctx, opener); err != nil {
			return err
		}
	}

	e, err := interactor.EventRepo.Add(ctx, event)
	if err != nil {
		return err
	}
	log.Printf("Created or found event %v", e)
	return nil
}

func (interactor *EventInteractor) DeleteEvent(ctx Context, event Event) (error) {
    return interactor.EventRepo.Delete(ctx, event)
}

func (interactor *EventInteractor) ListEvents(ctx Context) (*[]Event, error) {
    return interactor.EventRepo.FindAll(ctx)
}
