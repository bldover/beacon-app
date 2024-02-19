package event

import (
	"concert-manager/cli"
	"concert-manager/cli/format"
	"concert-manager/data"
	"concert-manager/out"
	"context"
	"math"
	"slices"
)

type EventDeleter interface {
    DeleteEvent(context.Context, data.Event) error
}

type Deleter struct {
	Events   *[]data.Event
	Deleter EventDeleter
	Viewer cli.Screen
	startIdx int
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
		out.Displayln("No concerts found")
	}
}

func (d Deleter) Actions() []string {
	actions := []string{}
	pageEvents := []data.Event{}
	var maxNameLen int
	for i := d.startIdx; i < d.startIdx + d.displayCount; i++ {
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

func (d *Deleter) NextScreen(i int) cli.Screen {
	if i != d.displayCount + 1 {
		eventIdx := d.startIdx + i - 1
		d.Deleter.DeleteEvent(context.Background(), (*d.Events)[eventIdx])
		*d.Events = slices.Delete(*d.Events, eventIdx, eventIdx + 1)
	}
    return d.Viewer
}
