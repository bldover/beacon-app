package screens

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/output"
	"fmt"
	"slices"
)

type OpenerRemover struct {
	AddEventScreen Screen
	openers   *[]data.Artist
}

func NewOpenerRemoveScreen() *OpenerRemover {
	or := OpenerRemover{}
	or.openers = &[]data.Artist{}
    return &or
}

func (or *OpenerRemover) AddOpenerContext(openers *[]data.Artist) {
	or.openers = openers
}

func (or OpenerRemover) Title() string {
    return "Remove Opener"
}

func (or OpenerRemover) DisplayData() {
	if len(*or.openers) == 0 {
		output.Displayln("No openers to remove!")
	}
}

func (or OpenerRemover) Actions() []string {
	actions := []string{}

	for _, opener := range *or.openers {
		actions = append(actions, fmt.Sprintf("%v", opener))
	}

	actions = append(actions, "Back")
    return actions
}

func (or *OpenerRemover) NextScreen(i int) Screen {
	if i == len(*or.openers) + 1 {
		return or.AddEventScreen
	}
	*or.openers = slices.Delete(*or.openers, i - 1, i)
	return or
}
