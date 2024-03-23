package screens

import (
	"concert-manager/ui/terminal/output"
	"os"
)

type MainMenu struct {
    Children map[int]ContextScreen
}

func NewMainMenu() *MainMenu {
	mm := MainMenu{}
	mm.Children = make(map[int]ContextScreen)
    return &mm
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

func (mm MainMenu) NextScreen(i int) (Screen, *ScreenContext) {
	if i == len(mm.Children) + 1 {
		output.Displayln("Received exit request, terminating...")
		os.Exit(0)
	}
	next := mm.Children[i]
	return next, NewScreenContext(mm)
}
