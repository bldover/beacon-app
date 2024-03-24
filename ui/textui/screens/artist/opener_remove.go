package artist

import (
	"concert-manager/data"
	"concert-manager/ui/textui/output"
	"concert-manager/ui/textui/screens"
	"fmt"
	"slices"
)

type OpenerRemover struct {
	openers      *[]data.Artist
	returnScreen screens.Screen
}

func NewOpenerRemoveScreen() *OpenerRemover {
	return &OpenerRemover{openers: new([]data.Artist)}
}

func (or *OpenerRemover) AddContext(context screens.ScreenContext) {
	or.returnScreen = context.ReturnScreen
	or.openers = context.Props[0].(*[]data.Artist)
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

func (or *OpenerRemover) NextScreen(i int) (screens.Screen, *screens.ScreenContext) {
	if i == len(*or.openers)+1 {
		return or.returnScreen, nil
	}
	*or.openers = slices.Delete(*or.openers, i-1, i)
	return or, nil
}
