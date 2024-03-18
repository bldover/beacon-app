package event

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/ui/terminal/input"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"concert-manager/ui/terminal/screens/format"
	"context"
	"slices"
)

type eventDbAdder interface {
	AddEventRecursive(context.Context, data.Event) error
}

type artistEditScreen interface {
	screens.Screen
	AddArtistContext(*data.Artist)
}

type venueEditScreen interface {
	screens.Screen
	AddVenueContext(*data.Venue)
}

type openerRemoveScreen interface {
	screens.Screen
	AddOpenerContext(*[]data.Artist)
}

type Adder struct {
	Events        *[]data.Event
	NewEvent      data.Event
	Database      eventDbAdder
	ArtistEditor  artistEditScreen
	VenueEditor   venueEditScreen
	OpenerRemover openerRemoveScreen
	Viewer        screens.Screen
	actions       []string
	futureEvents  bool
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

func NewAddScreen(futureEvents bool) *Adder {
	a := Adder{}
	a.futureEvents = futureEvents
	if futureEvents {
		a.actions = []string{"Edit Main Act", "Add Opener", "Remove Opener", "Edit Venue", "Edit Date", "Toggle Purchased", "Save Event", "Cancel"}
	} else {
		a.actions = []string{"Edit Main Act", "Add Opener", "Remove Opener", "Edit Venue", "Edit Date", "Save Event", "Cancel"}
	}
	return &a
}

func (a *Adder) AddContext(details data.EventDetails) {
	a.NewEvent = details.Event
}

func (a Adder) Title() string {
	return "Add Concert"
}

func (a *Adder) DisplayData() {
	for i, op := range a.NewEvent.Openers {
		if !op.Populated() {
			// at most one could have been added since the last check
			a.NewEvent.Openers = slices.Delete(a.NewEvent.Openers, i, i+1)
			break
		}
	}
	output.Displayln(format.FormatEventExpanded(a.NewEvent, a.futureEvents))
}

func (a Adder) Actions() []string {
	return a.actions
}

func (a *Adder) NextScreen(i int) screens.Screen {
	if !a.futureEvents && i >= togglePurchased {
		i += 1
	}

	switch i {
	case editMainAct:
		a.ArtistEditor.AddArtistContext(&a.NewEvent.MainAct)
		return a.ArtistEditor

	case addOpener:
		if len(a.NewEvent.Openers) >= maxOpeners {
			output.Displayln("Max number of openers is already reached!")
			return a
		}
		a.NewEvent.Openers = append(a.NewEvent.Openers, data.Artist{})
		log.Debugf("add event - op: %p", &a.NewEvent.Openers)
		log.Debugf("add event - op: %+v", a.NewEvent.Openers)
		log.Debugf("add event - op: %p", &a.NewEvent.Openers[0])
		log.Debugf("add event - op: %+v", a.NewEvent.Openers[0])
		a.ArtistEditor.AddArtistContext(&a.NewEvent.Openers[len(a.NewEvent.Openers)-1])
		return a.ArtistEditor

	case removeOpener:
		a.OpenerRemover.AddOpenerContext(&a.NewEvent.Openers)
		return a.OpenerRemover

	case editVenue:
		a.VenueEditor.AddVenueContext(&a.NewEvent.Venue)
		return a.VenueEditor

	case editDate:
		if a.futureEvents {
			a.NewEvent.Date = input.PromptAndGetInput("event date (mm/dd/yyyy)", input.FutureDateValidation)
		} else {
			a.NewEvent.Date = input.PromptAndGetInput("event date (mm/dd/yyyy)", input.PastDateValidation)
		}

	case togglePurchased:
		a.NewEvent.Purchased = !a.NewEvent.Purchased

	case saveEvent:
		if !a.futureEvents {
			a.NewEvent.Purchased = true
		}
		if err := a.Database.AddEventRecursive(context.Background(), a.NewEvent); err != nil {
			output.Displayf("Failed to save event: %v\n", err)
			return a
		}
		*a.Events = append(*a.Events, a.NewEvent)
		return a.Viewer

	case cancelAddEvent:
		a.NewEvent = data.Event{}
		return a.Viewer
	}

	return a
}
