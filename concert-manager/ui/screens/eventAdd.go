package screens

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/ui/input"
	"concert-manager/ui/output"
	"concert-manager/util"
	"slices"
)

type eventAddCache interface {
	AddSavedEvent(data.Event) (*data.Event, error)
}

type artistEditor interface {
	Screen
	SetArtist(*data.Artist)
}

type venueEditor interface {
	Screen
	SetVenue(*data.Venue)
}

type EventAdder struct {
	newEvent         data.Event
	dateType         dateType
	ArtistEditor     artistEditor
	VenueEditor      venueEditor
	Cache            eventAddCache
	actions          []string
	beforeSaveAction extraAction
}

type extraAction func() error

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

type dateType int

const (
	future = iota
	past
)

const maxOpeners = 20

func NewEventAddScreen() *EventAdder {
	a := EventAdder{}
	a.actions = []string{
		"Edit Main Act",
		"Add Opener",
		"Remove Opener",
		"Edit Venue",
		"Edit Date",
		"Toggle Purchased",
		"Save Event",
		"Cancel",
	}
	return &a
}

func (a *EventAdder) WithBeforeSaveAction(action extraAction) {
	a.beforeSaveAction = action
}

func (a EventAdder) Title() string {
	return "Add Concert"
}

func (a *EventAdder) DisplayData() {
	for i, op := range a.newEvent.Openers {
		if !op.Populated() {
			// at most one could have been added since the last check
			a.newEvent.Openers = slices.Delete(a.newEvent.Openers, i, i+1)
			break
		}
	}

	output.Displayln(util.FormatEventExpanded(a.newEvent))
}

func (a EventAdder) Actions() []string {
	return a.actions
}

func (a *EventAdder) NextScreen(i int) Screen {
	switch i {
	case editMainAct:
		a.ArtistEditor.SetArtist(a.newEvent.MainAct)
		return a.ArtistEditor
	case addOpener:
		if len(a.newEvent.Openers) >= maxOpeners {
			output.Displayln("Max number of openers is already reached!")
			return a
		}
		a.newEvent.Openers = append(a.newEvent.Openers, data.Artist{})
		a.ArtistEditor.SetArtist(&a.newEvent.Openers[len(a.newEvent.Openers)-1])
		return a.ArtistEditor
	case removeOpener:
		selectScreen := &Selector[data.Artist]{
			ScreenTitle: "Remove Opener",
			Next:        a,
			Options:     a.newEvent.Openers,
			HandleSelect: func(artist data.Artist) {
				a.newEvent.Openers = slices.DeleteFunc(a.newEvent.Openers, artist.EqualsFields)
			},
			Formatter: util.FormatArtists,
		}
		return selectScreen
	case editVenue:
		a.VenueEditor.SetVenue(&a.newEvent.Venue)
		return a.VenueEditor
	case editDate:
		a.newEvent.Date = input.PromptAndGetInput("event date (mm/dd/yyyy)", input.DateValidation)
		if util.FutureDate(a.newEvent.Date) {
			a.dateType = future
		} else {
			a.dateType = past
			a.newEvent.Purchased = true
		}
	case togglePurchased:
		if a.dateType == past {
			output.Displayln("Past events must be purchased")
		} else {
			a.newEvent.Purchased = !a.newEvent.Purchased
		}
	case saveEvent:
		if a.beforeSaveAction != nil {
			if err := a.beforeSaveAction(); err != nil {
				log.Error("Before save action failed:", err)
				output.Displayf("Failed to save event: %v\n", err)
				return a
			}
		}
		if _, err := a.Cache.AddSavedEvent(a.newEvent); err != nil {
			output.Displayf("Failed to save event: %v\n", err)
			return a
		}
		fallthrough
	case cancelAddEvent:
		a.newEvent = data.Event{}
		a.dateType = future
		a.beforeSaveAction = nil
		return nil
	}
	return a
}
