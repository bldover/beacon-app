package venue

import (
	"concert-manager/data"
	"concert-manager/ui/textui/input"
	"concert-manager/ui/textui/output"
	"concert-manager/ui/textui/screens"
)

type Search interface {
    FindFuzzyVenueMatchesByName(string) []data.Venue
}

type VenueEditor struct {
	Search       Search
	SelectScreen screens.Screen
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

func (e *VenueEditor) AddContext(context screens.ScreenContext) {
	if context.ContextType == screens.Selector {
		e.tempVenue = context.Props[0].(data.Venue)
		return
	}

	e.returnScreen = context.ReturnScreen
	props := context.Props
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

func (e *VenueEditor) NextScreen(i int) (screens.Screen, *screens.ScreenContext) {
	switch i {
	case search:
		name := input.PromptAndGetInput("venue name", input.NoValidation)
		matches := e.Search.FindFuzzyVenueMatchesByName(name)
		return e.SelectScreen, screens.NewScreenContext(e, matches)
	case setName:
		e.tempVenue.Name = input.PromptAndGetInput("venue name", input.NoValidation)
	case setCity:
		e.tempVenue.City = input.PromptAndGetInput("venue city", input.NoValidation)
	case setState:
		e.tempVenue.State = input.PromptAndGetInput("venue state", input.NoValidation)
	case save:
		if e.tempVenue.Populated() {
			*e.venue = e.tempVenue
			return e.returnScreen, nil
		} else {
			output.Displayln("Failed to save venue: all fields are required")
		}
	case cancelEdit:
		return e.returnScreen, nil
	}
	return e, nil
}
