package util

import (
	"concert-manager/data"
	"concert-manager/log"
	"fmt"
	"math"
	"slices"
	"strings"
)

func FormatArtist(artist data.Artist) string {
	artistFmt := "%s - %s"
    return fmt.Sprintf(artistFmt, artist.Name, artist.Genre)
}

func FormatArtists(artists []data.Artist) []string {
    formattedArtists := []string{}
	for _, artist := range artists {
		formattedArtists = append(formattedArtists, FormatArtist(artist))
	}
	return formattedArtists
}

func FormatArtistExpanded(artist data.Artist) string {
    artistFmt := "Name: %s\nGenre: %s"
	return fmt.Sprintf(artistFmt, artist.Name, artist.Genre)
}

func FormatVenue(venue data.Venue) string {
    venueFmt := "%s - %s, %s"
	return fmt.Sprintf(venueFmt, venue.Name, venue.City, venue.State)
}

func FormatVenues(venues []data.Venue) []string {
    formattedVenues := []string{}
	for _, venue := range venues {
		formattedVenues = append(formattedVenues, FormatVenue(venue))
	}
	return formattedVenues
}

func FormatVenueExpanded(venue data.Venue) string {
    venueFmt := "Name: %s\nCity: %s\nState: %s"
	return fmt.Sprintf(venueFmt, venue.Name, venue.City, venue.State)
}

func FormatEvent(e data.Event) string {
	fmtParts := []any{}

	date := FormatDate(e.Date)
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
	if FutureDate(e.Date) {
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
		date := FormatDate(event.Date)
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

func FormatEventExpanded(e data.Event) string {
	eventFmt := "%s\n%s\n%s\n%s\n%s\n"
	mainActFmt := "Main Act: %s"
	mainActNaFmt := "Main Act: N/A"
	openerFmt := "Openers:  %s"
	openerNaFmt := "Openers:  N/A"
	venueFmt := "Venue: %s"
	venueNaFmt := "Venue: N/A"
	dateFmt := "Date: %s"
	dateNaFmt := "Date: N/A"
	purchasedFmt := "Purchased: %v"

	mainAct := mainActNaFmt
	if e.MainAct.Name != "" {
		mainAct = fmt.Sprintf(mainActFmt, FormatArtist(e.MainAct))
	}

	openers := openerNaFmt
	if len(e.Openers) == 1 {
		opener := FormatArtist(e.Openers[0])
		openers = fmt.Sprintf(openerFmt, opener)
	} else if len(e.Openers) > 1 {
		allOpeners := strings.Builder{}
		opener := FormatArtist(e.Openers[0])
		allOpeners.WriteString(opener)
		for _, op := range e.Openers[1:] {
			opener := FormatArtist(op)
			allOpeners.WriteString("\n          " + opener)
		}
		openers = fmt.Sprintf(openerFmt, allOpeners.String())
	}

	venue := venueNaFmt
	if e.Venue.Populated() {
		venue = fmt.Sprintf(venueFmt, FormatVenueExpanded(e.Venue))
	}

	date := dateNaFmt
	if ValidDate(e.Date) {
		date = fmt.Sprintf(dateFmt, FormatDate(e.Date))
	}

	purchased := fmt.Sprintf(purchasedFmt, e.Purchased)

	fmtParts := []any{}
	fmtParts = append(fmtParts, mainAct)
	fmtParts = append(fmtParts, openers)
	fmtParts = append(fmtParts, venue)
	fmtParts = append(fmtParts, date)
	fmtParts = append(fmtParts, purchased)

    return fmt.Sprintf(eventFmt, fmtParts...)
}

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

func FormatEventRank(e data.EventRank) string {
	fmtParts := []any{}
	event := e.Event.Event

	artists := []string{}
	genres := []string{}
	if event.MainAct.Name != "" {
		artists = append(artists, event.MainAct.Name)
		mainActGenre := event.MainAct.Genre
		if mainActGenre == "" {
			mainActGenre = e.Event.EventGenre
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
					openerGenre = e.Event.EventGenre
				}
				if !slices.Contains(genres, openerGenre) && openerGenre != "" {
					genres = append(genres, openerGenre)
				}
			}
		}
		artistStr := strings.Join(artists, ", ")
		fmtParts = append(fmtParts, "Artists", artistStr)
	} else {
		eventName := e.Event.Name
		genres = append(genres, e.Event.EventGenre)
		fmtParts = append(fmtParts, "Event", eventName)
	}
	genreStr := strings.Join(genres, ", ")
	if genreStr == "" {
		genreStr = "Unknown"
	}
	fmtParts = append(fmtParts, genreStr)

	fmtParts = append(fmtParts, event.Venue.Name)

	price := e.Event.Price
	fmtParts = append(fmtParts, price)

	fmtParts = append(fmtParts, e.Rank)

	format := "%s: %s\n\tGenres: %s\n\tVenue: %s\n\tPrice: %v\n\tRank: %v\n"

	similar := []string{}
	for _, artistRank := range e.ArtistRanks {
		for _, relatedArtist := range artistRank.Related {
			if matches := SearchStrings(relatedArtist, artists, 1, ExactTolerance); len(matches) == 0 {
				similar = append(similar, relatedArtist)
			} else {
				log.Debugf("Hiding similar artist %s due to being part of the event", relatedArtist)
			}
		}
	}
	if len(similar) > 0 {
		slices.Sort(similar)
		similar = slices.Compact(similar)
		similarStr := strings.Join(similar, ", ")
		fmtParts = append(fmtParts, similarStr)
		format += "\tSimilar: %s\n"
	}

	return fmt.Sprintf(format, fmtParts...)
}
