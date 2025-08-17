package domain

import (
	"slices"
)

func CloneArtist(artist Artist) Artist {
	clone := artist
	clone.Genres = CloneGenreInfo(artist.Genres)
	return clone
}

func CloneArtists(artists []Artist) []Artist {
	clone := []Artist{}
	for _, artist := range artists {
		clone = append(clone, CloneArtist(artist))
	}
	return clone
}

func CloneGenreInfo(genres GenreInfo) GenreInfo {
	clone := GenreInfo{
		LastFm:       []string{},
		Spotify:      []string{},
		Ticketmaster: []string{},
		User:         []string{},
	}
	if genres.LastFm != nil {
		clone.LastFm = append(clone.LastFm, genres.LastFm...)
	}
	if genres.Spotify != nil {
		clone.Spotify = append(clone.Spotify, genres.Spotify...)
	}
	if genres.Ticketmaster != nil {
		clone.Ticketmaster = append(clone.Ticketmaster, genres.Ticketmaster...)
	}
	if genres.User != nil {
		clone.User = append(clone.User, genres.User...)
	}
	return clone
}

// currently everything is copied in the arg, but this provides easy support to
// add (and then clone) complex fields in the future
func CloneVenue(venue Venue) Venue {
	return venue
}

func CloneVenues(venues []Venue) []Venue {
	clone := []Venue{}
	for _, venue := range venues {
		clone = append(clone, CloneVenue(venue))
	}
	return clone
}

func CloneEvent(event Event) Event {
	clone := event
	if event.MainAct != nil {
		mainActClone := CloneArtist(*clone.MainAct)
		clone.MainAct = &mainActClone
	}
	clone.Openers = slices.Clone(event.Openers)
	return clone
}

func CloneEvents(events []Event) []Event {
	clone := []Event{}
	for _, event := range events {
		clone = append(clone, CloneEvent(event))
	}
	return clone
}

func CloneEventDetail(event EventDetails) EventDetails {
	clone := event
	clone.Event = CloneEvent(event.Event)
	if event.Ranks != nil {
		ranksClone := CloneRankInfo(*event.Ranks)
		clone.Ranks = &ranksClone
	}
	return clone
}

func CloneEventDetails(events []EventDetails) []EventDetails {
	clone := []EventDetails{}
	for _, event := range events {
		clone = append(clone, CloneEventDetail(event))
	}
	return clone
}

func CloneRankInfo(rankInfo RankInfo) RankInfo {
	clone := rankInfo
	clone.ArtistRanks = CloneArtistRanks(rankInfo.ArtistRanks)
	return clone
}

func CloneArtistRanks(artists map[string]ArtistRank) map[string]ArtistRank {
	clone := map[string]ArtistRank{}
	for artist, rank := range artists {
		clone[artist] = CloneArtistRank(rank)
	}
	return clone
}

func CloneArtistRank(artist ArtistRank) ArtistRank {
	clone := artist
	clone.Related = slices.Clone(artist.Related)
	return clone
}
