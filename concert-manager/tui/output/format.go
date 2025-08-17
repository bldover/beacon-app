package output

import (
	"concert-manager/domain"
	"concert-manager/log"
	"concert-manager/search"
	"concert-manager/util"
	"fmt"
	"math"
	"slices"
	"strings"
)

func FormatArtist(artist domain.Artist) string {
	artistFmt := "%s - %v"
	return fmt.Sprintf(artistFmt, artist.Name, artist.Genres.Genres())
}

func FormatArtists(artists []domain.Artist) []string {
	formattedArtists := []string{}
	for _, artist := range artists {
		formattedArtists = append(formattedArtists, FormatArtist(artist))
	}
	return formattedArtists
}

func FormatArtistExpanded(artist domain.Artist) string {
	artistFmt := "Name: %s\nGenres: %v"
	return fmt.Sprintf(artistFmt, artist.Name, artist.Genres.Genres())
}

func FormatVenue(venue domain.Venue) string {
	venueFmt := "%s - %s, %s"
	return fmt.Sprintf(venueFmt, venue.Name, venue.City, venue.State)
}

func FormatVenues(venues []domain.Venue) []string {
	formattedVenues := []string{}
	for _, venue := range venues {
		formattedVenues = append(formattedVenues, FormatVenue(venue))
	}
	return formattedVenues
}

func FormatVenueExpanded(venue domain.Venue) string {
	venueFmt := "Name: %s\nCity: %s\nState: %s"
	return fmt.Sprintf(venueFmt, venue.Name, venue.City, venue.State)
}

func FormatEvent(e domain.Event) string {
	fmtParts := []any{}

	date := util.FormatDate(e.Date)
	fmtParts = append(fmtParts, date)

	location := fmt.Sprintf("%s, %s, %s", e.Venue.Name, e.Venue.City, e.Venue.State)
	fmtParts = append(fmtParts, location)

	artists := []string{}
	genres := []string{}
	if e.MainAct.Populated() {
		artists = append(artists, e.MainAct.Name)
		genres = append(genres, e.MainAct.Genres.Genres()...)
	}
	for _, artist := range e.Openers {
		if artist.Populated() {
			if !slices.Contains(artists, artist.Name) {
				artists = append(artists, artist.Name)
			}
			for _, genre := range artist.Genres.Genres() {
				if !slices.Contains(genres, genre) {
					genres = append(genres, genre)
				}
			}
		}
	}
	artistStr := strings.Join(artists, ", ")
	genreStr := strings.Join(genres, ", ")
	fmtParts = append(fmtParts, artistStr, genreStr)

	format := "%v @ %s\n\tArtists: %s\n\tGenres: %s\n"
	if util.FutureDate(e.Date) {
		format += "\tPurchased: %v\n"
		fmtParts = append(fmtParts, e.Purchased)
	}

	return fmt.Sprintf(format, fmtParts...)
}

func FormatEventsShort(events []domain.Event) []string {
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
		date := util.FormatDate(event.Date)
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

func FormatEventExpanded(e domain.Event) string {
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
	if e.MainAct != nil {
		mainAct = fmt.Sprintf(mainActFmt, FormatArtist(*e.MainAct))
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
	if util.ValidDate(e.Date) {
		date = fmt.Sprintf(dateFmt, util.FormatDate(e.Date))
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

func FormatEventDetails(d domain.EventDetails) string {
	fmtParts := []any{}
	event := d.Event

	date := util.FormatDate(event.Date)
	fmtParts = append(fmtParts, date)

	fmtParts = append(fmtParts, event.Venue.Name)

	artists := []string{}
	genres := []string{}
	if event.MainAct.Name != "" {
		artists = append(artists, event.MainAct.Name)
		mainActGenres := event.MainAct.Genres.Genres()
		genres = append(genres, mainActGenres...)
		for _, artist := range event.Openers {
			if artist.Name != "" {
				if !slices.Contains(artists, artist.Name) {
					artists = append(artists, artist.Name)
				}
				openerGenres := artist.Genres.Genres()
				for _, genre := range openerGenres {
					if !slices.Contains(genres, genre) {
						genres = append(genres, genre)
					}
				}
			}
		}
		artistStr := strings.Join(artists, ", ")
		fmtParts = append(fmtParts, "Artists", artistStr)
	} else {
		eventName := d.Name
		fmtParts = append(fmtParts, "Event", eventName)
	}

	if len(genres) == 0 {
		genres = append(genres, d.EventGenre)
	}
	genreStr := strings.Join(genres, ", ")
	if genreStr == "" {
		genreStr = "Unknown"
	}
	fmtParts = append(fmtParts, genreStr)

	format := "%v @ %s\n\t%s: %s\n\tGenres: %s\n"
	return fmt.Sprintf(format, fmtParts...)
}

func FormatEventDetailsShort(details []domain.EventDetails) []string {
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
		date := util.FormatDate(detail.Event.Date)
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

func FormatRankedEvent(e domain.EventDetails) string {
	fmtParts := []any{}
	event := e.Event

	artists := []string{}
	genres := []string{}
	if event.MainAct.Name != "" {
		artists = append(artists, event.MainAct.Name)
		mainActGenres := event.MainAct.Genres.Genres()
		if len(mainActGenres) == 0 {
			mainActGenres = []string{e.EventGenre}
		}
		genres = append(genres, mainActGenres...)
		for _, artist := range event.Openers {
			if artist.Name != "" {
				if !slices.Contains(artists, artist.Name) {
					artists = append(artists, artist.Name)
				}
				openerGenres := artist.Genres.Genres()
				if len(openerGenres) == 0 {
					openerGenres = []string{e.EventGenre}
				}
				for _, genre := range openerGenres {
					if !slices.Contains(genres, genre) {
						genres = append(genres, genre)
					}
				}
			}
		}
		artistStr := strings.Join(artists, ", ")
		fmtParts = append(fmtParts, "Artists", artistStr)
	} else {
		eventName := e.Name
		genres = append(genres, e.EventGenre)
		fmtParts = append(fmtParts, "Event", eventName)
	}
	genreStr := strings.Join(genres, ", ")
	if genreStr == "" {
		genreStr = "Unknown"
	}
	fmtParts = append(fmtParts, genreStr)

	fmtParts = append(fmtParts, event.Venue.Name)

	fmtParts = append(fmtParts, e.Ranks.Rank)

	format := "%s: %s\n\tGenres: %s\n\tVenue: %s\n\tRank: %v\n"

	similar := []string{}
	for _, artistRank := range e.Ranks.ArtistRanks {
		for _, relatedArtist := range artistRank.Related {
			if matches := search.SearchStrings(relatedArtist, artists, 1, search.ExactTolerance); len(matches) == 0 {
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
