package event

import (
	"concert-manager/cli"
	"concert-manager/cli/format"
	"concert-manager/data"
	"concert-manager/out"
	"fmt"
	"math"
	"slices"
	"strings"
)

const pageSize = 10

type Viewer struct {
	Events   *[]data.Event
	AddEventScreen cli.Screen
	DeleteEventScreen EventDeleteScreen
	MainMenu cli.Screen
	page     int
	actions []string
	title string
}

type EventDeleteScreen interface {
	cli.Screen
	AddDeleteContext(int, int)
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

func NewViewerScreen(title string) *Viewer {
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
		out.Displayln("No concerts found")
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
	out.Displayln(eventData.String())
}

func (v Viewer) Actions() []string {
	return v.actions
}

func (v *Viewer) NextScreen(i int) cli.Screen {
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
		v.page = cli.PromptAndGetInputNumeric("page number", 1, v.numPages() + 1) - 1
	case toggleSort:
		sortByDate := eventSort()
		if slices.IsSortedFunc(*v.Events, sortByDate) {
			slices.Reverse(*v.Events)
		} else {
			slices.SortFunc(*v.Events, sortByDate)
		}
		v.page = 0
	case addEvent:
		return v.AddEventScreen
	case deleteEvent:
		v.DeleteEventScreen.AddDeleteContext(pageSize * v.page, pageSize)
		return v.DeleteEventScreen
	case mainMenu:
		v.page = 0
		return v.MainMenu
	}
	return v
}

func eventSort() func(a, b data.Event) int {
	return func(a, b data.Event) int {
		if data.Timestamp(a.Date).Before(data.Timestamp(b.Date)) {
			return -1
		} else if data.Timestamp(a.Date).After(data.Timestamp(b.Date)) {
			return 1
		} else {
			return 0
		}
	}
}
