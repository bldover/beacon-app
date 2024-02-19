package data

type (
	Venue struct {
		Name    string
		City    string
		State   string
	}
	Artist struct {
		Name  string
		Genre string
	}
	Event struct {
		MainAct Artist
		Openers []Artist
		Venue   Venue
		Date    string
		Purchased bool
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
	return populated && !invalidArtist && e.Venue.Populated() && ValidDate(e.Date)
}

func allNotEmpty(fields ...string) bool {
	for _, f := range fields {
		if len(f) == 0 {
			return false
		}
	}
	return true
}
