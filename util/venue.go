package util

import (
	"concert-manager/data"
	"fmt"
)

const venueFmt = "%s - %s, %s"

func FormatVenue(venues []data.Venue) []string {
    formattedVenues := []string{}
	for _, venue := range venues {
		formattedVenue := fmt.Sprintf(venueFmt, venue.Name, venue.City, venue.State)
		formattedVenues = append(formattedVenues, formattedVenue)
	}
	return formattedVenues
}
