package analytics

import (
	"concert-manager/domain"
	"concert-manager/util"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Count struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Summary struct {
	TotalEvents int     `json:"totalEvents"`
	TopYears    []Count `json:"topYears"`
	TopMonths   []Count `json:"topMonths"`
	TopArtists  []Count `json:"topArtists"`
	TopVenues   []Count `json:"topVenues"`
	TopGenres   []Count `json:"topGenres"`
}

type EventsResponse struct {
	Count  int            `json:"count"`
	Events []domain.Event `json:"events"`
}

const TopN = 3

var monthNames = []string{
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

// PastEvents returns events whose date is before today.
func PastEvents(events []domain.Event) []domain.Event {
	out := []domain.Event{}
	for _, e := range events {
		if util.PastDate(e.Date) {
			out = append(out, e)
		}
	}
	return out
}

// analyticsGenres returns the effective genres for an artist using the
// user > spotify > lastFm hierarchy, matching the Android client. The
// domain GenreInfo.Genres() method uses a different priority and is left
// untouched.
func analyticsGenres(a domain.Artist) []string {
	if len(a.Genres.User) > 0 {
		return a.Genres.User
	}
	if len(a.Genres.Spotify) > 0 {
		return a.Genres.Spotify
	}
	return a.Genres.LastFm
}

func parseDate(date string) (year, month int, ok bool) {
	parts := strings.Split(date, "/")
	if len(parts) != 3 {
		return 0, 0, false
	}
	m, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, false
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, false
	}
	return y, m, true
}

type aggregator struct {
	counts       map[string]int
	displayNames map[string]string
}

func newAggregator() *aggregator {
	return &aggregator{
		counts:       map[string]int{},
		displayNames: map[string]string{},
	}
}

func (a *aggregator) add(key, name string) {
	a.counts[key]++
	if _, ok := a.displayNames[key]; !ok {
		a.displayNames[key] = name
	}
}

func (a *aggregator) toCounts() []Count {
	counts := make([]Count, 0, len(a.counts))
	for key, count := range a.counts {
		counts = append(counts, Count{
			Key:   key,
			Name:  a.displayNames[key],
			Count: count,
		})
	}
	sortCounts(counts)
	return counts
}

func sortCounts(counts []Count) {
	sort.SliceStable(counts, func(i, j int) bool {
		if counts[i].Count != counts[j].Count {
			return counts[i].Count > counts[j].Count
		}
		return counts[i].Name < counts[j].Name
	})
}

func top(counts []Count, n int) []Count {
	if len(counts) <= n {
		return counts
	}
	return counts[:n]
}

func CountByYear(events []domain.Event) []Count {
	agg := newAggregator()
	for _, e := range events {
		y, _, ok := parseDate(e.Date)
		if !ok {
			continue
		}
		key := strconv.Itoa(y)
		agg.add(key, key)
	}
	return agg.toCounts()
}

func CountByMonth(events []domain.Event) []Count {
	agg := newAggregator()
	for _, e := range events {
		y, m, ok := parseDate(e.Date)
		if !ok || m < 1 || m > 12 {
			continue
		}
		key := fmt.Sprintf("%04d-%02d", y, m)
		name := fmt.Sprintf("%s %d", monthNames[m-1], y)
		agg.add(key, name)
	}
	return agg.toCounts()
}

func CountByArtist(events []domain.Event) []Count {
	agg := newAggregator()
	for _, e := range events {
		for _, a := range e.Artists() {
			if a.ID.Primary == "" {
				continue
			}
			agg.add(a.ID.Primary, a.Name)
		}
	}
	return agg.toCounts()
}

func CountByVenue(events []domain.Event) []Count {
	agg := newAggregator()
	for _, e := range events {
		if e.Venue.ID.Primary == "" {
			continue
		}
		agg.add(e.Venue.ID.Primary, e.Venue.Name)
	}
	return agg.toCounts()
}

// CountByGenre groups past events by genre. Each genre is counted at most
// once per event, even if multiple artists in the lineup share it.
func CountByGenre(events []domain.Event) []Count {
	agg := newAggregator()
	for _, e := range events {
		seen := map[string]bool{}
		for _, artist := range e.Artists() {
			for _, genre := range analyticsGenres(artist) {
				key := strings.ToLower(strings.TrimSpace(genre))
				if key == "" || seen[key] {
					continue
				}
				seen[key] = true
				agg.add(key, genre)
			}
		}
	}
	return agg.toCounts()
}

func BuildSummary(events []domain.Event) Summary {
	return Summary{
		TotalEvents: len(events),
		TopYears:    top(CountByYear(events), TopN),
		TopMonths:   top(CountByMonth(events), TopN),
		TopArtists:  top(CountByArtist(events), TopN),
		TopVenues:   top(CountByVenue(events), TopN),
		TopGenres:   top(CountByGenre(events), TopN),
	}
}

func FilterEventsByYear(events []domain.Event, yearKey string) []domain.Event {
	out := []domain.Event{}
	for _, e := range events {
		y, _, ok := parseDate(e.Date)
		if !ok {
			continue
		}
		if strconv.Itoa(y) == yearKey {
			out = append(out, e)
		}
	}
	return out
}

func FilterEventsByMonth(events []domain.Event, monthKey string) []domain.Event {
	parts := strings.Split(monthKey, "-")
	if len(parts) != 2 {
		return []domain.Event{}
	}
	y, err := strconv.Atoi(parts[0])
	if err != nil {
		return []domain.Event{}
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return []domain.Event{}
	}
	out := []domain.Event{}
	for _, e := range events {
		ey, em, ok := parseDate(e.Date)
		if !ok {
			continue
		}
		if ey == y && em == m {
			out = append(out, e)
		}
	}
	return out
}

func FilterEventsByArtist(events []domain.Event, artistID string) []domain.Event {
	out := []domain.Event{}
	for _, e := range events {
		for _, a := range e.Artists() {
			if a.ID.Primary == artistID {
				out = append(out, e)
				break
			}
		}
	}
	return out
}

func FilterEventsByVenue(events []domain.Event, venueID string) []domain.Event {
	out := []domain.Event{}
	for _, e := range events {
		if e.Venue.ID.Primary == venueID {
			out = append(out, e)
		}
	}
	return out
}

// FilterEventsByGenre matches events containing any artist with the given
// genre, comparing case-insensitively.
func FilterEventsByGenre(events []domain.Event, genreKey string) []domain.Event {
	key := strings.ToLower(strings.TrimSpace(genreKey))
	out := []domain.Event{}
	for _, e := range events {
		matched := false
		for _, artist := range e.Artists() {
			for _, genre := range analyticsGenres(artist) {
				if strings.ToLower(strings.TrimSpace(genre)) == key {
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}
		if matched {
			out = append(out, e)
		}
	}
	return out
}
