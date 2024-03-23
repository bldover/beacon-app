package event

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/input"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"concert-manager/ui/terminal/screens/format"
	"slices"
)

type eventAddCache interface {
	AddEvent(data.Event) error
}

type Adder struct {
	actions       map[data.EventType][]string
	eventType     data.EventType
	cache         eventAddCache
	newEvent      data.Event
	artistEditor  screens.ContextScreen
	venueEditor   screens.ContextScreen
	openerRemover screens.ContextScreen
	returnScreen  screens.Screen
}

const (
	editMainAct = iota + 1
	addOpener
	removeOpener
	editVenue
	editDate
	togglePurchased
	saveEvent
	cancelAddEvent
)

const maxOpeners = 20

func NewAddScreen(artistEditor, venueEditor, openerRemover screens.ContextScreen, cache eventAddCache) *Adder {
	a := Adder{}
	a.actions = make(map[data.EventType][]string)
	a.actions[data.Future] = []string{"Edit Main Act", "Add Opener", "Remove Opener", "Edit Venue",
		"Edit Date", "Toggle Purchased", "Save Event", "Cancel"}
	a.actions[data.Past] = []string{"Edit Main Act", "Add Opener", "Remove Opener", "Edit Venue",
		"Edit Date", "Save Event", "Cancel"}
	a.artistEditor = artistEditor
	a.venueEditor = venueEditor
	a.openerRemover = openerRemover
	a.cache = cache
	return &a
}

func (a *Adder) AddContext(returnScreen screens.Screen, props ...any) {
	a.returnScreen = returnScreen
	a.eventType = props[0].(data.EventType)
	if len(props) > 1 {
		a.newEvent = props[1].(data.Event)
	}
}

func (a Adder) Title() string {
	return "Add Concert"
}

func (a *Adder) DisplayData() {
	for i, op := range a.newEvent.Openers {
		if !op.Populated() {
			// at most one could have been added since the last check
			a.newEvent.Openers = slices.Delete(a.newEvent.Openers, i, i+1)
			break
		}
	}

	isFuture := a.eventType == data.Future
	output.Displayln(format.FormatEventExpanded(a.newEvent, isFuture))
}

func (a Adder) Actions() []string {
	return a.actions[a.eventType]
}

func (a *Adder) NextScreen(i int) screens.Screen {
	if a.eventType == data.Past && i >= togglePurchased {
		i += 1
	}

	switch i {
	case editMainAct:
		a.artistEditor.AddContext(a, &a.newEvent.MainAct)
		return a.artistEditor
	case addOpener:
		if len(a.newEvent.Openers) >= maxOpeners {
			output.Displayln("Max number of openers is already reached!")
			return a
		}
		a.newEvent.Openers = append(a.newEvent.Openers, data.Artist{})
		a.artistEditor.AddContext(a, &a.newEvent.Openers[len(a.newEvent.Openers)-1])
		return a.artistEditor
	case removeOpener:
		a.openerRemover.AddContext(a, &a.newEvent.Openers)
		return a.openerRemover
	case editVenue:
		a.venueEditor.AddContext(a, &a.newEvent.Venue)
		return a.venueEditor
	case editDate:
		if a.eventType == data.Future {
			a.newEvent.Date = input.PromptAndGetInput("event date (mm/dd/yyyy)", input.FutureDateValidation)
		} else {
			a.newEvent.Date = input.PromptAndGetInput("event date (mm/dd/yyyy)", input.PastDateValidation)
		}
	case togglePurchased:
		a.newEvent.Purchased = !a.newEvent.Purchased
	case saveEvent:
		if a.eventType == data.Past {
			a.newEvent.Purchased = true
		}
		if err := a.cache.AddEvent(a.newEvent); err != nil {
			output.Displayf("Failed to save event: %v\n", err)
			return a
		}
		return a.returnScreen
	case cancelAddEvent:
		a.newEvent = data.Event{}
		return a.returnScreen
	}
	return a
}
