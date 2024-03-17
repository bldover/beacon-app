package event

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"concert-manager/ui/terminal/screens/format"
	"context"
	"math"
	"slices"
)

type eventDbDeleter interface {
	DeleteEvent(context.Context, data.Event) error
}

type EventDeleter struct {
	Events       *[]data.Event
	Database     eventDbDeleter
	Viewer       screens.Screen
	startIdx     int
	displayCount int
}

func NewDeleteScreen() *EventDeleter {
	return &EventDeleter{}
}

func (d *EventDeleter) AddDeleteContext(startIdx int, displayCount int) {
	d.startIdx = startIdx
	remainingEvents := len(*d.Events) - startIdx
	d.displayCount = int(math.Min(float64(displayCount), float64(remainingEvents)))
}

func (d EventDeleter) Title() string {
	return "Delete Concert"
}

func (d EventDeleter) DisplayData() {
	if len(*d.Events) == 0 {
		output.Displayln("No concerts found")
	}
}

func (d EventDeleter) Actions() []string {
	actions := []string{}
	pageEvents := []data.Event{}
	var maxNameLen int
	for i := d.startIdx; i < d.startIdx+d.displayCount; i++ {
		event := (*d.Events)[i]
		pageEvents = append(pageEvents, event)
		var artist string
		if event.MainAct.Populated() {
			artist = event.MainAct.Name
		} else {
			artist = event.Openers[0].Name
		}
		maxNameLen = int(math.Max(float64(maxNameLen), float64(len(artist))))
	}

	for _, event := range pageEvents {
		actions = append(actions, format.FormatEventShort(event, maxNameLen))
	}

	actions = append(actions, "Cancel")
	return actions
}

func (d *EventDeleter) NextScreen(i int) screens.Screen {
	if i != d.displayCount+1 {
		eventIdx := d.startIdx + i - 1
		d.Database.DeleteEvent(context.Background(), (*d.Events)[eventIdx])
		*d.Events = slices.Delete(*d.Events, eventIdx, eventIdx+1)
	}
	return d.Viewer
}
