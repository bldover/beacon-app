package util

import (
	"concert-manager/data"
	"math"
	"slices"
	"strings"
)

type Search struct {
    Cache cache
	maxCount int
}

type cache interface {
    GetArtists() []data.Artist
	GetVenues() []data.Venue
	GetSavedEvents() []data.Event
	GetUpcomingEvents(string, string) []data.EventDetails
}

const threshold = 0.4

func NewSearch() *Search {
	s := &Search{}
	s.resetMaxCount()
	return s
}

func (s *Search) WithMaxCount(maxCount int) {
    s.maxCount = maxCount
}

func (s *Search) resetMaxCount() {
    s.maxCount = 999999
}

func (s *Search) FindFuzzyArtistMatches(artist data.Artist) []data.Artist {
	return s.FindFuzzyArtistMatchesByName(artist.Name)
}

func (s *Search) FindFuzzyArtistMatchesByName(name string) []data.Artist {
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
	returnCount := int(math.Min(float64(len(eligibleMatches)), float64(s.maxCount)))
	s.resetMaxCount()
	return eligibleMatches[:returnCount]
}

func (s *Search) FindFuzzyVenueMatches(venue data.Venue) []data.Venue {
    return s.FindFuzzyVenueMatchesByName(venue.Name)
}

func (s *Search) FindFuzzyVenueMatchesByName(name string) []data.Venue {
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
	returnCount := int(math.Min(float64(len(eligibleMatches)), float64(s.maxCount)))
	s.resetMaxCount()
	return eligibleMatches[:returnCount]
}

func (s *Search) FindFuzzyEventMatchesByArtist(name string) []data.Event {
	eligibleMatches := []data.Event{}
	maxDistance := int(math.Ceil(float64(len(name)) * threshold))

	existingEvents := s.Cache.GetSavedEvents()
	distances := map[string]int{}
	for _, event := range existingEvents {
		artistNames := []string{event.MainAct.Name}
		for _, opener := range event.Openers {
			artistNames = append(artistNames, opener.Name)
		}

		eventDist := 99999
		for _, artistName := range artistNames {
			artistDistance := getLevenshteinDistance(name, artistName)
			eventDist = int(math.Min(float64(eventDist), float64(artistDistance)))
		}

		if eventDist <= maxDistance {
			eventKey := getEventKey(event)
			distances[eventKey] = eventDist
			eligibleMatches = append(eligibleMatches, event)
		}
	}

	slices.SortFunc(eligibleMatches, compareDistancesUsingKeyFuncAsc[data.Event, string](distances, getEventKey))
	returnCount := int(math.Min(float64(len(eligibleMatches)), float64(s.maxCount)))
	s.resetMaxCount()
	return eligibleMatches[:returnCount]
}

func (s *Search) FindFuzzyEventMatchesByVenue(venue string) []data.Event {
	eligibleMatches := []data.Event{}
	maxDistance := int(math.Ceil(float64(len(venue)) * threshold))

	existingEvents := s.Cache.GetSavedEvents()
	distances := map[string]int{}
	for _, event := range existingEvents {
		distance := getLevenshteinDistance(venue, event.Venue.Name)
		if distance <= maxDistance {
			eventKey := getEventKey(event)
			distances[eventKey] = distance
			eligibleMatches = append(eligibleMatches, event)
		}
	}

	slices.SortFunc(eligibleMatches, compareDistancesUsingKeyFuncAsc[data.Event, string](distances, getEventKey))
	returnCount := int(math.Min(float64(len(eligibleMatches)), float64(s.maxCount)))
	s.resetMaxCount()
	return eligibleMatches[:returnCount]
}

func getEventKey(event data.Event) string {
    key := strings.Builder{}
	key.WriteString(event.MainAct.Name)
	key.WriteString(event.Venue.Name)
	key.WriteString(event.Date)
	return key.String()
}

func (s *Search) FindFuzzyEventDetailsMatchesByArtist(name, city, state string) []data.EventDetails {
	eligibleMatches := []data.EventDetails{}
	maxDistance := int(math.Ceil(float64(len(name)) * threshold))

	existingEvents := s.Cache.GetUpcomingEvents(city, state)
	distances := map[string]int{}
	for _, event := range existingEvents {
		artistNames := []string{event.Event.MainAct.Name}
		for _, opener := range event.Event.Openers {
			artistNames = append(artistNames, opener.Name)
		}

		eventDist := 99999
		for _, artistName := range artistNames {
			artistDistance := getLevenshteinDistance(name, artistName)
			eventDist = int(math.Min(float64(eventDist), float64(artistDistance)))
		}

		if eventDist <= maxDistance {
			eventKey := getEventDetailsKey(event)
			distances[eventKey] = eventDist
			eligibleMatches = append(eligibleMatches, event)
		}
	}

	slices.SortFunc(eligibleMatches, compareDistancesUsingKeyFuncAsc[data.EventDetails, string](distances, getEventDetailsKey))
	returnCount := int(math.Min(float64(len(eligibleMatches)), float64(s.maxCount)))
	s.resetMaxCount()
	return eligibleMatches[:returnCount]
}

func (s *Search) FindFuzzyEventDetailsMatchesByVenue(venue, city, state string) []data.EventDetails {
	eligibleMatches := []data.EventDetails{}
	maxDistance := int(math.Ceil(float64(len(venue)) * threshold))

	existingEvents := s.Cache.GetUpcomingEvents(city, state)
	distances := map[string]int{}
	for _, event := range existingEvents {
		distance := getLevenshteinDistance(venue, event.Event.Venue.Name)
		if distance <= maxDistance {
			eventKey := getEventDetailsKey(event)
			distances[eventKey] = distance
			eligibleMatches = append(eligibleMatches, event)
		}
	}

	slices.SortFunc(eligibleMatches, compareDistancesUsingKeyFuncAsc[data.EventDetails, string](distances, getEventDetailsKey))
	returnCount := int(math.Min(float64(len(eligibleMatches)), float64(s.maxCount)))
	s.resetMaxCount()
	return eligibleMatches[:returnCount]
}

func getEventDetailsKey(event data.EventDetails) string {
	return getEventKey(event.Event)
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

func compareDistancesUsingKeyFuncAsc[T any, U comparable](distances map[U]int, keyFunc func(T) U) func(T, T) int {
	return func(x, y T) int {
		xKey := keyFunc(x)
		yKey := keyFunc(y)
		if distances[xKey] > distances[yKey] {
			return 1
		} else if distances[xKey] < distances[yKey] {
			return -1
		} else {
			return 0
		}
	}
}
