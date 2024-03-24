package search

import (
	"concert-manager/data"
	"math"
	"slices"
)

type Search struct {
    Cache Cache
}

type Cache interface {
    GetArtists() []data.Artist
	GetVenues() []data.Venue
}

const threshold = 0.3
const maxCount = 10

func (s Search) FindFuzzyArtistMatches(artist data.Artist) []data.Artist {
	return s.FindFuzzyArtistMatchesByName(artist.Name)
}

func (s Search) FindFuzzyArtistMatchesByName(name string) []data.Artist {
	eligibleMatches := []data.Artist{}
	maxDistance := int(math.Ceil(float64(len(name)) * threshold))

	existingArtists := s.Cache.GetArtists()
	distances := map[data.Artist]int{}
	for _, artist := range existingArtists {
		distance := getLevenshteinDistance(name, artist.Name)
		if distance <= maxDistance {
			distances[artist] = distance
			eligibleMatches = append(eligibleMatches, artist)
		}
	}
	slices.SortFunc(eligibleMatches, compareDistancesAsc[data.Artist](distances))
	returnCount := int(math.Min(float64(len(eligibleMatches)), float64(maxCount)))
	return eligibleMatches[:returnCount]
}

func (s Search) FindFuzzyVenueMatches(venue data.Venue) []data.Venue {
    return s.FindFuzzyVenueMatchesByName(venue.Name)
}

func (s Search) FindFuzzyVenueMatchesByName(name string) []data.Venue {
	eligibleMatches := []data.Venue{}
	maxDistance := int(math.Ceil(float64(len(name)) * threshold))

	existingVenues := s.Cache.GetVenues()
	distances := map[data.Venue]int{}
	for _, venue := range existingVenues {
		distance := getLevenshteinDistance(name, venue.Name)
		if distance <= maxDistance {
			distances[venue] = distance
			eligibleMatches = append(eligibleMatches, venue)
		}
	}
	slices.SortFunc(eligibleMatches, compareDistancesAsc[data.Venue](distances))
	returnCount := int(math.Min(float64(len(eligibleMatches)), float64(maxCount)))
	return eligibleMatches[:returnCount]
}

func compareDistancesAsc[T comparable](distances map[T]int) func(T, T) int {
	return func(x, y T) int {
		if distances[x] > distances[y] {
			return 1
		} else if distances[x] < distances[y] {
			return -1
		} else {
			return 0
		}
	}
}
