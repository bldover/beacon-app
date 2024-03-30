package util

import (
	"concert-manager/data"
	"fmt"
	"math"
	"slices"
	"strings"
)

func FormatEventDetails(d data.EventDetails) string {
	fmtParts := []any{}
	event := d.Event

	date := FormatDate(event.Date)
	fmtParts = append(fmtParts, date)

	fmtParts = append(fmtParts, event.Venue.Name)

	artists := []string{}
	genres := []string{}
	if event.MainAct.Name != "" {
		artists = append(artists, event.MainAct.Name)
		mainActGenre := event.MainAct.Genre
		if mainActGenre == "" {
			mainActGenre = d.EventGenre
		}
		if mainActGenre != "" {
			genres = append(genres, mainActGenre)
		}
		for _, artist := range event.Openers {
			if artist.Name != "" {
				if !slices.Contains(artists, artist.Name) {
					artists = append(artists, artist.Name)
				}
				openerGenre := artist.Genre
				if openerGenre == "" {
					openerGenre = d.EventGenre
				}
				if !slices.Contains(genres, openerGenre) && openerGenre != "" {
					genres = append(genres, openerGenre)
				}
			}
		}
		artistStr := strings.Join(artists, ", ")
		fmtParts = append(fmtParts, "Artists", artistStr)
	} else {
		eventName := d.Name
		genres = append(genres, d.EventGenre)
		fmtParts = append(fmtParts, "Event", eventName)
	}
	genreStr := strings.Join(genres, ", ")
	if genreStr == "" {
		genreStr = "Unknown"
	}
	fmtParts = append(fmtParts, genreStr)

	price := d.Price
	fmtParts = append(fmtParts, price)

	format := "%v @ %s\n\t%s: %s\n\tGenres: %s\n\tPrice: %v\n"
	return fmt.Sprintf(format, fmtParts...)
}

func FormatEventDetailsShort(details []data.EventDetails) []string {
	eventNames := []string{}
	maxNameLen := 0
	for _, detail := range details {
		var eventName string
		if detail.Event.MainAct.Name != "" {
			eventName = detail.Event.MainAct.Name
		} else {
			eventName = detail.Name
		}
		eventNames = append(eventNames, eventName)
		maxNameLen = int(math.Max(float64(maxNameLen), float64(len(eventName))))
	}

	formattedEvents := []string{}
	for i, detail := range details {
		eventName := eventNames[i]
		date := FormatDate(detail.Event.Date)
		venue := detail.Event.Venue.Name

		var spacing strings.Builder
		for i := len(eventName); i < maxNameLen; i++ {
			spacing.WriteString(" ")
		}

		formattedEvent := fmt.Sprintf("%s %s%v @ %s", eventName, spacing.String(), date, venue)
		formattedEvents = append(formattedEvents, formattedEvent)
	}
	return formattedEvents
}

func EventDetailsSorterDateAsc() func(a, b data.EventDetails) int {
	return func(a, b data.EventDetails) int {
		return EventSorterDateAsc()(a.Event, b.Event)
	}
}

func EventDetailsSorterDateDesc() func(a, b data.EventDetails) int {
	return func(a, b data.EventDetails) int {
		return EventSorterDateDesc()(a.Event, b.Event)
	}
}
