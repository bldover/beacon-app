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

type Deleter struct {
	Events       *[]data.Event
	Database     eventDbDeleter
	Viewer       screens.Screen
	startIdx     int
	displayCount int
}

func NewDeleteScreen() *Deleter {
	return &Deleter{}
}

func (d *Deleter) AddDeleteContext(startIdx int, displayCount int) {
	d.startIdx = startIdx
	remainingEvents := len(*d.Events) - startIdx
	d.displayCount = int(math.Min(float64(displayCount), float64(remainingEvents)))
}

func (d Deleter) Title() string {
	return "Delete Concert"
}

func (d Deleter) DisplayData() {
	if len(*d.Events) == 0 {
		output.Displayln("No concerts found")
	}
}

func (d Deleter) Actions() []string {
	actions := []string{}
	pageEvents := (*d.Events)[d.startIdx : d.startIdx+d.displayCount]
	actions = append(actions, format.FormatEventsShort(pageEvents)...)
	actions = append(actions, "Cancel")
	return actions
}

func (d *Deleter) NextScreen(i int) screens.Screen {
	if i != d.displayCount+1 {
		eventIdx := d.startIdx + i - 1
		d.Database.DeleteEvent(context.Background(), (*d.Events)[eventIdx])
		*d.Events = slices.Delete(*d.Events, eventIdx, eventIdx+1)
	}
	return d.Viewer
}
