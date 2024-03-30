package screens

import (
	"concert-manager/data"
	"concert-manager/ui/textui/input"
	"concert-manager/ui/textui/output"
	"concert-manager/util"
)

type venueSearch interface {
	FindFuzzyVenueMatchesByName(string) []data.Venue
	WithMaxCount(int)
}

type VenueEditor struct {
	Search       venueSearch
	ReturnScreen Screen
	actions      []string
	venue        *data.Venue
	tempVenue    data.Venue
}

const (
	searchVenue = iota + 1
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

func (e *VenueEditor) SetVenue(venue *data.Venue) {
	e.venue = venue
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

func (e *VenueEditor) NextScreen(i int) Screen {
	switch i {
	case searchVenue:
		name := input.PromptAndGetInput("venue name to search", input.NoValidation)
		e.Search.WithMaxCount(pageSize)
		matches := e.Search.FindFuzzyVenueMatchesByName(name)
		selectScreen := &Selector[data.Venue]{
			ScreenTitle: "Select Venue",
			Next:        e.ReturnScreen,
			Options:     matches,
			HandleSelect: func(v data.Venue) {
				*e.venue = v
			},
			Formatter: util.FormatVenue,
		}
		return selectScreen
	case setVenueName:
		e.tempVenue.Name = input.PromptAndGetInput("venue name", input.NoValidation)
	case setVenueCity:
		e.tempVenue.City = input.PromptAndGetInput("venue city", input.NoValidation)
	case setVenueState:
		e.tempVenue.State = input.PromptAndGetInput("venue state", input.NoValidation)
	case saveVenue:
		if e.tempVenue.Populated() {
			*e.venue = e.tempVenue
			return nil
		} else {
			output.Displayln("Failed to save venue: all fields are required")
		}
	case cancelVenueEdit:
		return nil
	}
	return e
}
