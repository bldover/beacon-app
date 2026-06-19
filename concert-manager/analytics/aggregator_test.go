package analytics

import (
	"concert-manager/domain"
	"reflect"
	"testing"
)

func artist(id, name string, userGenres, spotifyGenres, lastFmGenres []string) domain.Artist {
	return domain.Artist{
		Name: name,
		ID:   domain.ID{Primary: id},
		Genres: domain.GenreInfo{
			User:    userGenres,
			Spotify: spotifyGenres,
			LastFm:  lastFmGenres,
		},
	}
}

func venue(id, name string) domain.Venue {
	return domain.Venue{Name: name, ID: domain.ID{Primary: id}}
}

func event(date string, v domain.Venue, main domain.Artist, openers ...domain.Artist) domain.Event {
	m := main
	return domain.Event{
		MainAct: &m,
		Openers: openers,
		Venue:   v,
		Date:    date,
		ID:      domain.ID{Primary: date + "|" + v.Name + "|" + main.Name},
	}
}

func keys(counts []Count) []string {
	out := make([]string, len(counts))
	for i, c := range counts {
		out[i] = c.Key
	}
	return out
}

func TestCountByYear(t *testing.T) {
	v := venue("v1", "Venue 1")
	a := artist("a1", "Artist 1", []string{"rock"}, nil, nil)
	events := []domain.Event{
		event("5/1/2024", v, a),
		event("6/2/2024", v, a),
		event("3/3/2025", v, a),
	}
	got := CountByYear(events)
	if len(got) != 2 {
		t.Fatalf("expected 2 year buckets, got %d", len(got))
	}
	if got[0].Key != "2024" || got[0].Count != 2 {
		t.Errorf("expected top year 2024 with count 2, got %+v", got[0])
	}
	if got[1].Key != "2025" || got[1].Count != 1 {
		t.Errorf("expected second year 2025 with count 1, got %+v", got[1])
	}
}

func TestCountByMonth(t *testing.T) {
	v := venue("v1", "Venue 1")
	a := artist("a1", "Artist 1", []string{"rock"}, nil, nil)
	events := []domain.Event{
		event("5/1/2024", v, a),
		event("5/15/2024", v, a),
		event("3/3/2025", v, a),
	}
	got := CountByMonth(events)
	if len(got) != 2 {
		t.Fatalf("expected 2 month buckets, got %d", len(got))
	}
	if got[0].Key != "2024-05" || got[0].Name != "May 2024" || got[0].Count != 2 {
		t.Errorf("expected May 2024 first, got %+v", got[0])
	}
	if got[1].Key != "2025-03" || got[1].Name != "March 2025" {
		t.Errorf("expected March 2025 second, got %+v", got[1])
	}
}

func TestCountByArtistIncludesOpeners(t *testing.T) {
	v := venue("v1", "Venue 1")
	headliner := artist("h1", "Headliner", nil, nil, nil)
	opener := artist("o1", "Opener", nil, nil, nil)
	events := []domain.Event{
		event("5/1/2024", v, headliner, opener),
		event("6/2/2024", v, headliner),
	}
	got := CountByArtist(events)
	if len(got) != 2 {
		t.Fatalf("expected 2 artist buckets, got %d", len(got))
	}
	if got[0].Key != "h1" || got[0].Count != 2 {
		t.Errorf("expected headliner first with count 2, got %+v", got[0])
	}
	if got[1].Key != "o1" || got[1].Count != 1 {
		t.Errorf("expected opener with count 1, got %+v", got[1])
	}
}

func TestCountByVenue(t *testing.T) {
	a := artist("a1", "Artist 1", nil, nil, nil)
	v1 := venue("v1", "Bravo")
	v2 := venue("v2", "Alpha")
	events := []domain.Event{
		event("5/1/2024", v1, a),
		event("6/2/2024", v2, a),
		event("7/3/2024", v2, a),
	}
	got := CountByVenue(events)
	if len(got) != 2 {
		t.Fatalf("expected 2 venue buckets, got %d", len(got))
	}
	if got[0].Key != "v2" || got[0].Count != 2 {
		t.Errorf("expected v2 first with count 2, got %+v", got[0])
	}
	if got[1].Key != "v1" || got[1].Count != 1 {
		t.Errorf("expected v1 with count 1, got %+v", got[1])
	}
}

func TestCountByGenreDedupsPerEvent(t *testing.T) {
	v := venue("v1", "Venue 1")
	a := artist("a1", "Artist 1", []string{"Rock", "Indie"}, nil, nil)
	b := artist("a2", "Artist 2", []string{"rock"}, nil, nil)
	events := []domain.Event{
		event("5/1/2024", v, a, b),
		event("6/2/2024", v, a),
	}
	got := CountByGenre(events)
	if len(got) != 2 {
		t.Fatalf("expected 2 genre buckets, got %d", len(got))
	}
	rockIdx := -1
	indieIdx := -1
	for i, c := range got {
		if c.Key == "rock" {
			rockIdx = i
		}
		if c.Key == "indie" {
			indieIdx = i
		}
	}
	if rockIdx == -1 || indieIdx == -1 {
		t.Fatalf("missing rock or indie bucket: %+v", got)
	}
	if got[rockIdx].Count != 2 {
		t.Errorf("expected rock count 2 (deduped within event 1), got %d", got[rockIdx].Count)
	}
	if got[indieIdx].Count != 2 {
		t.Errorf("expected indie count 2, got %d", got[indieIdx].Count)
	}
}

