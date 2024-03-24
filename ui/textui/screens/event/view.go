package event

import (
	"concert-manager/data"
	"concert-manager/ui/textui/input"
	"concert-manager/ui/textui/output"
	"concert-manager/ui/textui/screens"
	"concert-manager/util/format"
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
	ScreenTitle       string
	ViewType          data.EventType
	AddEventScreen    screens.Screen
	DeleteEventScreen screens.Screen
	Cache             eventViewCache
	actions           []string
	sortType          sortType
	events            []data.Event
	page              int
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

func NewViewScreen() *Viewer {
	view := Viewer{}
	view.actions = []string{"Next Page", "Prev Page", "Goto Page", "Toggle Sort", "Add Event", "Delete Event", "Main Menu"}
	view.sortType = dateAsc
	return &view
}

func (v *Viewer) AddContext(context screens.ScreenContext) {
	v.returnScreen = context.ReturnScreen
}

func (v Viewer) Title() string {
	return v.ScreenTitle
}

func (v *Viewer) Refresh() {
	if v.ViewType == data.Past {
		v.events = v.Cache.GetPastEvents()
	} else {
		v.events = v.Cache.GetFutureEvents()
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
		return v.AddEventScreen, screens.NewScreenContext(v)
	case deleteEvent:
		eventIdx := v.page * pageSize
		pageEvents := int(math.Min(float64(pageSize), float64(len(v.events)-eventIdx)))
		context := screens.NewScreenContext(v, v.events[eventIdx:eventIdx+pageEvents])
		return v.DeleteEventScreen, context
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
