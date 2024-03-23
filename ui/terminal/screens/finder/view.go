package finder

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
	"time"
)

const pageSize = 10
const reloadTimeFormat = "2006-01-02T15:04:05"

type eventRetrievalCache interface {
	GetUpcomingEvents(string, string) []data.EventDetails
	ReloadUpcomingEvents(string, string) error
}

type Finder struct {
	title                  string
	city                   string
	state                  string
	actions                []string
	sortType               sortType
	cache                  eventRetrievalCache
	events                 []data.EventDetails
	page                   int
	loaded                 bool
	lastLoad               string
	addEventSelectorScreen screens.ContextScreen
	returnScreen           screens.Screen
}

const (
	nextPage = iota + 1
	prevPage
	gotoPage
	toggleSort
	addEvent
	changeLocation
	refreshEvents
	finderMenu
)

type sortType int

const (
	dateAsc = iota
	dateDesc
)

func NewViewScreen(title, city, state string, addSelector screens.ContextScreen, cache eventRetrievalCache) *Finder {
	view := Finder{}
	view.title = title
	view.city = city
	view.state = state
	view.actions = []string{"Next Page", "Prev Page", "Goto Page", "Toggle Sort",
		"Save Concert", "Change Location", "Refresh Concerts", "Finder Menu"}
	view.sortType = dateAsc
	view.addEventSelectorScreen = addSelector
	view.cache = cache
	return &view
}

func (f *Finder) AddContext(context screens.ScreenContext) {
    f.returnScreen = context.ReturnScreen
}

func (f Finder) Title() string {
	return f.title
}

func (f *Finder) DisplayData() {
	if !f.loaded {
		f.reloadEvents()
	}
	if len(f.events) == 0 {
		output.Displayln("No events found")
	}

	var eventData strings.Builder
	pageIndicator := fmt.Sprintf("Page %d/%d - Last reloaded: %v\n", f.page+1, f.numPages(), f.lastLoad)
	eventData.WriteString(pageIndicator)

	startEvent := (f.page * pageSize)
	endEvent := startEvent + pageSize
	if endEvent > len(f.events) {
		endEvent = len(f.events)
	}

	for i := startEvent; i < endEvent; i++ {
		eventData.WriteString(format.FormatEventDetails((f.events)[i]))
	}
	output.Displayln(eventData.String())
}

func (f Finder) Actions() []string {
	return f.actions
}

func (f *Finder) NextScreen(i int) (screens.Screen, *screens.ScreenContext) {
	switch i {
	case nextPage:
		if (f.page + 1) < f.numPages() {
			f.page++
		}
		return f, nil
	case prevPage:
		if f.page > 0 {
			f.page--
		}
	case gotoPage:
		f.page = input.PromptAndGetInputNumeric("page number", 1, f.numPages()+1) - 1
	case toggleSort:
		if f.sortType == dateAsc {
			f.sortType = dateDesc
		} else {
			f.sortType = dateAsc
		}
		f.sort()
		f.page = 0
	case addEvent:
		startIdx := pageSize * f.page
		endIdx := startIdx + pageSize
		return f.addEventSelectorScreen, screens.NewScreenContext(f, f.events[startIdx:endIdx])
	case changeLocation:
		f.getNewLocation()
	case refreshEvents:
		f.reloadEvents()
	case finderMenu:
		f.page = 0
		return f.returnScreen, nil
	}
	return f, nil
}

func (f *Finder) getNewLocation() {
	f.city = input.PromptAndGetInput("city", input.OnlyLettersOrSpacesValidation)
	f.state = input.PromptAndGetInput("state code", input.StateValidation)
	f.events = f.cache.GetUpcomingEvents(f.city, f.state)
	f.page = 0
}

func (f *Finder) reloadEvents() error {
	output.Displayf("Reloading concerts for %s, %s...", f.city, f.state)
	err := f.cache.ReloadUpcomingEvents(f.city, f.state)
	f.events = f.cache.GetUpcomingEvents(f.city, f.state)
	f.loaded = true
	f.lastLoad = time.Now().Format(reloadTimeFormat)
	f.page = 0
	output.ClearCurrentLine()
	return err
}

func (f Finder) numPages() int {
	return int(math.Ceil(float64(len(f.events)) / float64(pageSize)))
}

func (f *Finder) sort() {
	sortFunc := data.EventDetailsSorterDateAsc()
	if f.sortType == dateDesc {
		sortFunc = data.EventDetailsSorterDateDesc()
	}
	slices.SortFunc(f.events, sortFunc)
}
