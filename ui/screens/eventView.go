package screens

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/ui/input"
	"concert-manager/ui/output"
	"concert-manager/util"
	"fmt"
	"math"
	"slices"
	"strings"
)

type eventViewCache interface {
	GetSavedEvents() []data.Event
	DeleteSavedEvent(string) error
}

type SavedEventViewer struct {
	SearchResultScreen *EventSearchResult
	AddEventScreen     *EventAdder
	Cache              eventViewCache
	actions            []string
	sortType           sortType
	events             []data.Event
	page               int
}

const (
	nextEventPage = iota + 1
	prevEventPage
	gotoEventPage
	toggleEventSort
	addEvent
	deleteEvent
	searchSavedEvents
	eventViewToMainMenu
)

func NewSavedEventViewScreen() *SavedEventViewer {
	view := SavedEventViewer{}
	view.actions = []string{"Next Page", "Prev Page", "Goto Page", "Toggle Sort", "Add Event",
		"Delete Event", "Search Events", "Main Menu"}
	view.sortType = dateAsc
	return &view
}

func (v SavedEventViewer) Title() string {
	return "All Saved Events"
}

func (v *SavedEventViewer) Refresh() {
	v.events = v.Cache.GetSavedEvents()
	v.sort()
}

func (v SavedEventViewer) DisplayData() {
	var eventData strings.Builder
	pageIndicator := fmt.Sprintf("Page %d/%d\n", v.page+1, v.numPages())
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
		eventData.WriteString(util.FormatEvent((v.events)[i]))
	}
	output.Displayln(eventData.String())
}

func (v SavedEventViewer) Actions() []string {
	return v.actions
}

func (v *SavedEventViewer) NextScreen(i int) Screen {
	switch i {
	case nextEventPage:
		if (v.page + 1) < v.numPages() {
			v.page++
		}
	case prevEventPage:
		if v.page > 0 {
			v.page--
		}
	case gotoEventPage:
		v.page = input.PromptAndGetInputNumeric("page number", 1, v.numPages()+1) - 1
	case toggleEventSort:
		if v.sortType == dateAsc {
			v.sortType = dateDesc
		} else {
			v.sortType = dateAsc
		}
		v.sort()
		v.page = 0
	case addEvent:
		return v.AddEventScreen
	case deleteEvent:
		startIdx := v.page * pageSize
		endIdx := int(math.Min(float64(startIdx + pageSize), float64(len(v.events))))
		selectScreen := &Selector[data.Event]{
			ScreenTitle: "Delete Event",
			Next:        v,
			Options:     v.events[startIdx : endIdx],
			HandleSelect: func(e data.Event) {
				if err := v.Cache.DeleteSavedEvent(e.Id); err != nil {
					output.Displayf("Failed to delete event: %v\n", err)
				}
			},
			Formatter: util.FormatEventsShort,
		}
		return selectScreen
	case searchSavedEvents:
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
					v.SearchResultScreen.Events = util.SearchEventsByArtists(name, v.Cache.GetSavedEvents(), util.NoMaxResults, util.LenientTolerance)
				case searchByVenue:
					name := input.PromptAndGetInput("venue name to search", input.NoValidation)
					v.SearchResultScreen.Events = util.SearchEventsByVenue(name, v.Cache.GetSavedEvents(), util.NoMaxResults, util.LenientTolerance)
				default:
					output.Display("Internal error! Check the logs")
					log.Error("Invalid search type selection:", s)
					v.SearchResultScreen.Events = []data.Event{}
				}
			},
			Formatter: IdentityTransform[string],
		}
		return selectScreen
	case eventViewToMainMenu:
		v.page = 0
		return nil
	}
	return v
}

func (v SavedEventViewer) numPages() int {
	return int(math.Ceil(float64(len(v.events)) / float64(pageSize)))
}

func (v *SavedEventViewer) sort() {
	sortFunc := util.EventSorterDateAsc()
	if v.sortType == dateDesc {
		sortFunc = util.EventSorterDateDesc()
	}
	slices.SortFunc(v.events, sortFunc)
}
