package data

type (
	Venue struct {
		Name  string `json:"name"`
		City  string `json:"city"`
		State string `json:"state"`
		Id    string `json:"id"`
	}
	Artist struct {
		Name  string `json:"name"`
		Genre string `json:"genre"`
		Id    string `json:"id"`
	}
	Event struct {
		MainAct   Artist   `json:"mainAct"`
		Openers   []Artist `json:"openers"`
		Venue     Venue    `json:"venue"`
		Date      string   `json:"date"`
		Purchased bool     `json:"purchased"`
		Id        string   `json:"id"`
		TmId      string   `json:"tmId"`
	}
	EventDetails struct {
		Name       string `json:"name"`
		EventGenre string `json:"genre"`
		Price      string `json:"price"`
		Event      Event  `json:"event"`
	}
	EventRank struct {
		Event       EventDetails `json:"event"`
		ArtistRanks []ArtistRank `json:"artistRanks"`
		Rank        float64      `json:"rank"`
	}
	ArtistRank struct {
		Artist  Artist   `json:"artist"`
		Rank    float64  `json:"rank"`
		Related []string `json:"related"`
	}
)

func (v *Venue) Populated() bool {
	return allNotEmpty(v.Name, v.City, v.State)
}

func (a *Artist) Populated() bool {
	return allNotEmpty(a.Name, a.Genre)
}

func (a *Artist) Invalid() bool {
	return allNotEmpty(a.Name) != allNotEmpty(a.Genre)
}

func (e *Event) Populated() bool {
	invalidArtist := e.MainAct.Invalid()
	populated := e.MainAct.Populated()
	for _, opener := range e.Openers {
		invalidArtist = invalidArtist || opener.Invalid()
		populated = populated || opener.Populated()
	}
	return populated && !invalidArtist && e.Venue.Populated() && e.Date != ""
}

func (e Event) Equals(o Event) bool {
	return e.MainAct.Equals(o.MainAct) && e.Venue.Equals(o.Venue) && e.Date == o.Date
}

func (a Artist) Equals(o Artist) bool {
	return a.Name == o.Name && a.Genre == o.Genre
}

func (v Venue) Equals(o Venue) bool {
	return v.Name == o.Name && v.City == o.City && v.State == o.State
}

func allNotEmpty(fields ...string) bool {
	for _, f := range fields {
		if len(f) == 0 {
			return false
		}
	}
	return true
}
