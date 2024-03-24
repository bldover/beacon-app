package event

import (
	"concert-manager/data"
	"concert-manager/ui/textui/output"
	"concert-manager/ui/textui/screens"
	"concert-manager/util/format"
)

type eventDeleteCache interface {
	DeleteEvent(data.Event) error
}

type Deleter struct {
	Cache        eventDeleteCache
	events       []data.Event
	returnScreen screens.Screen
}

func (d *Deleter) AddContext(context screens.ScreenContext) {
	d.returnScreen = context.ReturnScreen
	d.events = context.Props[0].([]data.Event)
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

func (d *Deleter) NextScreen(i int) (screens.Screen, *screens.ScreenContext) {
	if i != len(d.events)+1 {
		if err := d.Cache.DeleteEvent(d.events[i-1]); err != nil {
			output.Displayf("Failed to delete event: %v\n", err)
			return d, nil
		}
	}
	return d.returnScreen, nil
}
