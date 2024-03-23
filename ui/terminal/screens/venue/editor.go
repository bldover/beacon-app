package venue

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/input"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
)

type venueAddCache interface {
	AddVenue(data.Venue) error
}

type VenueEditor struct {
	actions      []string
	venue        *data.Venue
	tempVenue    data.Venue
	returnScreen screens.Screen
}

const (
	search = iota + 1
	setName
	setCity
	setState
	save
	cancelEdit
)

func NewEditScreen() *VenueEditor {
	e := VenueEditor{}
	e.actions = []string{"Search Venues", "Set Name", "Set City", "Set State", "Save Venue", "Cancel"}
	return &e
}

func (e *VenueEditor) AddContext(returnScreen screens.Screen, props ...any) {
	e.returnScreen = returnScreen
	e.venue = props[0].(*data.Venue)
	e.tempVenue.Name = e.venue.Name
	e.tempVenue.City = e.venue.City
	e.tempVenue.State = e.venue.State
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

func (e *VenueEditor) NextScreen(i int) screens.Screen {
	switch i {
	case search:
		e.handleVenueSearch()
	case setName:
		e.tempVenue.Name = input.PromptAndGetInput("venue name", input.NoValidation)
	case setCity:
		e.tempVenue.City = input.PromptAndGetInput("venue city", input.NoValidation)
	case setState:
		e.tempVenue.State = input.PromptAndGetInput("venue state", input.NoValidation)
	case save:
		if e.tempVenue.Populated() {
			*e.venue = e.tempVenue
			return e.returnScreen
		} else {
			output.Displayln("Failed to save venue: all fields are required")
		}
	case cancelEdit:
		return e.returnScreen
	}
	return e
}

func (e *VenueEditor) handleVenueSearch() {
	output.Displayln("Not yet implemented!")
	//TODO: add search functionality
}
