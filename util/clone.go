package util

import (
	"concert-manager/data"
	"slices"
)

// gives us easy support to add (and then clone) complex fields in the future
func CloneArtist(artist data.Artist) data.Artist {
    return artist
}

func CloneArtists(artists []data.Artist) []data.Artist {
	clone := []data.Artist{}
	for _, artist := range artists {
		clone = append(clone, CloneArtist(artist))
	}
	return clone
}

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
	return clone
}

func CloneEventDetails(events []data.EventDetails) []data.EventDetails {
    clone := []data.EventDetails{}
	for _, event := range events {
		clone = append(clone, CloneEventDetail(event))
	}
	return clone
}

func CloneArtistRank(artist data.ArtistRank) data.ArtistRank {
	clone := artist
	clone.Artist = CloneArtist(artist.Artist)
	clone.Related = slices.Clone(artist.Related)
	return clone
}

func CloneArtistRanks(artists []data.ArtistRank) []data.ArtistRank {
	clone := []data.ArtistRank{}
	for _, artist := range artists {
		clone = append(clone, CloneArtistRank(artist))
	}
	return clone
}

func CloneEventRank(event data.EventRank) data.EventRank {
	clone := event
	clone.Event = CloneEventDetail(event.Event)
	clone.ArtistRanks = CloneArtistRanks(event.ArtistRanks)
	return clone
}

func CloneEventRanks(events []data.EventRank) []data.EventRank {
    clone := []data.EventRank{}
	for _, event := range events {
		clone = append(clone, CloneEventRank(event))
	}
	return clone
}
