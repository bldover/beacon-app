package event

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/input"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"concert-manager/ui/terminal/screens/format"
	"fmt"
	"math"
	"slices"
	"strings"
)

const pageSize = 10

type eventDeleteScreen interface {
	screens.Screen
	AddDeleteContext(int, int)
}

type Viewer struct {
	Events            *[]data.Event
	AddEventScreen    screens.Screen
	DeleteEventScreen eventDeleteScreen
	MainMenu          screens.Screen
	page              int
	actions           []string
	title             string
}

const (
	nextPage = iota + 1
	prevPage
	gotoPage
	toggleSort
	addEvent
	deleteEvent
	mainMenu
)

func NewViewScreen(title string) *Viewer {
	view := Viewer{}
	view.title = title
	view.actions = []string{"Next Page", "Prev Page", "Goto Page", "Toggle Sort", "Add Event", "Delete Event", "Main Menu"}
	return &view
}

func (v Viewer) numPages() int {
	return int(math.Ceil(float64(len(*v.Events)) / float64(pageSize)))
}

func (v Viewer) Title() string {
	return v.title
}

func (v Viewer) DisplayData() {
	if len(*v.Events) == 0 {
		output.Displayln("No concerts found")
	}

	var eventData strings.Builder
	pageIndicator := fmt.Sprintf("Page %d/%d\n", v.page+1, v.numPages())
	eventData.WriteString(pageIndicator)
	startEvent := (v.page * pageSize)
	endEvent := startEvent + pageSize
	if endEvent > len(*v.Events) {
		endEvent = len(*v.Events)
	}

	for i := startEvent; i < endEvent; i++ {
		eventData.WriteString(format.FormatEvent((*v.Events)[i]))
	}
	output.Displayln(eventData.String())
}

func (v Viewer) Actions() []string {
	return v.actions
}

func (v *Viewer) NextScreen(i int) screens.Screen {
	switch i {
	case nextPage:
		if (v.page + 1) < v.numPages() {
			v.page++
		}
		return v
	case prevPage:
		if v.page > 0 {
			v.page--
		}
	case gotoPage:
		v.page = input.PromptAndGetInputNumeric("page number", 1, v.numPages()+1) - 1
	case toggleSort:
		sortByDate := data.EventSorter()
		if slices.IsSortedFunc(*v.Events, sortByDate) {
			slices.Reverse(*v.Events)
		} else {
			slices.SortFunc(*v.Events, sortByDate)
		}
		v.page = 0
	case addEvent:
		return v.AddEventScreen
	case deleteEvent:
		v.DeleteEventScreen.AddDeleteContext(pageSize*v.page, pageSize)
		return v.DeleteEventScreen
	case mainMenu:
		v.page = 0
		return v.MainMenu
	}
	return v
}
