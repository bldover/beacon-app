package artist

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"fmt"
	"slices"
)

type OpenerRemover struct {
	openers      *[]data.Artist
	returnScreen screens.Screen
}

func NewOpenerRemoveScreen() *OpenerRemover {
	or := OpenerRemover{}
	or.openers = &[]data.Artist{}
	return &or
}

func (or *OpenerRemover) AddContext(returnScreen screens.Screen, props ...any) {
	or.returnScreen = returnScreen
	or.openers = props[0].(*[]data.Artist)
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

func (or *OpenerRemover) NextScreen(i int) screens.Screen {
	if i == len(*or.openers)+1 {
		return or.returnScreen
	}
	*or.openers = slices.Delete(*or.openers, i-1, i)
	return or
}
