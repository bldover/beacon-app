package finder

import (
	"concert-manager/data"
	"concert-manager/log"
	"errors"
	"strings"
)

type eventRetriever interface {
    GetUpcomingEvents(string, string) ([]data.EventDetails, error)
}

type EventFinder struct {
	Ticketmaster eventRetriever
}

func NewEventFinder() *EventFinder {
	finder := EventFinder{}
	return &finder
}

func (f EventFinder) FindAllEvents(city string, state string) ([]data.EventDetails, error) {
	anyError := false
	events, err := f.Ticketmaster.GetUpcomingEvents(city, state)
	if err != nil {
		log.Error("Failed to retrieve all events from Ticketmaster", err)
		anyError = true
	}
	log.Debug("Total retrieved event count:", len(events))

	postProcess(events)

	if anyError {
		return events, errors.New("some events were unable to be retrieved")
	}
	return events, nil
}

// venues sometimes have weird names from non-partnered ticketing sites
func postProcess(events []data.EventDetails) {
    for i := range events {
		event := events[i].Event
		venue := event.Venue.Name
		if strings.Contains(venue, "Eastern") {
			event.Venue.Name = "The Eastern"
		} else if strings.Contains(venue, "Cadence") {
			event.Venue.Name = "Cadence Bank Ampitheatre"
		} else if strings.Contains(venue, "Altar") {
			event.Venue.Name = "The Masquerade - Altar"
		}
	}
}
