package screens

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/ui/textui/input"
	"concert-manager/ui/textui/output"
	"concert-manager/util"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"
)

const reloadTimeFormat = "2006-01-02T15:04:05"

type eventRetrievalCache interface {
	GetUpcomingEvents(string, string) []data.EventDetails
	ReloadUpcomingEvents(string, string) error
}

type upcomingEventSearch interface {
	FindFuzzyEventDetailsMatchesByArtist(string, string, string) []data.EventDetails
	FindFuzzyEventDetailsMatchesByVenue(string, string, string) []data.EventDetails
}

type DiscoveryViewer struct {
	City               string
	State              string
	Search             upcomingEventSearch
	SearchResultScreen *DiscoverySearchResult
	AddEventScreen     *EventAdder
	Cache              eventRetrievalCache
	actions            []string
	events             []data.EventDetails
	sortType           sortType
	page               int
	loaded             bool
	lastLoad           string
}

const (
	nextDiscoveryViewerPage = iota + 1
	prevDiscoveryViewerPage
	gotoDiscoveryViewerPage
	toggleDiscoveryViewerSort
	addDiscoveryViewerEvent
	searchDiscoveryEvents
	changeLocation
	refreshEvents
	discoveryViewToMenu
)

func NewDiscoveryViewScreen() *DiscoveryViewer {
	view := DiscoveryViewer{}
	view.actions = []string{"Next Page", "Prev Page", "Goto Page", "Toggle Sort",
		"Save Event", "Search Events", "Change Location", "Refresh Event", "Discovery Menu"}
	view.sortType = dateAsc
	return &view
}

func (v DiscoveryViewer) Title() string {
	return "All Upcoming Events"
}

func (v *DiscoveryViewer) DisplayData() {
	if !v.loaded {
		v.reloadEvents()
	}

	var eventData strings.Builder
	pageIndicator := fmt.Sprintf("Page %d/%d - Last reloaded: %v\n", v.page+1, v.numPages(), v.lastLoad)
	eventData.WriteString(pageIndicator)

	if len(v.events) == 0 {
		output.Displayln("No events found")
	}

	startEvent := (v.page * pageSize)
	endEvent := startEvent + pageSize
	if endEvent > len(v.events) {
		endEvent = len(v.events)
	}

	for i := startEvent; i < endEvent; i++ {
		eventData.WriteString(util.FormatEventDetails((v.events)[i]))
	}
	output.Displayln(eventData.String())
}

func (v DiscoveryViewer) Actions() []string {
	return v.actions
}

func (v *DiscoveryViewer) NextScreen(i int) Screen {
	switch i {
	case nextDiscoveryViewerPage:
		if (v.page + 1) < v.numPages() {
			v.page++
		}
	case prevDiscoveryViewerPage:
		if v.page > 0 {
			v.page--
		}
	case gotoDiscoveryViewerPage:
		v.page = input.PromptAndGetInputNumeric("page number", 1, v.numPages()+1) - 1
	case toggleDiscoveryViewerSort:
		if v.sortType == dateAsc {
			v.sortType = dateDesc
		} else {
			v.sortType = dateAsc
		}
		v.sort()
		v.page = 0
	case addDiscoveryViewerEvent:
		startIdx := pageSize * v.page
		endIdx := int(math.Min(float64(startIdx + pageSize), float64(len(v.events))))
		selectScreen := &Selector[data.EventDetails]{
			ScreenTitle: "Select Event",
			Next:        v.AddEventScreen,
			Options:     v.events[startIdx:endIdx],
			HandleSelect: func(e data.EventDetails) {
				v.AddEventScreen.newEvent = e.Event
			},
			Formatter: util.FormatEventDetailsShort,
		}
		return selectScreen
	case searchDiscoveryEvents:
		const searchByArtist = "Search by Artist"
		const searchByVenue = "Search by Venue"
		selectScreen := &Selector[string]{
			ScreenTitle: "Select Search Type",
			Next:        v.SearchResultScreen,
			Options:     []string{searchByArtist, searchByVenue},
			HandleSelect: func(s string) {
				switch s {
				case searchByArtist:
					name := input.PromptAndGetInput("artist name to search", input.NoValidation)
					v.SearchResultScreen.Events = v.Search.FindFuzzyEventDetailsMatchesByArtist(name, v.City, v.State)
				case searchByVenue:
					name := input.PromptAndGetInput("venue name to search", input.NoValidation)
					v.SearchResultScreen.Events = v.Search.FindFuzzyEventDetailsMatchesByVenue(name, v.City, v.State)
				default:
					output.Display("Internal error! Check the logs")
					log.Error("Invalid search type selection:", s)
					v.SearchResultScreen.Events = []data.EventDetails{}
				}
			},
			Formatter: IdentityTransform[string],
		}
		return selectScreen
	case changeLocation:
		v.getNewLocation()
	case refreshEvents:
		v.reloadEvents()
	case discoveryViewToMenu:
		v.page = 0
		return nil
	}
	return v
}

func (v *DiscoveryViewer) getNewLocation() {
	v.City = input.PromptAndGetInput("city", input.OnlyLettersOrSpacesValidation)
	v.State = input.PromptAndGetInput("state code", input.StateValidation)
	output.Displayf("Retrieving concerts for %s, %s...", v.City, v.State)
	v.events = v.Cache.GetUpcomingEvents(v.City, v.State)
	output.ClearCurrentLine()
	v.page = 0
}

func (v *DiscoveryViewer) reloadEvents() error {
	output.Displayf("Retrieving concerts for %s, %s...", v.City, v.State)
	err := v.Cache.ReloadUpcomingEvents(v.City, v.State)
	v.events = v.Cache.GetUpcomingEvents(v.City, v.State)
	v.loaded = true
	v.lastLoad = time.Now().Format(reloadTimeFormat)
	v.page = 0
	output.ClearCurrentLine()
	return err
}

func (v DiscoveryViewer) numPages() int {
	return int(math.Ceil(float64(len(v.events)) / float64(pageSize)))
}

func (v *DiscoveryViewer) sort() {
	sortFunc := util.EventDetailsSorterDateAsc()
	if v.sortType == dateDesc {
		sortFunc = util.EventDetailsSorterDateDesc()
	}
	slices.SortFunc(v.events, sortFunc)
}
