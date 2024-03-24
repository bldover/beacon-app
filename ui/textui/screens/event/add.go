package event

import (
	"concert-manager/data"
	"concert-manager/ui/textui/input"
	"concert-manager/ui/textui/output"
	"concert-manager/ui/textui/screens"
	"concert-manager/util/format"
	"slices"
)

type eventAddCache interface {
	AddEvent(data.Event) error
}

type Adder struct {
	ArtistEditor  screens.Screen
	VenueEditor   screens.Screen
	OpenerRemover screens.Screen
	Cache         eventAddCache
	actions       map[data.EventType][]string
	eventType     data.EventType
	newEvent      data.Event
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

func NewAddScreen() *Adder {
	a := Adder{}
	a.actions = make(map[data.EventType][]string)
	a.actions[data.Future] = []string{"Edit Main Act", "Add Opener", "Remove Opener", "Edit Venue",
		"Edit Date", "Toggle Purchased", "Save Event", "Cancel"}
	a.actions[data.Past] = []string{"Edit Main Act", "Add Opener", "Remove Opener", "Edit Venue",
		"Edit Date", "Save Event", "Cancel"}
	return &a
}

func (a *Adder) AddContext(context screens.ScreenContext) {
	a.returnScreen = context.ReturnScreen
	props := context.Props
	if len(props) > 0 {
		a.newEvent = props[0].(data.Event)
	}

	// TODO: reconsider this since it requires changes here for new screens using this screen
	switch a.returnScreen.Title() {
	case "Concert History":
		a.eventType = data.Past
	case "Future Concerts", "All Upcoming Events":
		a.eventType = data.Future
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

func (a *Adder) NextScreen(i int) (screens.Screen, *screens.ScreenContext) {
	if a.eventType == data.Past && i >= togglePurchased {
		i += 1
	}

	switch i {
	case editMainAct:
		return a.ArtistEditor, screens.NewScreenContext(a, &a.newEvent.MainAct)
	case addOpener:
		if len(a.newEvent.Openers) >= maxOpeners {
			output.Displayln("Max number of openers is already reached!")
			return a, nil
		}
		a.newEvent.Openers = append(a.newEvent.Openers, data.Artist{})
		context := screens.NewScreenContext(a, &a.newEvent.Openers[len(a.newEvent.Openers)-1])
		return a.ArtistEditor, context
	case removeOpener:
		return a.OpenerRemover, screens.NewScreenContext(a, &a.newEvent.Openers)
	case editVenue:
		return a.VenueEditor, screens.NewScreenContext(a, &a.newEvent.Venue)
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
		if err := a.Cache.AddEvent(a.newEvent); err != nil {
			output.Displayf("Failed to save event: %v\n", err)
			return a, nil
		}
		return a.returnScreen, nil
	case cancelAddEvent:
		a.newEvent = data.Event{}
		return a.returnScreen, nil
	}
	return a, nil
}
