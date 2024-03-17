package finder

import (
	"concert-manager/data"
	"concert-manager/finder/ticketmaster"
	"concert-manager/log"
	"errors"
	"strings"
)

type FindEventRequest struct {
    City string
	State string
}

func FindAllEvents(request FindEventRequest) ([]data.EventDetails, error) {
	anyError := false
	allEvents := []data.EventDetails{}

	ticketmasterRequest := ticketmaster.UpcomingEventsRequest{
		City: request.City,
		State: request.State,
	}
    ticketmasterEvents, err := ticketmaster.GetUpcomingEvents(ticketmasterRequest)
	if err != nil {
		log.Error("Failed to retrieve all events from Ticketmaster", err)
		anyError = true
	}
	allEvents = append(allEvents, ticketmasterEvents...)

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
