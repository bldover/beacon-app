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
		return m.upcomingEventViewer, screens.NewScreenContext(m)
	case viewRecommended:
		output.Displayln("Not yet implemented!")
	case mainMenu:
		return m.returnScreen, nil
	}
	return m, nil
}
