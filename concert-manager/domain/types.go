package domain

import "fmt"

type (
	Venue struct {
		Name  string `json:"name"`
		City  string `json:"city"`
		State string `json:"state"`
		ID    ID     `json:"id"`
	}
	Artist struct {
		Name   string    `json:"name"`
		Genre  string    // legacy, remove after migration to Genres
		Genres GenreInfo `json:"genres"`
		ID     ID        `json:"id"`
	}
	GenreInfo struct {
		Spotify      []string `json:"spotify"`
		LastFm       []string `json:"lastFm"`
		Ticketmaster []string `json:"ticketmaster"`
		User         []string `json:"user"`
	}
	ID struct {
		Primary      string `json:"primary"`
		Spotify      string `json:"spotify"`
		Ticketmaster string `json:"ticketmaster"`
		MusicBrainz  string `json:"musicbrainz"`
	}
	Event struct {
		MainAct   *Artist  `json:"mainAct"`
		Openers   []Artist `json:"openers"`
		Venue     Venue    `json:"venue"`
		Date      string   `json:"date"`
		Purchased bool     `json:"purchased"`
		ID        ID       `json:"id"`
	}
	EventDetails struct {
		Name       string    `json:"name"`
		EventGenre string    `json:"genre"`
		Event      Event     `json:"event"`
		Ranks      *RankInfo `json:"ranks"`
	}
	RankInfo struct {
		Rank           float64               `json:"rank"`
		Recommendation string                `json:"recommendation"`
		ArtistRanks    map[string]ArtistRank `json:"artistRanks"`
	}
	ArtistRank struct {
		Rank    float64  `json:"rank"`
		Related []string `json:"related"`
	}
	GenreResponse struct {
		User         []string `json:"user"`
		Spotify      []string `json:"spotify"`
		LastFm       []string `json:"lastFm"`
		Ticketmaster []string `json:"ticketmaster"`
	}
)

func (e *Event) Artists() []Artist {
	artists := []Artist{}
	if e.MainAct != nil {
		artists = append(artists, *e.MainAct)
	}
	if e.Openers != nil {
		artists = append(artists, e.Openers...)
	}
	return artists
}

func (e *Event) ArtistsMut() []*Artist {
	artists := []*Artist{}
	if e.MainAct != nil {
		artists = append(artists, e.MainAct)
	}
	for i := range e.Openers {
		artists = append(artists, &e.Openers[i])
	}
	return artists
}

func (e *Event) String() string {
	str := ""
	if e.MainAct != nil {
		str += fmt.Sprintf("%+v", *(e.MainAct))
	}
	str += fmt.Sprintf("%+v", e.Openers)
	str += fmt.Sprintf("%+v", e.Venue)
	str += fmt.Sprintf("%+v", e.Purchased)
	str += fmt.Sprintf("%+v", e.ID)
	return str
}

func (g *GenreInfo) Genres() []string {
	if len(g.Spotify) > 0 {
		return g.Spotify
	} else if len(g.LastFm) > 0 {
		return g.LastFm
	} else {
		return g.User
	}
}

func (v *Venue) Populated() bool {
	return v != nil && allNotEmpty(v.Name, v.City, v.State)
}

func (a *Artist) Populated() bool {
	return a != nil && a.Name != ""
}

func (e *Event) Populated() bool {
	artistsPopulated := e.MainAct.Populated()
	for _, opener := range e.Openers {
		artistsPopulated = artistsPopulated || opener.Populated()
	}
	return artistsPopulated && e.Venue.Populated() && e.Date != ""
}

func allNotEmpty(fields ...string) bool {
	for _, f := range fields {
		if len(f) == 0 {
			return false
		}
	}
	return true
}
