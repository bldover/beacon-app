package util

import (
	"concert-manager/data"
	"slices"
)

func CloneArtist(artist data.Artist) data.Artist {
	clone := artist
	clone.Genres = CloneGenreInfo(artist.Genres)
    return clone
}

func CloneArtists(artists []data.Artist) []data.Artist {
	clone := []data.Artist{}
	for _, artist := range artists {
		clone = append(clone, CloneArtist(artist))
	}
	return clone
}

func CloneGenreInfo(genres data.GenreInfo) data.GenreInfo {
    clone := data.GenreInfo{LfmGenres: []string{}, UserGenres: []string{}}
	if genres.LfmGenres != nil {
		clone.LfmGenres = append(clone.LfmGenres, genres.LfmGenres...)
	}
	if genres.UserGenres != nil {
		clone.UserGenres = append(clone.UserGenres, genres.UserGenres...)
	}
	return clone
}

// gives us easy support to add (and then clone) complex fields in the future
func CloneVenue(venue data.Venue) data.Venue {
    return venue
}

func CloneVenues(venues []data.Venue) []data.Venue {
	clone := []data.Venue{}
	for _, venue := range venues {
		clone = append(clone, CloneVenue(venue))
	}
	return clone
}

func CloneEvent(event data.Event) data.Event {
	clone := event
	clone.Openers = slices.Clone(event.Openers)
	return clone
}

func CloneEvents(events []data.Event) []data.Event {
    clone := []data.Event{}
	for _, event := range events {
		clone = append(clone, CloneEvent(event))
	}
	return clone
}

func CloneEventDetail(event data.EventDetails) data.EventDetails {
	clone := event
	clone.Event.Openers = slices.Clone(event.Event.Openers)
	ranksClone := CloneRankInfo(*event.Ranks)
	clone.Ranks = &ranksClone
	return clone
}

func CloneEventDetails(events []data.EventDetails) []data.EventDetails {
    clone := []data.EventDetails{}
	for _, event := range events {
		clone = append(clone, CloneEventDetail(event))
	}
	return clone
}

func CloneRankInfo(rankInfo data.RankInfo) data.RankInfo {
    clone := rankInfo
	clone.ArtistRanks = CloneArtistRanks(rankInfo.ArtistRanks)
	return clone
}

func CloneArtistRanks(artists map[string]data.ArtistRank) map[string]data.ArtistRank {
	clone := map[string]data.ArtistRank{}
	for artist, rank := range artists {
		clone[artist] = CloneArtistRank(rank)
	}
	return clone
}

func CloneArtistRank(artist data.ArtistRank) data.ArtistRank {
	clone := artist
	clone.Related = slices.Clone(artist.Related)
	return clone
}
