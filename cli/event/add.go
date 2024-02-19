package event

import (
	"concert-manager/cli"
	"concert-manager/cli/format"
	"concert-manager/data"
	"concert-manager/out"
	"context"
	"slices"
)

type Adder struct {
	Events *[]data.Event
	NewEvent data.Event
	EventAdder EventAdder
	ArtistEditor ArtistEditScreen
	VenueEditor VenueEditScreen
	OpenerRemover OpenerRemoveScreen
	Viewer cli.Screen
	actions []string
	futureEvents bool
}

type EventAdder interface {
    AddEvent(context.Context, data.Event) error
}

type ArtistEditScreen interface {
	cli.Screen
	AddArtistContext(*data.Artist)
}

type VenueEditScreen interface {
    cli.Screen
	AddVenueContext(*data.Venue)
}

type OpenerRemoveScreen interface {
    cli.Screen
	AddOpenerContext(*[]data.Artist)
}

const (
	mainAct = iota + 1
	addOpener
	removeOpener
	venue
	date
	togglePurchased
	save
	cancel
)

const maxOpeners = 10

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

func (a Adder) Title() string {
    return "Add Concert"
}

func (a *Adder) DisplayData() {
	for i, op := range a.NewEvent.Openers {
		if !op.Populated() {
			// at most one could have been added since the last check
			a.NewEvent.Openers = slices.Delete(a.NewEvent.Openers, i, i + 1)
			break
		}
	}
    out.Displayln(format.FormatEventExpanded(a.NewEvent, a.futureEvents))
}

func (a Adder) Actions() []string {
	return a.actions
}

func (a *Adder) NextScreen(i int) cli.Screen {
	if !a.futureEvents && i >= togglePurchased {
		i += 1
	}

	switch i {
    case mainAct:
		a.ArtistEditor.AddArtistContext(&a.NewEvent.MainAct)
		return a.ArtistEditor

	case addOpener:
		if len(a.NewEvent.Openers) >= maxOpeners {
			out.Displayln("Max number of openers is already reached!")
			return a
		}
		a.NewEvent.Openers = append(a.NewEvent.Openers, data.Artist{})
		a.ArtistEditor.AddArtistContext(&a.NewEvent.Openers[len(a.NewEvent.Openers) - 1])
		return a.ArtistEditor

	case removeOpener:
		a.OpenerRemover.AddOpenerContext(&a.NewEvent.Openers)
		return a.OpenerRemover

	case venue:
		a.VenueEditor.AddVenueContext(&a.NewEvent.Venue)
		return a.VenueEditor

	case date:
		if a.futureEvents {
			a.NewEvent.Date = cli.PromptAndGetInput("event date (mm/dd/yyyy)", data.ValidFutureDate)
		} else {
			a.NewEvent.Date = cli.PromptAndGetInput("event date (mm/dd/yyyy)", data.ValidPastDate)
		}

	case togglePurchased:
		a.NewEvent.Purchased = !a.NewEvent.Purchased

	case save:
		if !a.futureEvents {
			a.NewEvent.Purchased = true
		}
		if err := a.EventAdder.AddEvent(context.Background(), a.NewEvent); err != nil {
			out.Displayf("Failed to save event: %v\n", err)
			return a
		}
		*a.Events = append(*a.Events, a.NewEvent)
		return a.Viewer

	case cancel:
		a.NewEvent = data.Event{}
		return a.Viewer
	}

	return a
}
