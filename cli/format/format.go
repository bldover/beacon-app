package format

import (
	"concert-manager/data"
	"fmt"
	"slices"
	"strings"
)

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
	futureDate := data.ValidFutureDate(e.Date)
	if futureDate {
		format += "\tPurchased: %v\n"
		fmtParts = append(fmtParts, e.Purchased)
	}

	return fmt.Sprintf(format, fmtParts...)
}

// adds leading zeros if needed
func FormatDate(date string) string {
	parts := strings.Split(date, "/")
	month := parts[0]
	day := parts[1]
	year := parts[2]
	if len(month) == 1 {
		month = "0" + month
	}
	if len(day) == 1 {
		day = "0" + day
	}
	return fmt.Sprintf("%s/%s/%s", month, day, year)
}

func FormatEventShort(e data.Event, maxNameLen int) string {
	var artist string
	if e.MainAct.Populated() {
		artist = e.MainAct.Name
	} else {
		artist = e.Openers[0].Name
	}

	date := FormatDate(e.Date)
	var spacing strings.Builder
	for i := len(artist); i < maxNameLen; i++ {
		spacing.WriteString(" ")
	}
	return fmt.Sprintf("%s %s%v @ %s", artist, spacing.String(), date, e.Venue.Name)
}

func FormatEventExpanded(e data.Event, future bool) string {
	mainActFmt := "Main Act: %+v"
	mainActNaFmt := "Main Act: N/A"
	openerFmt := "Openers: %s"
	openerNaFmt := "Openers: N/A"
	venueFmt := "Venue: %s"
	venueNaFmt := "Venue: N/A"
	dateFmt := "Date: %s"
	dateNaFmt := "Date: N/A"
	purchasedFmt := "Purchased: %v"

	mainAct := mainActNaFmt
	if e.MainAct.Populated() {
		mainAct = fmt.Sprintf(mainActFmt, e.MainAct)
	}

	openers := openerNaFmt
	if len(e.Openers) > 0 {
		allOpeners := strings.Builder{}
		for _, op := range e.Openers {
			allOpeners.WriteString(fmt.Sprintf("\n\t%+v", op))
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
