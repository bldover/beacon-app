package screens

import (
	"concert-manager/cache"
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/ui/input"
	"concert-manager/ui/output"
	"concert-manager/util"
	"fmt"
	"strings"
	"time"
)

type recommendationCache interface {
	GetRecommendedEvents(cache.Threshold) []data.EventRank
	ChangeLocation(string, string)
	GetLocation() cache.Location
	Invalidate()
}

type savedEventCache interface {
	GetSavedEvents() []data.Event
}

type RecommendationViewer struct {
	AddEventScreen      *EventAdder
	RecommendationCache recommendationCache
	SavedCache          savedEventCache
	actions             []string
	date                time.Time
	recs                map[string][]data.EventRank
	firstRecDate        time.Time
	lastRecDate         time.Time
	threshold           cache.Threshold
}

const (
	nextDate = iota + 1
	prevDate
	gotoDate
	saveRecEvent
	changeThreshold
	changeRecLocation
	refreshRecommendations
	recToDiscoveryMenu
)

func NewRecommendationScreen() *RecommendationViewer {
	view := RecommendationViewer{}
	view.actions = []string{"Next Date", "Prev Date", "Goto Date", "Save Event", "Change Recommendation Threshold", "Change Location", "Refresh Events", "Discovery Menu"}
	view.threshold = cache.LowThreshold
	return &view
}

func (v RecommendationViewer) Title() string {
	return "Recommended Events"
}

func (v *RecommendationViewer) Refresh() {
	output.Displayf("Retrieving recommendations for %s...", v.RecommendationCache.GetLocation())
	events := v.RecommendationCache.GetRecommendedEvents(v.threshold)
	log.Debugf("Found %v recommendations for threshold %s\n", len(events), v.threshold.Level())
	v.recs = map[string][]data.EventRank{}
	for _, e := range events {
		date := util.Timestamp(e.Event.Event.Date)
		key := getRecKey(date)
		eventsForDate := v.recs[key]
		if eventsForDate == nil {
			eventsForDate = []data.EventRank{}
		}
		eventsForDate = append(eventsForDate, e)
		v.recs[key] = eventsForDate
	}

	firstDate, lastDate := time.Time{}.AddDate(9999, 0, 0), time.Time{}
	for _, e := range events {
		eventDate := util.Timestamp(e.Event.Event.Date)
		if eventDate.Before(firstDate) {
			firstDate = eventDate
		}
		if eventDate.After(lastDate) {
			lastDate = eventDate
		}
	}
	v.firstRecDate = firstDate
	v.lastRecDate = lastDate
	if v.date.IsZero() {
		v.date = firstDate
	}
	log.Debugf("Updating rec dates as firstRec: %s, lastRec: %s\n", firstDate, lastDate)
	output.ClearCurrentLine()
}

func (v RecommendationViewer) DisplayData() {
	if v.recs == nil {
		v.Refresh()
	}

	var eventData strings.Builder
	eventData.WriteString(fmt.Sprintf("Filter Threshold: %s\n", v.threshold.Level()))

	weekday := v.date.Weekday().String()
	formattedDate := util.Date(v.date)
	dateInd := fmt.Sprintf("Date - %s, %s\n", weekday, formattedDate)
	eventData.WriteString(dateInd)

	eventData.WriteString("--Saved Events--\n")
	savedEvents := v.getSavedEventsForDate(v.date)
	for _, event := range savedEvents {
		eventData.WriteString(util.FormatEvent(event))
	}
	if len(savedEvents) == 0 {
		eventData.WriteString("(none)\n")
	}
	eventData.WriteString("\n")

	eventData.WriteString("--Recommended Events--\n")
	recs := v.recs[getRecKey(v.date)]
	if recs == nil {
		recs = []data.EventRank{}
	}
	nonSavedRecs := getNonSavedEvents(recs, savedEvents)
	for _, rec := range nonSavedRecs {
		eventData.WriteString(util.FormatEventRank(rec))
	}
	if len(nonSavedRecs) == 0 {
		eventData.WriteString("(none)\n")
	}

	output.Displayln(eventData.String())
}

