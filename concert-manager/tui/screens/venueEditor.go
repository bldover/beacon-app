package screens

import (
	"concert-manager/domain"
	"concert-manager/search"
	"concert-manager/tui/input"
	"concert-manager/tui/output"
)

type venueCache interface {
	GetVenues() []domain.Venue
}

type VenueEditor struct {
	VenueCache   venueCache
	ReturnScreen Screen
	actions      []string
	venue        *domain.Venue
	tempVenue    domain.Venue
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

func (e *VenueEditor) SetVenue(venue *domain.Venue) {
	e.venue = venue
	e.tempVenue.Name = e.venue.Name
	e.tempVenue.City = e.venue.City
	e.tempVenue.State = e.venue.State
}

func (e VenueEditor) Title() string {
	return "Edit Venue"
}

func (e VenueEditor) DisplayData() {
	output.Displayf("%+v\n\n", output.FormatVenueExpanded(e.tempVenue))
}

func (e VenueEditor) Actions() []string {
	return e.actions
}

func (e *VenueEditor) NextScreen(i int) Screen {
	switch i {
	case searchVenue:
		name := input.PromptAndGetInput("venue name to search", input.NoValidation)
		matches := search.SearchVenues(name, e.VenueCache.GetVenues(), pageSize, search.LenientTolerance)
		selectScreen := &Selector[domain.Venue]{
			ScreenTitle: "Select Venue",
			Next:        e.ReturnScreen,
			Options:     matches,
			HandleSelect: func(v domain.Venue) {
				*e.venue = v
			},
			Formatter: output.FormatVenues,
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
