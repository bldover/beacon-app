package data

import (
	"testing"
	"time"
)

func TestVenuePopulatedValid(t *testing.T) {
    v := Venue{Name: "name", City: "city", State: "state"}
	if !v.Populated() {
		t.Error("populated venue marked as not populated")
	}
}

func TestVenuePopulatedInvalid(t *testing.T) {
    v := Venue{City: "city", State: "state"}
	if v.Populated() {
		t.Error("not populated venue marked as populated")
	}
}

func TestArtistPopulatedValid(t *testing.T) {
    a := Artist{Name: "Name", Genre: "Genre"}
	if !a.Populated() {
		t.Error("populated artist marked as not populated")
	}
}

func TestArtistPopulatedInvalid(t *testing.T) {
    a := Artist{Genre: "Genre"}
	if a.Populated() {
		t.Error("not populated artist marked as populated")
	}
}

func TestEventPopulatedValidOnlyMainAct(t *testing.T) {
    v := Venue{Name: "name", City: "city", State: "state"}
	a := Artist{Name: "Name", Genre: "Genre"}
	e := Event{MainAct: a, Venue: v, Date: time.Now()}
	if !e.Populated() {
		t.Error("populated event marked as not populated")
	}
}

func TestEventPopulatedValidNoMainAct(t *testing.T) {
    v := Venue{Name: "name", City: "city", State: "state"}
	a := Artist{Name: "Name", Genre: "Genre"}
	e := Event{Openers: []Artist{a}, Venue: v, Date: time.Now()}
	if !e.Populated() {
		t.Error("populated event marked as not populated")
	}
}

func TestEventPopulatedInvalidNoArtists(t *testing.T) {
    v := Venue{Name: "name", City: "city", State: "state"}
	e := Event{Venue: v, Date: time.Now()}
	if e.Populated() {
		t.Error("event with no artists marked as populated")
	}
}

func TestEventPopulatedInvalidArtist(t *testing.T) {
    v := Venue{Name: "name", City: "city", State: "state"}
	a := Artist{Genre: "Genre"}
	e := Event{Openers: []Artist{a}, Venue: v, Date: time.Now()}
	if e.Populated() {
		t.Error("event with invalid artist marked as populated")
	}
}

func TestEventPopulatedInvalidVenue(t *testing.T) {
    v := Venue{City: "city", State: "state"}
	a := Artist{Name: "Name", Genre: "Genre"}
	e := Event{Openers: []Artist{a}, Venue: v, Date: time.Now()}
	if e.Populated() {
		t.Error("event with invalid artist marked as populated")
	}
}

func TestEventPopulatedInvalidDate(t *testing.T) {
    v := Venue{Name: "name", City: "city", State: "state"}
	a := Artist{Name: "Name", Genre: "Genre"}
	e := Event{MainAct: a, Venue: v}
	if e.Populated() {
		t.Error("event with no date marked as populated")
	}
}
