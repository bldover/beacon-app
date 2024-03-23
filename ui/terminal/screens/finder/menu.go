package finder

import (
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
)

type Menu struct {
	actions []string
	upcomingEventViewer screens.ContextScreen
	recommendedEventViewer screens.ContextScreen
	returnScreen screens.Screen
}

const (
	viewAllUpcoming = iota + 1
	viewRecommended
	mainMenu
)

func NewMenu(upcomingViewer screens.ContextScreen, recommendedViewer screens.ContextScreen) *Menu {
	menu := Menu{}
	menu.upcomingEventViewer = upcomingViewer
	menu.recommendedEventViewer = recommendedViewer
	menu.actions = []string{"All Upcoming Concerts", "Recommended Concerts", "Main Menu"}
    return &menu
}

func (m *Menu) AddContext(returnScreen screens.Screen, _ ...any) {
    m.returnScreen = returnScreen
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
		m.upcomingEventViewer.AddContext(m)
		return m.upcomingEventViewer
	case viewRecommended:
		output.Displayln("Not yet implemented!")
	case mainMenu:
		return m.returnScreen
	}
	return m
}
