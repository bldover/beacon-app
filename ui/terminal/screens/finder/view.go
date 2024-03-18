package finder

import (
	"concert-manager/data"
	"concert-manager/finder"
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

type addEventSelectorScreen interface {
	screens.Screen
	AddContext([]data.EventDetails)
}

type eventFinder interface {
	FindAllEvents(finder.FindEventRequest) ([]data.EventDetails, error)
}

type Finder struct {
	AddEventSelectorScreen addEventSelectorScreen
	FinderMenu             screens.Screen
	EventFinder            eventFinder
	events                 []data.EventDetails
	page                   int
	actions                []string
	title                  string
	city                   string
	state                  string
	loaded                 bool
	lastLoad               string
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

func NewViewerScreen(title string, city string, state string) *Finder {
	view := Finder{}
	view.title = title
	view.city = city
	view.state = state
	view.actions = []string{"Next Page", "Prev Page", "Goto Page", "Toggle Sort",
		"Save Concert", "Change Location", "Refresh Concerts", "Finder Menu"}
	return &view
}

func (f Finder) numPages() int {
	return int(math.Ceil(float64(len(f.events)) / float64(pageSize)))
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

func (f *Finder) NextScreen(i int) screens.Screen {
	switch i {
	case nextPage:
		if (f.page + 1) < f.numPages() {
			f.page++
		}
		return f
	case prevPage:
		if f.page > 0 {
			f.page--
		}
	case gotoPage:
		f.page = input.PromptAndGetInputNumeric("page number", 1, f.numPages()+1) - 1
	case toggleSort:
		sortByDate := data.EventDetailsSorter()
		if slices.IsSortedFunc(f.events, sortByDate) {
			slices.Reverse(f.events)
		} else {
			slices.SortFunc(f.events, sortByDate)
		}
		f.page = 0
	case addEvent:
		startIdx := pageSize * f.page
		endIdx := startIdx + pageSize
		f.AddEventSelectorScreen.AddContext(f.events[startIdx:endIdx])
		return f.AddEventSelectorScreen
	case changeLocation:
		f.getNewLocation()
	case refreshEvents:
		f.reloadEvents()
	case finderMenu:
		f.page = 0
		return f.FinderMenu
	}
	return f
}

func (f *Finder) getNewLocation() {
	f.city = input.PromptAndGetInput("city", input.OnlyLettersOrSpacesValidation)
	f.state = input.PromptAndGetInput("state code", input.StateValidation)
	f.reloadEvents()
}

func (f *Finder) reloadEvents() error {
	output.Displayf("Reloading concerts for %s, %s...", f.city, f.state)
	request := finder.FindEventRequest{City: f.city, State: f.state}
	events, err := f.EventFinder.FindAllEvents(request)
	f.events = events
	f.loaded = true
	f.lastLoad = time.Now().Format(reloadTimeFormat)
	f.page = 0
	output.ClearCurrentLine()
	return err
}
