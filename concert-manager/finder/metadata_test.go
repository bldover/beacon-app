package finder

import (
	"concert-manager/domain"
	"concert-manager/external"
	"testing"
)

type MockSavedArtistCache struct {
	getArtistsFunc func() []domain.Artist
	calls          struct {
		getArtists int
	}
}

func (m *MockSavedArtistCache) GetArtists() []domain.Artist {
	m.calls.getArtists++
	return m.getArtistsFunc()
}

type MockGenreProvider struct {
	artistInfoByIdsFunc  func(ids []string) (map[string]external.ArtistInfo, error)
	artistInfoByNameFunc func(name string) (external.ArtistInfo, error)
	calls                struct {
		artistInfoByIds  map[string]int
		artistInfoByName map[string]int
	}
}

func newMockGenreProvider() *MockGenreProvider {
	return &MockGenreProvider{
		calls: struct {
			artistInfoByIds  map[string]int
			artistInfoByName map[string]int
		}{
			artistInfoByIds:  make(map[string]int),
			artistInfoByName: make(map[string]int),
		},
	}
}

func (m *MockGenreProvider) ArtistInfoByIds(ids []string) (map[string]external.ArtistInfo, error) {
	key := ""
	for _, id := range ids {
		key += id + ","
	}
	if m.calls.artistInfoByIds == nil {
		m.calls.artistInfoByIds = make(map[string]int)
	}
	m.calls.artistInfoByIds[key]++
	return m.artistInfoByIdsFunc(ids)
}

func (m *MockGenreProvider) ArtistInfoByName(name string) (external.ArtistInfo, error) {
	if m.calls.artistInfoByName == nil {
		m.calls.artistInfoByName = make(map[string]int)
	}
	m.calls.artistInfoByName[name]++
	return m.artistInfoByNameFunc(name)
}

func TestPopulateGenres(t *testing.T) {
	mockCache := &MockSavedArtistCache{}
	mockProvider := newMockGenreProvider()

	savedArtists := []domain.Artist{
		{Name: "Saved Artist", ID: domain.ID{Spotify: "saved-id"}, Genres: domain.GenreInfo{Spotify: []string{"indie"}}},
	}

	events := []domain.EventDetails{
		{
			Event: domain.Event{
				MainAct: &domain.Artist{Name: "Main Artist", ID: domain.ID{Spotify: "main-id"}},
				Openers: []domain.Artist{
					{Name: "Opener", ID: domain.ID{Spotify: ""}},
				},
			},
		},
		{
			Event: domain.Event{
				MainAct: &domain.Artist{Name: "Saved Artist", ID: domain.ID{Spotify: "main-id"}},
			},
		},
	}

	mockCache.getArtistsFunc = func() []domain.Artist {
		return savedArtists
	}

	mockProvider.artistInfoByIdsFunc = func(ids []string) (map[string]external.ArtistInfo, error) {
		return map[string]external.ArtistInfo{
			"main-id": {Id: "main-id", Genres: []string{"rock", "alternative"}},
		}, nil
	}

	mockProvider.artistInfoByNameFunc = func(name string) (external.ArtistInfo, error) {
		return external.ArtistInfo{
			Id: "opener-id", Genres: []string{"pop"},
		}, nil
	}

	finder := ArtistInfoFinder{
		ArtistCache:   mockCache,
		GenreProvider: mockProvider,
	}

	result := finder.PopulateGenres(events)

	if len(result) != 2 {
		t.Errorf("Expected 2 events, got %d", len(result))
	}
	if len(result[0].Event.MainAct.Genres.Spotify) != 2 || result[0].Event.MainAct.Genres.Spotify[0] != "rock" {
		t.Errorf("Expected main act genres [rock, alternative], got %v", result[0].Event.MainAct.Genres.Spotify)
	}
	if len(result[0].Event.Openers[0].Genres.Spotify) != 1 || result[0].Event.Openers[0].Genres.Spotify[0] != "pop" {
		t.Errorf("Expected opener genres [pop], got %v", result[0].Event.Openers[0].Genres.Spotify)
	}
	if result[0].Event.Openers[0].ID.Spotify != "opener-id" {
		t.Errorf("Expected opener ID opener-id, got %s", result[0].Event.Openers[0].ID.Spotify)
	}
	if len(result[1].Event.MainAct.Genres.Spotify) != 1 || result[1].Event.MainAct.Genres.Spotify[0] != "indie" {
		t.Errorf("Expected indie genre for second event, got %v", result[1].Event.MainAct.Genres.Spotify)
	}
}

func TestReloadGenres(t *testing.T) {
	mockCache := &MockSavedArtistCache{}
	mockProvider := newMockGenreProvider()

	artists := []domain.Artist{
		{Name: "Artist1", ID: domain.ID{Spotify: "id1"}},
		{Name: "Artist1"},
	}

	mockProvider.artistInfoByIdsFunc = func(ids []string) (map[string]external.ArtistInfo, error) {
		return map[string]external.ArtistInfo{
			"id1": {Id: "id1", Genres: []string{"rock"}},
		}, nil
	}

	mockProvider.artistInfoByNameFunc = func(name string) (external.ArtistInfo, error) {
		return external.ArtistInfo{
			Id: "id2", Genres: []string{"pop"},
		}, nil
	}

	finder := ArtistInfoFinder{
		ArtistCache:   mockCache,
		GenreProvider: mockProvider,
	}

	result, err := finder.ReloadGenres(artists)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 artists, got %d", len(result))
	}
	if result[0].ID.Spotify != "id1" {
		t.Errorf("Expected spotify ID id1, got %s", result[0].ID.Spotify)
	}
	if len(result[0].Genres.Spotify) != 1 || result[0].Genres.Spotify[0] != "rock" {
		t.Errorf("Expected genres [rock], got %v", result[0].Genres.Spotify)
	}
	if result[1].ID.Spotify != "id2" {
		t.Errorf("Expected spotify ID id2, got %s", result[1].ID.Spotify)
	}
	if len(result[1].Genres.Spotify) != 1 || result[1].Genres.Spotify[0] != "pop" {
		t.Errorf("Expected genres [pop], got %v", result[1].Genres.Spotify)
	}
}
