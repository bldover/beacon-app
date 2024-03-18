package finder

import (
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
)

type Menu struct {
	UpcomingEventViewer screens.Screen
	RecommendedEventViewer screens.Screen
	MainMenu screens.Screen
	actions []string
}

const (
	viewAllUpcoming = iota + 1
	viewRecommended
	mainMenu
)

func NewMenu() *Menu {
	menu := Menu{}
	menu.actions = []string{"All Upcoming Concerts", "Recommended Concerts", "Main Menu"}
    return &menu
}

func (m Menu) Title() string {
    return "Concert Finder"
}

func (m Menu) Actions() []string {
    return m.actions
}

func (m Menu) NextScreen(i int) screens.Screen {
	switch i {
    case viewAllUpcoming:
		return m.UpcomingEventViewer
	case viewRecommended:
		output.Displayln("Not yet implemented!")
	case mainMenu:
		return m.MainMenu
	}
	return m
}
