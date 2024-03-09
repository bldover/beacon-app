package screens

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/input"
	"concert-manager/ui/terminal/output"
	"context"
)

type venueAdder interface {
    AddVenue(context.Context, data.Venue) error
}

type VenueEditor struct {
    venue *data.Venue
	tempVenue data.Venue
	Venues *[]data.Venue
	VenueAdder venueAdder
	AddEventScreen Screen
	actions []string
}

const (
	venueSearch = iota + 1
	setVenueName
	setVenueCity
	setVenueState
	saveVenue
	cancelVenueEdit
)

func NewVenueEditScreen() *VenueEditor {
	e := VenueEditor{}
	e.actions = []string{"Search Venues", "Set Name", "Set City", "Set State", "Save Venue", "Cancel"}
    return &e
}

func (e *VenueEditor) AddVenueContext(venue *data.Venue) {
    e.venue = venue
	e.tempVenue.Name = venue.Name
	e.tempVenue.City = venue.City
	e.tempVenue.State = venue.State
}

func (e VenueEditor) Title() string {
    return "Edit Venue"
}

func (e VenueEditor) DisplayData() {
    output.Displayf("%+v\n", e.tempVenue)
}

func (e VenueEditor) Actions() []string {
    return e.actions
}

func (e *VenueEditor) NextScreen(i int) Screen {
	switch i {
    case venueSearch:
		e.handleVenueSearch()
	case setName:
		e.tempVenue.Name = input.PromptAndGetInput("venue name", input.NoValidation)
	case setVenueCity:
		e.tempVenue.City = input.PromptAndGetInput("venue city", input.NoValidation)
	case setVenueState:
		e.tempVenue.State = input.PromptAndGetInput("venue state", input.NoValidation)
	case saveVenue:
		if err := e.VenueAdder.AddVenue(context.Background(), e.tempVenue); err != nil {
			output.Displayf("Failed to save venue: %v\n", err)
		} else {
			*e.venue = e.tempVenue
			return e.AddEventScreen
		}
	case cancelVenueEdit:
		return e.AddEventScreen
	}
	return e
}

func (e *VenueEditor) handleVenueSearch() {
	output.Displayln("Not yet implemented!")
	//TODO: add search functionality
}
