package cli

import "os"

type MainMenu struct {
    Children []Screen
}

func (mm MainMenu) Title() string {
    return "Main Menu"
}

func (mm MainMenu) Data() string {
    return ""
}

func (mm MainMenu) Actions() []string {
	actions := []string{}
	for _, child := range mm.Children {
		actions = append(actions, child.Title())
	}
	actions = append(actions, "Exit")
    return actions
}

func (mm MainMenu) NextScreen(i int) Screen {
	if i <= len(mm.Children) {
		return mm.Children[i - 1]
	}
	if i == len(mm.Children) + 1 {
		os.Exit(0)
	}
	return mm
}

func (mm MainMenu) Parent() Screen {
    return nil
}
