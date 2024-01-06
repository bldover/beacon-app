package svc

import (
	"context"
	"errors"
	"testing"
	"time"
)

type venueRepo struct{}
func (venueRepo) Add(context.Context, Venue) (string, error) {
    return "", nil
}
func (venueRepo) Delete(context.Context, Venue) error {
    return nil
}
func (venueRepo) Exists(context.Context, Venue) (bool, error) {
    return true, nil
}
func (venueRepo) FindAll(context.Context) (*[]Venue, error) {
    return nil, nil
}

type artistRepo struct{}
func (artistRepo) Add(context.Context, Artist) (string, error) {
    return "", nil
}
func (artistRepo) Delete(context.Context, Artist) error {
    return nil
}
func (artistRepo) Exists(context.Context, Artist) (bool, error) {
    return true, nil
}
func (artistRepo) FindAll(context.Context) (*[]Artist, error) {
    return nil, nil
}

type eventRepo struct{}
func (eventRepo) Add(context.Context, Event) (string, error) {
    return "", nil
}
func (eventRepo) Delete(context.Context, Event) error {
    return nil
}
func (eventRepo) Exists(context.Context, Event) (bool, error) {
    return true, nil
}
func (eventRepo) FindAll(context.Context) (*[]Event, error) {
    return nil, nil
}

func TestAddVenueValid(t *testing.T) {
	interactor := &EventInteractor{VenueRepo: venueRepo{}, ArtistRepo: artistRepo{}, EventRepo: eventRepo{}}
	venue := Venue{Name: "name", City: "city", State: "state", Country: "country", Outside: true}

	if err := interactor.AddVenue(context.Background(), venue); err != nil {
		t.Error("enexpected error:", err)
	}
}

func TestAddVenueInvalid(t *testing.T) {
	interactor := &EventInteractor{VenueRepo: venueRepo{}, ArtistRepo: artistRepo{}, EventRepo: eventRepo{}}
	venue := Venue{City: "city", State: "state", Country: "country", Outside: true}

	if err := interactor.AddVenue(context.Background(), venue); err == nil {
		t.Error("error expected")
	}
}

func TestAddArtistValid(t *testing.T) {
	interactor := &EventInteractor{VenueRepo: venueRepo{}, ArtistRepo: artistRepo{}, EventRepo: eventRepo{}}
	artist := Artist{Name: "name", Genre: "genre"}

	if err := interactor.AddArtist(context.Background(), artist); err != nil {
		t.Error("enexpected error:", err)
	}
}

func TestAddArtistInvalid(t *testing.T) {
	interactor := &EventInteractor{VenueRepo: venueRepo{}, ArtistRepo: artistRepo{}, EventRepo: eventRepo{}}
	artist := Artist{Genre: "genre"}

	if err := interactor.AddArtist(context.Background(), artist); err == nil {
		t.Error("error expected")
	}
}

func TestAddEventValid(t *testing.T) {
	interactor := &EventInteractor{VenueRepo: venueRepo{}, ArtistRepo: artistRepo{}, EventRepo: eventRepo{}}
	venue := Venue{Name: "name", City: "city", State: "state", Country: "country", Outside: true}
	artist := Artist{Name: "name", Genre: "genre"}
	event := Event{MainAct: artist, Openers: []Artist{artist}, Venue: venue, Date: time.Now()}

	if err := interactor.AddEvent(context.Background(), event); err != nil {
		t.Error("enexpected error:", err)
	}
}

func TestAddEventInvalid(t *testing.T) {
	interactor := &EventInteractor{VenueRepo: venueRepo{}, ArtistRepo: artistRepo{}, EventRepo: eventRepo{}}
	venue := Venue{Name: "name", City: "city", State: "state", Country: "country", Outside: true}
	artist := Artist{Name: "name", Genre: "genre"}
 	event := Event{Openers: []Artist{artist}, Venue: venue, Date: time.Now()}

	if err := interactor.AddEvent(context.Background(), event); err == nil {
		t.Error("error expected")
	}
}

type venueRepoErr struct{}
func (venueRepoErr) Add(context.Context, Venue) (string, error) {
    return "", errors.New("")
}
func (venueRepoErr) Delete(context.Context, Venue) error {
    return nil
}
func (venueRepoErr) Exists(context.Context, Venue) (bool, error) {
    return true, nil
}
func (venueRepoErr) FindAll(context.Context) (*[]Venue, error) {
    return nil, nil
}

type artistRepoErr struct{}
func (artistRepoErr) Add(context.Context, Artist) (string, error) {
    return "", errors.New("")
}
func (artistRepoErr) Delete(context.Context, Artist) error {
    return nil
}
func (artistRepoErr) Exists(context.Context, Artist) (bool, error) {
    return true, nil
}
func (artistRepoErr) FindAll(context.Context) (*[]Artist, error) {
    return nil, nil
}

type eventRepoErr struct{}
func (eventRepoErr) Add(context.Context, Event) (string, error) {
    return "", errors.New("")
}
func (eventRepoErr) Delete(context.Context, Event) error {
    return nil
}
func (eventRepoErr) Exists(context.Context, Event) (bool, error) {
    return true, nil
}
func (eventRepoErr) FindAll(context.Context) (*[]Event, error) {
    return nil, nil
}

func TestAddVenueErr(t *testing.T) {
	interactor := &EventInteractor{VenueRepo: venueRepoErr{}, ArtistRepo: artistRepoErr{}, EventRepo: eventRepoErr{}}
	venue := Venue{Name: "name", City: "city", State: "state", Country: "country", Outside: true}

	if err := interactor.AddVenue(context.Background(), venue); err == nil {
		t.Error("expected error")
	}
}

func TestAddArtistErr(t *testing.T) {
	interactor := &EventInteractor{VenueRepo: venueRepoErr{}, ArtistRepo: artistRepoErr{}, EventRepo: eventRepoErr{}}
	artist := Artist{Name: "name", Genre: "genre"}

	if err := interactor.AddArtist(context.Background(), artist); err == nil {
		t.Error("expected error")
	}
}

func TestAddEventError(t *testing.T) {
	interactor := &EventInteractor{VenueRepo: venueRepoErr{}, ArtistRepo: artistRepoErr{}, EventRepo: eventRepoErr{}}
	venue := Venue{Name: "name", City: "city", State: "state", Country: "country", Outside: true}
	artist := Artist{Name: "name", Genre: "genre"}
	event := Event{MainAct: artist, Openers: []Artist{artist}, Venue: venue, Date: time.Now()}

	if err := interactor.AddEvent(context.Background(), event); err == nil {
		t.Error("expected error")
	}
}
