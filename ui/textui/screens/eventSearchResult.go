package screens

import (
	"concert-manager/data"
	"concert-manager/ui/textui/input"
	"concert-manager/ui/textui/output"
	"concert-manager/util"
	"fmt"
	"math"
	"strings"
)

type eventSearchResultCache interface {
	DeleteSavedEvent(data.Event) error
}

type EventSearchResult struct {
	AddEventScreen   *EventAdder
	Cache            eventSearchResultCache
	Events           []data.Event
	actions          []string
	page             int
}

const (
	nextEventSearchResultPage = iota + 1
	prevEventSearchResultPage
	gotoEventSearchResultPage
	addEventSearchResultEvent
	deleteEventSearchResultEvent
	allSavedEvents
)

func NewEventSearchResultScreen() *EventSearchResult {
	view := EventSearchResult{}
	view.actions = []string{"Next Page", "Prev Page", "Goto Page", "Save Event", "Delete Event", "All Saved Events"}
	return &view
}

func (s EventSearchResult) Title() string {
	return "Search Results"
}

func (s *EventSearchResult) DisplayData() {
	var eventData strings.Builder
	pageIndicator := fmt.Sprintf("Page %d/%d\n", s.page+1, s.numPages())
	eventData.WriteString(pageIndicator)

	if len(s.Events) == 0 {
		eventData.WriteString("No events found")
	}

	startEvent := (s.page * pageSize)
	endEvent := startEvent + pageSize
	if endEvent > len(s.Events) {
		endEvent = len(s.Events)
	}

	for i := startEvent; i < endEvent; i++ {
		eventData.WriteString(util.FormatEvent((s.Events)[i]))
	}
	output.Displayln(eventData.String())
}

func (s EventSearchResult) Actions() []string {
	return s.actions
}

func (s *EventSearchResult) NextScreen(i int) Screen {
	switch i {
	case nextEventSearchResultPage:
		if (s.page + 1) < s.numPages() {
			s.page++
		}
		return s
	case prevEventSearchResultPage:
		if s.page > 0 {
			s.page--
		}
	case gotoEventSearchResultPage:
		s.page = input.PromptAndGetInputNumeric("page number", 1, s.numPages()+1) - 1
	case addEventSearchResultEvent:
		startIdx := pageSize * s.page
		endIdx := int(math.Min(float64(startIdx + pageSize), float64(len(s.Events))))
		selectScreen := &Selector[data.Event]{
			ScreenTitle: "Select Event",
			Next:        s.AddEventScreen,
			Options:     s.Events[startIdx:endIdx],
			HandleSelect: func(e data.Event) {
				s.AddEventScreen.newEvent = e
			},
			Formatter: util.FormatEventsShort,
		}
		return selectScreen
	case deleteEventSearchResultEvent:
		startIdx := s.page * pageSize
		endIdx := int(math.Min(float64(startIdx + pageSize), float64(len(s.Events))))
		selectScreen := &Selector[data.Event]{
			ScreenTitle: "Delete Event",
			Next:        s,
			Options:     s.Events[startIdx : endIdx],
			HandleSelect: func(e data.Event) {
				if err := s.Cache.DeleteSavedEvent(e); err != nil {
					output.Displayf("Failed to delete event: %v\n", err)
				}
			},
			Formatter: util.FormatEventsShort,
		}
		return selectScreen
	case allSavedEvents:
		s.page = 0
		return nil
	}
	return s
}

func (s EventSearchResult) numPages() int {
	return int(math.Ceil(float64(len(s.Events)) / float64(pageSize)))
}
