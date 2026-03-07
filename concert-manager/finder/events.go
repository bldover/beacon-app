package finder

import (
	"concert-manager/domain"
	"concert-manager/log"
	"errors"
	"strings"
)

type eventRetriever interface {
	GetUpcomingEvents(string, string) ([]domain.EventDetails, error)
}

type EventFinder struct {
	Ticketmaster eventRetriever
}

func NewEventFinder() *EventFinder {
	finder := EventFinder{}
	return &finder
}

func (f EventFinder) FindAllEvents(city string, state string) ([]domain.EventDetails, error) {
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

func postProcess(events []domain.EventDetails) {
	for i := range events {
		// venues sometimes have weird names from non-partnered ticketing sites
		venue := events[i].Event.Venue.Name
		if strings.Contains(venue, "The Eastern-GA") {
			events[i].Event.Venue.Name = "The Eastern"
		} else if strings.Contains(venue, "The Masquerade  - Altar") {
			// probably the typo will be fixed at some point and we can remove this?
			events[i].Event.Venue.Name = "The Masquerade - Altar"
		}
	}
}
