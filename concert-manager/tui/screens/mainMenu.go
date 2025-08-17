package screens

import (
	"concert-manager/tui/output"
	"os"
)

type MainMenu struct {
	Children map[int]Screen
}

func NewMainMenu() *MainMenu {
	mm := MainMenu{}
	mm.Children = make(map[int]Screen)
	return &mm
}

func (mm MainMenu) AddReturnContext(Screen) {
	// noop, this is the only screen that never needs to return
}

func (mm MainMenu) Title() string {
	return "Main Menu"
}

func (mm MainMenu) Actions() []string {
	actions := []string{}
	for i := 1; i <= len(mm.Children); i++ {
		actions = append(actions, mm.Children[i].Title())
	}
	actions = append(actions, "Exit")
	return actions
}

func (mm MainMenu) NextScreen(i int) Screen {
	if i == len(mm.Children)+1 {
		output.Displayln("Received exit request, terminating...")
		os.Exit(0)
	}
	next := mm.Children[i]
	return next
}