func TestGenreHierarchyPrefersUser(t *testing.T) {
	v := venue("v1", "Venue 1")
	a := artist("a1", "Artist 1", []string{"jazz"}, []string{"spotify-genre"}, []string{"lastfm-genre"})
	got := CountByGenre([]domain.Event{event("5/1/2024", v, a)})
	if len(got) != 1 || got[0].Key != "jazz" {
		t.Errorf("expected user genre to win, got %+v", got)
	}
}

func TestGenreHierarchyFallsBackToSpotifyThenLastFm(t *testing.T) {
	v := venue("v1", "Venue 1")
	a := artist("a1", "Artist 1", nil, []string{"spot"}, []string{"last"})
	got := CountByGenre([]domain.Event{event("5/1/2024", v, a)})
	if len(got) != 1 || got[0].Key != "spot" {
		t.Errorf("expected spotify fallback, got %+v", got)
	}

	b := artist("a2", "Artist 2", nil, nil, []string{"last"})
	got = CountByGenre([]domain.Event{event("5/1/2024", v, b)})
	if len(got) != 1 || got[0].Key != "last" {
		t.Errorf("expected lastfm fallback, got %+v", got)
	}
}

func TestCountSortOrderTiesByName(t *testing.T) {
	v := venue("v1", "Venue 1")
	bravo := artist("b", "Bravo", nil, nil, nil)
	alpha := artist("a", "Alpha", nil, nil, nil)
	charlie := artist("c", "Charlie", nil, nil, nil)
	events := []domain.Event{
		event("5/1/2024", v, alpha),
		event("5/2/2024", v, bravo),
		event("5/3/2024", v, charlie),
	}
	got := CountByArtist(events)
	if !reflect.DeepEqual(keys(got), []string{"a", "b", "c"}) {
		t.Errorf("expected alphabetic tie-break, got %+v", keys(got))
	}
}

func TestBuildSummaryTopN(t *testing.T) {
	v := venue("v1", "Venue 1")
	a1 := artist("a1", "A1", nil, nil, nil)
	a2 := artist("a2", "A2", nil, nil, nil)
	a3 := artist("a3", "A3", nil, nil, nil)
	a4 := artist("a4", "A4", nil, nil, nil)
	events := []domain.Event{
		event("5/1/2024", v, a1),
		event("5/2/2024", v, a2),
		event("5/3/2024", v, a3),
		event("5/4/2024", v, a4),
	}
	summary := BuildSummary(events)
	if summary.TotalEvents != 4 {
		t.Errorf("expected total 4, got %d", summary.TotalEvents)
	}
	if len(summary.TopArtists) != TopN {
		t.Errorf("expected top %d artists, got %d", TopN, len(summary.TopArtists))
	}
}

func TestFilterEventsByArtist(t *testing.T) {
	v := venue("v1", "Venue 1")
	a := artist("a1", "A", nil, nil, nil)
	b := artist("a2", "B", nil, nil, nil)
	events := []domain.Event{
		event("5/1/2024", v, a),
		event("5/2/2024", v, b, a),
		event("5/3/2024", v, b),
	}
	got := FilterEventsByArtist(events, "a1")
	if len(got) != 2 {
		t.Errorf("expected 2 events for artist a1, got %d", len(got))
	}
}

func TestFilterEventsByGenreCaseInsensitive(t *testing.T) {
	v := venue("v1", "Venue 1")
	a := artist("a1", "A", []string{"Rock"}, nil, nil)
	b := artist("a2", "B", []string{"jazz"}, nil, nil)
	events := []domain.Event{
		event("5/1/2024", v, a),
		event("5/2/2024", v, b),
	}
	got := FilterEventsByGenre(events, "rock")
	if len(got) != 1 {
		t.Errorf("expected 1 event for genre rock, got %d", len(got))
	}
}

func TestFilterEventsByMonth(t *testing.T) {
	v := venue("v1", "Venue 1")
	a := artist("a1", "A", nil, nil, nil)
	events := []domain.Event{
		event("5/1/2024", v, a),
		event("5/15/2024", v, a),
		event("6/1/2024", v, a),
	}
	got := FilterEventsByMonth(events, "2024-05")
	if len(got) != 2 {
		t.Errorf("expected 2 May events, got %d", len(got))
	}
	got = FilterEventsByMonth(events, "bad-key")
	if len(got) != 0 {
		t.Errorf("expected 0 for malformed key, got %d", len(got))
	}
}
