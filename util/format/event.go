package format

import (
	"concert-manager/data"
	"fmt"
	"math"
	"slices"
	"strings"
)

func FormatEvent(e data.Event) string {
	fmtParts := []any{}

	date := formatDate(e.Date)
	fmtParts = append(fmtParts, date)

	location := fmt.Sprintf("%s, %s, %s", e.Venue.Name, e.Venue.City, e.Venue.State)
	fmtParts = append(fmtParts, location)

	artists := []string{}
	genres := []string{}
	if e.MainAct.Populated() {
		artists = append(artists, e.MainAct.Name)
		genres = append(genres, e.MainAct.Genre)
	}
	for _, artist := range e.Openers {
		if artist.Populated() {
			if !slices.Contains(artists, artist.Name) {
				artists = append(artists, artist.Name)
			}
			if !slices.Contains(genres, artist.Genre) {
				genres = append(genres, artist.Genre)
			}
		}
	}
	artistStr := strings.Join(artists, ", ")
	genreStr := strings.Join(genres, ", ")
	fmtParts = append(fmtParts, artistStr, genreStr)

	format := "%v @ %s\n\tArtists: %s\n\tGenres: %s\n"
	if data.ValidFutureDate(e.Date) {
		format += "\tPurchased: %v\n"
		fmtParts = append(fmtParts, e.Purchased)
	}

	return fmt.Sprintf(format, fmtParts...)
}

func FormatEventsShort(events []data.Event) []string {
	artistNames := []string{}
	maxNameLen := 0
	for _, event := range events {
		var artist string
		if event.MainAct.Populated() {
			artist = event.MainAct.Name
		} else {
			artist = event.Openers[0].Name
		}
		artistNames = append(artistNames, artist)
		maxNameLen = int(math.Max(float64(maxNameLen), float64(len(artist))))
	}

	formattedEvents := []string{}
	for i, event := range events {
		artist := artistNames[i]
		date := event.Date
		venue := event.Venue.Name

		var spacing strings.Builder
		for i := len(artist); i < maxNameLen; i++ {
			spacing.WriteString(" ")
		}

		formattedEvent := fmt.Sprintf("%s %s%v @ %s", artist, spacing.String(), date, venue)
		formattedEvents = append(formattedEvents, formattedEvent)
	}
	return formattedEvents
}

func FormatEventExpanded(e data.Event, future bool) string {
	mainActFmt := "Main Act: %+v"
	mainActNaFmt := "Main Act: N/A"
	openerFmt := "Openers: %s"
	openerNaFmt := "Openers: N/A"
	venueFmt := "Venue: %+v"
	venueNaFmt := "Venue: N/A"
	dateFmt := "Date: %s"
	dateNaFmt := "Date: N/A"
	purchasedFmt := "Purchased: %v"

	mainAct := mainActNaFmt
	if e.MainAct.Populated() {
		mainAct = fmt.Sprintf(mainActFmt, e.MainAct)
	}

	openers := openerNaFmt
	if len(e.Openers) == 1 {
		opener := fmt.Sprintf("%+v", e.Openers[0])
		openers = fmt.Sprintf(openerFmt, opener)
	} else if len(e.Openers) > 1 {
		allOpeners := strings.Builder{}
		allOpeners.WriteString(fmt.Sprintf("%+v", e.Openers[0]))
		for _, op := range e.Openers[1:] {
			allOpeners.WriteString("\n         ")
			allOpeners.WriteString(fmt.Sprintf("%+v", op))
		}
		openers = fmt.Sprintf(openerFmt, allOpeners.String())
	}

	venue := venueNaFmt
	if e.Venue.Populated() {
		venue = fmt.Sprintf(venueFmt, e.Venue)
	}

	date := dateNaFmt
	if data.ValidDate(e.Date) {
		date = fmt.Sprintf(dateFmt, e.Date)
	}

	fmtParts := []any{}
	fmtParts = append(fmtParts, mainAct)
	fmtParts = append(fmtParts, openers)
	fmtParts = append(fmtParts, venue)
	fmtParts = append(fmtParts, date)

	finalFmt := "%s\n%s\n%s\n%s\n"
	if future {
		finalFmt = "%s\n%s\n%s\n%s\n%s\n"
		purchased := fmt.Sprintf(purchasedFmt, e.Purchased)
		fmtParts = append(fmtParts, purchased)
	}

    return fmt.Sprintf(finalFmt, fmtParts...)
}
