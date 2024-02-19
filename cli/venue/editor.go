package venue

import (
	"concert-manager/cli"
	"concert-manager/data"
	"concert-manager/out"
	"context"
)

type Editor struct {
    venue *data.Venue
	tempVenue data.Venue
	Venues *[]data.Venue
	VenueAdder VenueAdder
	AddEventScreen cli.Screen
	actions []string
}

const (
	search = iota + 1
	setName
	setCity
	setState
	save
	cancel
)

type VenueAdder interface {
    AddVenue(context.Context, data.Venue) error
}

func NewEditScreen() *Editor {
	e := Editor{}
	e.actions = []string{"Search Venues", "Set Name", "Set City", "Set State", "Save Venue", "Cancel"}
    return &e
}

func (e *Editor) AddVenueContext(venue *data.Venue) {
    e.venue = venue
	e.tempVenue.Name = venue.Name
	e.tempVenue.City = venue.City
	e.tempVenue.State = venue.State
}

func (e Editor) Title() string {
    return "Edit Venue"
}

func (e Editor) DisplayData() {
    out.Displayf("%+v\n", e.tempVenue)
}

func (e Editor) Actions() []string {
    return e.actions
}

func (e *Editor) NextScreen(i int) cli.Screen {
	switch i {
    case search:
		e.handleSearch()
	case setName:
		e.tempVenue.Name = cli.PromptAndGetInput("venue name", cli.NoValidation)
	case setCity:
		e.tempVenue.City = cli.PromptAndGetInput("venue city", cli.NoValidation)
	case setState:
		e.tempVenue.State = cli.PromptAndGetInput("venue state", cli.NoValidation)
	case save:
		if err := e.VenueAdder.AddVenue(context.Background(), e.tempVenue); err != nil {
			out.Displayf("Failed to save venue: %v\n", err)
		} else {
			*e.venue = e.tempVenue
			return e.AddEventScreen
		}
	case cancel:
		return e.AddEventScreen
	}
	return e
}

func (e *Editor) handleSearch() {
	out.Displayln("Not yet implemented!")
	//TODO: add search functionality
}
