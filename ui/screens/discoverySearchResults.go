package screens

import (
	"concert-manager/data"
	"concert-manager/ui/input"
	"concert-manager/ui/output"
	"concert-manager/util"
	"fmt"
	"math"
	"strings"
)

type DiscoverySearchResult struct {
	AddEventScreen   *EventAdder
	Events           []data.EventDetails
	actions          []string
	page             int
}

const (
	nextDiscoverySearchPage = iota + 1
	prevDiscoverySearchPage
	gotoDiscoverySearchPage
	addDiscoverySearchEvent
	allUpcomingEvents
)

func NewDiscoverySearchResultScreen() *DiscoverySearchResult {
	view := DiscoverySearchResult{}
	view.actions = []string{"Next Page", "Prev Page", "Goto Page",  "Save Event", "All Upcoming Events"}
	return &view
}

func (s DiscoverySearchResult) Title() string {
	return "Search Results"
}

func (s *DiscoverySearchResult) DisplayData() {
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
		eventData.WriteString(util.FormatEventDetails((s.Events)[i]))
	}
	output.Displayln(eventData.String())
}

func (s DiscoverySearchResult) Actions() []string {
	return s.actions
}

func (s *DiscoverySearchResult) NextScreen(i int) Screen {
	switch i {
	case nextDiscoverySearchPage:
		if (s.page + 1) < s.numPages() {
			s.page++
		}
		return s
	case prevDiscoverySearchPage:
		if s.page > 0 {
			s.page--
		}
	case gotoDiscoverySearchPage:
		s.page = input.PromptAndGetInputNumeric("page number", 1, s.numPages()+1) - 1
	case addDiscoverySearchEvent:
		startIdx := pageSize * s.page
		endIdx := int(math.Min(float64(startIdx + pageSize), float64(len(s.Events))))
		selectScreen := &Selector[data.EventDetails]{
			ScreenTitle: "Select Event",
			Next:        s.AddEventScreen,
			Options:     s.Events[startIdx:endIdx],
			HandleSelect: func(e data.EventDetails) {
				s.AddEventScreen.newEvent = e.Event
			},
			Formatter: util.FormatEventDetailsShort,
		}
		return selectScreen
	case allUpcomingEvents:
		s.page = 0
		return nil
	}
	return s
}

func (s DiscoverySearchResult) numPages() int {
	return int(math.Ceil(float64(len(s.Events)) / float64(pageSize)))
}
