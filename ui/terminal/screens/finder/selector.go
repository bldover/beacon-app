package finder

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"concert-manager/ui/terminal/screens/format"
)

type eventAdderScreen interface {
	screens.Screen
	AddContext(data.EventDetails)
}

type Selector struct {
	EventAddScreen eventAdderScreen
	ViewerScreen   screens.Screen
	eventDetails   []data.EventDetails
}

func NewSelectorScreen() *Selector {
	return &Selector{}
}

func (s *Selector) AddContext(eventDetails []data.EventDetails) {
	s.eventDetails = eventDetails
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
		s.EventAddScreen.AddContext(s.eventDetails[eventIdx])
		return s.EventAddScreen
	}
	return s.ViewerScreen
}
