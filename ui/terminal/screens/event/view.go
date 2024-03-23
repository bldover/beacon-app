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

type eventViewCache interface {
	GetPastEvents() []data.Event
	GetFutureEvents() []data.Event
}

type Viewer struct {
	title             string
	actions           []string
	sortType          sortType
	viewType          data.EventType
	cache             eventViewCache
	events            []data.Event
	page              int
	addEventScreen    screens.ContextScreen
	deleteEventScreen screens.ContextScreen
	returnScreen      screens.Screen
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

type sortType int

const (
	dateAsc = iota
	dateDesc
)

func NewViewScreen(title string, viewType data.EventType, addScreen, deleteScreen screens.ContextScreen, cache eventViewCache) *Viewer {
	view := Viewer{}
	view.title = title
	view.viewType = viewType
	view.actions = []string{"Next Page", "Prev Page", "Goto Page", "Toggle Sort", "Add Event", "Delete Event", "Main Menu"}
	view.sortType = dateAsc
	view.addEventScreen = addScreen
	view.deleteEventScreen = deleteScreen
	view.cache = cache
	return &view
}

func (v *Viewer) AddContext(context screens.ScreenContext) {
    v.returnScreen = context.ReturnScreen
}

func (v Viewer) Title() string {
	return v.title
}

func (v *Viewer) Refresh() {
	if v.viewType == data.Past {
		v.events = v.cache.GetPastEvents()
	} else {
		v.events = v.cache.GetFutureEvents()
	}

	v.sort()
}

func (v Viewer) DisplayData() {
	if len(v.events) == 0 {
		output.Displayln("No concerts found")
	}

	var eventData strings.Builder
	pageIndicator := fmt.Sprintf("Page %d/%d\n", v.page+1, v.numPages())
	eventData.WriteString(pageIndicator)
	startEvent := (v.page * pageSize)
	endEvent := startEvent + pageSize
	if endEvent > len(v.events) {
		endEvent = len(v.events)
	}

	for i := startEvent; i < endEvent; i++ {
		eventData.WriteString(format.FormatEvent((v.events)[i]))
	}
	output.Displayln(eventData.String())
}

func (v Viewer) Actions() []string {
	return v.actions
}

func (v *Viewer) NextScreen(i int) (screens.Screen, *screens.ScreenContext) {
	switch i {
	case nextPage:
		if (v.page + 1) < v.numPages() {
			v.page++
		}
		return v, nil
	case prevPage:
		if v.page > 0 {
			v.page--
		}
	case gotoPage:
		v.page = input.PromptAndGetInputNumeric("page number", 1, v.numPages()+1) - 1
	case toggleSort:
		if v.sortType == dateAsc {
			v.sortType = dateDesc
		} else {
			v.sortType = dateAsc
		}
		v.sort()
		v.page = 0
	case addEvent:
		return v.addEventScreen, screens.NewScreenContext(v, v.viewType)
	case deleteEvent:
		eventIdx := v.page * pageSize
		pageEvents := int(math.Min(float64(pageSize), float64(len(v.events) - eventIdx)))
		context := screens.NewScreenContext(v, v.events[eventIdx:eventIdx+pageEvents])
		return v.deleteEventScreen, context
	case mainMenu:
		v.page = 0
		return v.returnScreen, nil
	}
	return v, nil
}

func (v Viewer) numPages() int {
	return int(math.Ceil(float64(len(v.events)) / float64(pageSize)))
}

func (v *Viewer) sort() {
	sortFunc := data.EventSorterDateAsc()
	if v.sortType == dateDesc {
		sortFunc = data.EventSorterDateDesc()
	}
	slices.SortFunc(v.events, sortFunc)
}
