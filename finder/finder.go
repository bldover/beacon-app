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

func (r FindEventRequest) GetCity() string {
    return r.City
}

func (r FindEventRequest) GetState() string {
    return r.State
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

	postProcess(allEvents)

	if anyError {
		return allEvents, errors.New("some events were unable to be retrieved")
	}
	return allEvents, nil
}

// venues sometimes have weird names from non-partnered ticketing sites
func postProcess(events []data.EventDetails) {
    for i, event := range events {
		if strings.Contains(event.Event.Venue.Name, "The Eastern") {
			events[i].Event.Venue.Name = "The Eastern"
		} else if strings.Contains(event.Event.Venue.Name, "Cadence Bank") {
			events[i].Event.Venue.Name = "Cadence Bank Ampitheatre"
		}
	}
}
