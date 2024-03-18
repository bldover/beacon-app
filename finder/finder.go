package finder

import (
	"concert-manager/data"
	"concert-manager/log"
	"errors"
	"strings"
)

type FindEventRequest struct {
    City string
	State string
}

type EventRetriever interface {
    GetUpcomingEvents(FindEventRequest) ([]data.EventDetails, error)
}

type EventFinder struct {
    retrievers map[string]EventRetriever
}

func NewEventFinder() *EventFinder {
	finder := EventFinder{}
	finder.retrievers = map[string]EventRetriever{}
	finder.retrievers["Ticketmaster"] = ticketmasterRetriever{}
	return &finder
}

func (finder EventFinder) FindAllEvents(request FindEventRequest) ([]data.EventDetails, error) {
	anyError := false
	allEvents := []data.EventDetails{}

	for name, retriever := range finder.retrievers {
		events, err := retriever.GetUpcomingEvents(request)
		if err != nil {
			log.Error("Failed to retrieve all events from", name, err)
			anyError = true
		}
		allEvents = append(allEvents, events...)
	}
	log.Debug("Total retrieved event count:", len(allEvents))

	postProcess(allEvents)

	if anyError {
		return allEvents, errors.New("some events were unable to be retrieved")
	}
	return allEvents, nil
}

// venues sometimes have weird names from non-partnered ticketing sites
func postProcess(events []data.EventDetails) {
    for i := range events {
		event := events[i].Event
		venue := event.Venue.Name
		if strings.Contains(venue, "The Eastern") {
			event.Venue.Name = "The Eastern"
		} else if strings.Contains(venue, "Cadence Bank") {
			event.Venue.Name = "Cadence Bank Ampitheatre"
		} else if strings.Contains(venue, "Altar") {
			event.Venue.Name = "The Masquerade - Altar"
		}
	}
}
