package finder

import (
	"concert-manager/ui/textui/output"
	"concert-manager/ui/textui/screens"
)

type Menu struct {
	UpcomingEventViewer    screens.Screen
	RecommendedEventViewer screens.Screen
	actions                []string
	returnScreen           screens.Screen
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

func (m *Menu) AddContext(context screens.ScreenContext) {
	m.returnScreen = context.ReturnScreen
}

func (m Menu) Title() string {
	return "Concert Finder"
}

func (m Menu) Actions() []string {
	return m.actions
}

func (m Menu) NextScreen(i int) (screens.Screen, *screens.ScreenContext) {
	switch i {
	case viewAllUpcoming:
		return m.UpcomingEventViewer, screens.NewScreenContext(m)
	case viewRecommended:
		output.Displayln("Not yet implemented!")
	case mainMenu:
		return m.returnScreen, nil
	}
	return m, nil
}
