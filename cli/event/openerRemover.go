package event

import (
	"concert-manager/cli"
	"concert-manager/data"
	"concert-manager/out"
	"fmt"
	"slices"
)

type OpenerRemove struct {
	AddEventScreen cli.Screen
	openers   *[]data.Artist
}

func NewOpenerRemovalScreen() *OpenerRemove {
	or := OpenerRemove{}
	or.openers = &[]data.Artist{}
    return &or
}

func (or *OpenerRemove) AddOpenerContext(openers *[]data.Artist) {
	or.openers = openers
}

func (or OpenerRemove) Title() string {
    return "Remove Opener"
}

func (or OpenerRemove) DisplayData() {
	if len(*or.openers) == 0 {
		out.Displayln("No openers to remove!")
	}
}

func (or OpenerRemove) Actions() []string {
	actions := []string{}

	for _, opener := range *or.openers {
		actions = append(actions, fmt.Sprintf("%v", opener))
	}

	actions = append(actions, "Back")
    return actions
}

func (or *OpenerRemove) NextScreen(i int) cli.Screen {
	if i == len(*or.openers) + 1 {
		return or.AddEventScreen
	}
	*or.openers = slices.Delete(*or.openers, i - 1, i)
	return or
}
