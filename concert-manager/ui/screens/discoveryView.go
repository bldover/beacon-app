package screens

import (
	"concert-manager/data"
	"concert-manager/finder"
	"concert-manager/log"
	"concert-manager/ui/input"
	"concert-manager/ui/output"
	"concert-manager/util"
	"fmt"
	"math"
	"slices"
	"strings"
)

type eventRetrievalCache interface {
	GetUpcomingEvents() []data.EventDetails
	Invalidate()
	ChangeLocation(string, string)
	GetLocation() finder.Location
}

type DiscoveryViewer struct {
	SearchResultScreen *DiscoverySearchResult
	AddEventScreen     *EventAdder
	Cache              eventRetrievalCache
	actions            []string
	events             []data.EventDetails
	sortType           sortType
	page               int
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
		"Save Event", "Search Events", "Change Location", "Refresh Events", "Discovery Menu"}
	view.sortType = dateAsc
	return &view
}

func (v DiscoveryViewer) Title() string {
	return "All Upcoming Events"
}

func (v *DiscoveryViewer) Refresh() {
	output.Displayf("Retrieving events for %s...", v.Cache.GetLocation())
	v.events = v.Cache.GetUpcomingEvents() // ignore possible refresh, this view maintains its own state
	v.sort()
	v.page = 0
	output.ClearCurrentLine()
}

func (v *DiscoveryViewer) DisplayData() {
	if v.events == nil {
		v.Refresh()
	}

	var eventData strings.Builder
	pageIndicator := fmt.Sprintf("Page %d/%d\n", v.page+1, v.numPages())
	eventData.WriteString(pageIndicator)

	if len(v.events) == 0 {
		eventData.WriteString("(none)")
		output.Displayln(eventData.String())
		return
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
					v.SearchResultScreen.Events = util.SearchEventDetailsByArtist(name, v.events, util.NoMaxResults, util.LenientTolerance)
				case searchByVenue:
					name := input.PromptAndGetInput("venue name to search", input.NoValidation)
					v.SearchResultScreen.Events = util.SearchEventDetailsByVenue(name, v.events, util.NoMaxResults, util.LenientTolerance)
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
		v.changeLocation()
	case refreshEvents:
		v.Cache.Invalidate()
		v.events = nil
	case discoveryViewToMenu:
		v.page = 0
		return nil
	}
	return v
}

func (v *DiscoveryViewer) changeLocation() {
	city := input.PromptAndGetInput("city", input.OnlyLettersOrSpacesValidation)
	state := input.PromptAndGetInput("state code", input.StateValidation)
	v.Cache.ChangeLocation(city, state)
	v.events = nil
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