func getNonSavedEvents(recs []data.EventRank, saveds []data.Event) []data.EventRank {
	nonSaved := []data.EventRank{}
	for _, rec := range recs {
		match := false
		recEvent := rec.Event.Event
		for _, saved := range saveds {
			if saved.Date == recEvent.Date && saved.Venue.Name == recEvent.Venue.Name {
				match = true
				break
			}
		}
		if !match {
			nonSaved = append(nonSaved, rec)
		}
	}
	return nonSaved
}

func (v RecommendationViewer) Actions() []string {
	return v.actions
}

func (v *RecommendationViewer) NextScreen(i int) Screen {
	switch i {
	case nextDate:
		for {
			v.date = v.date.AddDate(0, 0, 1)
			log.Debug("Next date: ", v.date)
			if v.date.After(v.lastRecDate) {
				log.Debug("Date is after lastRec date, setting date to lastRec ", v.lastRecDate)
				v.date = v.lastRecDate
			}
			if len(v.recs[getRecKey(v.date)]) > 0 {
				log.Debug("Found recommended events for date")
				break
			}
			log.Debugf("No recommended events for %s, trying next date\n", util.Date(v.date))
		}
	case prevDate:
		for {
			v.date = v.date.AddDate(0, 0, -1)
			log.Debug("Prev date: ", v.date)
			if v.date.Before(v.firstRecDate) {
				log.Debug("Date is before firstRec date, setting date to firstRec ", v.firstRecDate)
				v.date = v.firstRecDate
			}
			if len(v.recs[getRecKey(v.date)]) > 0 {
				log.Debug("Found recommended events for date")
				break
			}
			log.Debugf("No recommended events for %s, trying prev date\n", util.Date(v.date))
		}
	case gotoDate:
		newDate := input.PromptAndGetInput("date", input.DateValidation)
		v.date = util.Timestamp(newDate)
	case saveRecEvent:
		events := []data.EventDetails{}
		ranks := v.recs[getRecKey(v.date)]
		for _, rank := range ranks {
			events = append(events, rank.Event)
		}
		selectScreen := &Selector[data.EventDetails]{
			ScreenTitle: "Select Event",
			Next:        v.AddEventScreen,
			Options:     events,
			HandleSelect: func(e data.EventDetails) {
				v.AddEventScreen.newEvent = e.Event
			},
			Formatter: util.FormatEventDetailsShort,
		}
		return selectScreen
	case changeThreshold:
		const high = "High"
		const medium = "Medium"
		const low = "Low"
		selectScreen := &Selector[string]{
			ScreenTitle: "Select Recommendation Threshold",
			Next:        v,
			Options:     []string{"High", "Medium", "Low"},
			HandleSelect: func(s string) {
				switch s {
				case high:
					v.threshold = cache.HighThreshold
				case medium:
					v.threshold = cache.MediumThreshold
				case low:
					v.threshold = cache.LowThreshold
				}
				v.recs = nil
			},
			Formatter: IdentityTransform[string],
		}
		return selectScreen
	case changeRecLocation:
		v.changeLocation()
	case refreshRecommendations:
		v.RecommendationCache.Invalidate()
		v.recs = nil
	case recToDiscoveryMenu:
		v.date = time.Time{}
		return nil
	}
	return v
}

func (v RecommendationViewer) getSavedEventsForDate(date time.Time) []data.Event {
	log.Debug("Requesting saved events for date ", util.Date(date))
	events := []data.Event{}
	for _, event := range v.SavedCache.GetSavedEvents() {
		if date.Equal(util.Timestamp(event.Date)) {
			events = append(events, event)
		}
	}
	log.Debugf("Found %v saved events for date %s\n", len(events), util.Date(date))
	return events
}

func (v *RecommendationViewer) changeLocation() {
	city := input.PromptAndGetInput("city", input.OnlyLettersOrSpacesValidation)
	stateCode := input.PromptAndGetInput("state code", input.StateValidation)
	v.RecommendationCache.ChangeLocation(city, stateCode)
	v.recs = nil
}

func getRecKey(date time.Time) string {
    return util.Date(date.Round(0))
}
