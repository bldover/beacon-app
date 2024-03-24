package search

import (
	"concert-manager/data"
	"testing"
)

type MockCache struct {
    artists []data.Artist
	venues []data.Venue
}

func (c MockCache) GetArtists() []data.Artist {
    return c.artists
}

func (c MockCache) GetVenues() []data.Venue {
    return c.venues
}

var artists = []data.Artist{
	{Name: "cat"},
    {Name: "hat"},
	{Name: "mat"},
	{Name: "bat"},
	{Name: "rat"},
	{Name: "pat"},
	{Name: "sat"},
	{Name: "at"},
	{Name: "fat"},
	{Name: "tat"},
	{Name: "dt"},
}

func TestMaxCountArtistsReturned(t *testing.T) {
	cache := MockCache{}
	cache.artists = artists
	search := Search{cache}

	resp := search.FindFuzzyArtistMatchesByName("dat")
	expectedLen := maxCount
	if len(resp) != expectedLen {
		t.Errorf("Incorrect number of returned artists, expected: %v, actual: %v", expectedLen, len(resp))
	}
}

func TestLessThanMaxCountArtistsReturned(t *testing.T) {
	cache := MockCache{}
	cache.artists = artists
	search := Search{cache}

	resp := search.FindFuzzyArtistMatchesByName("dt")
	expectedLen := 2
	if len(resp) != expectedLen {
		t.Errorf("Incorrect number of returned artists, expected: %v, actual: %v", expectedLen, len(resp))
	}
	expectedResp := []data.Artist{{Name: "dt"}, {Name: "at"}}
	if resp[0] != expectedResp[0] || resp[1] != expectedResp[1] {
		t.Errorf("Incorrect artists returned, expected: %v, actual: %v", expectedResp, resp)
	}
}

func TestNoArtistsReturned(t *testing.T) {
	cache := MockCache{}
	search := Search{cache}

	resp := search.FindFuzzyArtistMatchesByName("dat")
	expectedLen := 0
	if len(resp) != expectedLen {
		t.Errorf("Incorrect number of returned artists, expected: %v, actual: %v", expectedLen, len(resp))
	}
}
