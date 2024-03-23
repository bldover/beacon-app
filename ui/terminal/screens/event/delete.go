package event

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"concert-manager/ui/terminal/screens/format"
)

type eventDeleteCache interface {
	DeleteEvent(data.Event) error
}

type Deleter struct {
	cache        eventDeleteCache
	events       []data.Event
	returnScreen screens.Screen
}

func NewDeleteScreen(cache eventDeleteCache) *Deleter {
	d := Deleter{cache: cache}
	d.cache = cache
	return &d
}

func (d *Deleter) AddContext(returnScreen screens.Screen, props ...any) {
	d.returnScreen = returnScreen
	d.events = props[0].([]data.Event)
}

func (d Deleter) Title() string {
	return "Delete Concert"
}

func (d Deleter) DisplayData() {
	if len(d.events) == 0 {
		output.Displayln("No concerts found")
	}
}

func (d Deleter) Actions() []string {
	actions := []string{}
	actions = append(actions, format.FormatEventsShort(d.events)...)
	actions = append(actions, "Cancel")
	return actions
}

func (d *Deleter) NextScreen(i int) screens.Screen {
	if i != len(d.events)+1 {
		if err := d.cache.DeleteEvent(d.events[i-1]); err != nil {
			output.Displayf("Failed to delete event: %v\n", err)
			return d
		}
	}
	return d.returnScreen
}
