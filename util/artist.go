package util

import (
	"concert-manager/data"
	"fmt"
)

const artistFmt = "%s - %s"

func FormatArtist(artists []data.Artist) []string {
    formattedArtists := []string{}
	for _, artist := range artists {
		formattedArtist := fmt.Sprintf(artistFmt, artist.Name, artist.Genre)
		formattedArtists = append(formattedArtists, formattedArtist)
	}
	return formattedArtists
}
