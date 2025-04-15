package data

import "fmt"

type (
	Venue struct {
		Name  string `json:"name"`
		City  string `json:"city"`
		State string `json:"state"`
		Id    string `json:"id"`
	}
	Artist struct {
		Name   string `json:"name"`
		Genre  string
		Genres GenreInfo `json:"genres"`
		Id     string    `json:"id"`
		MbId   string    `json:"mbId"`
	}
	GenreInfo struct {
		LfmGenres  []string
		TmGenres   []string
		UserGenres []string
	}
	Event struct {
		MainAct   *Artist  `json:"mainAct"`
		Openers   []Artist `json:"openers"`
		Venue     Venue    `json:"venue"`
		Date      string   `json:"date"`
		Purchased bool     `json:"purchased"`
		Id        string   `json:"id"`
		TmId      string   `json:"tmId"`
	}
	EventDetails struct {
		Name       string    `json:"name"`
		EventGenre string    `json:"genre"`
		Price      string    `json:"price"`
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
)

func (e *Event) String() string {
	str := ""
	if e.MainAct != nil {
		str += fmt.Sprintf("%v", *e.MainAct)
	}
	str += fmt.Sprintf("%v", e.Openers)
	str += fmt.Sprintf("%v", e.Venue)
	str += fmt.Sprintf("%v", e.Purchased)
	str += fmt.Sprintf("%v", e.Id)
	str += fmt.Sprintf("%v", e.TmId)
	return str
}

func (v *Venue) Populated() bool {
	return allNotEmpty(v.Name, v.City, v.State)
}

func (a *Artist) Populated() bool {
	return a.Name != ""
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
