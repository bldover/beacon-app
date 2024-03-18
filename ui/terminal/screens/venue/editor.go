package venue

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/input"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
)

type VenueEditor struct {
    venue *data.Venue
	tempVenue data.Venue
	Venues *[]data.Venue
	AddEventScreen screens.Screen
	actions []string
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
		if e.venue.Populated() {
			*e.venue = e.tempVenue
			return e.AddEventScreen
		} else {
			output.Displayln("Failed to save venue: all fields are required")
		}
	case cancelEdit:
		return e.AddEventScreen
	}
	return e
}

func (e *VenueEditor) handleVenueSearch() {
	output.Displayln("Not yet implemented!")
	//TODO: add search functionality
}
