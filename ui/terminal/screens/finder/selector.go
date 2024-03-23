package finder

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"concert-manager/ui/terminal/screens/format"
)

type Selector struct {
	eventDetails   []data.EventDetails
	eventAddScreen screens.ContextScreen
	returnScreen   screens.Screen
}

func NewSelectorScreen(eventAddScreen screens.ContextScreen) *Selector {
	selector := Selector{}
	selector.eventAddScreen = eventAddScreen
	return &selector
}

func (s *Selector) AddContext(returnScreen screens.Screen, props ...any) {
	s.returnScreen = returnScreen
	s.eventDetails = props[0].([]data.EventDetails)
}

func (s Selector) Title() string {
	return "Select Concert to Save"
}

func (s Selector) DisplayData() {
	if len(s.eventDetails) == 0 {
		output.Displayln("No concerts found")
	}
}

func (s Selector) Actions() []string {
	actions := []string{}
	pageEvents := s.eventDetails
	actions = append(actions, format.FormatEventDetailsShort(pageEvents)...)
	actions = append(actions, "Cancel")
	return actions
}

func (s *Selector) NextScreen(i int) screens.Screen {
	if i != len(s.eventDetails)+1 {
		eventIdx := i - 1
		s.eventAddScreen.AddContext(s.returnScreen, data.EventType(data.Future), s.eventDetails[eventIdx].Event)
		return s.eventAddScreen
	}
	return s.returnScreen
}
