package data

import (
	"strconv"
	"strings"
)

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
	return populated && !invalidArtist && e.Venue.Populated() && validDate(e.Date)
}

func allNotEmpty(fields ...string) bool {
	for _, f := range fields {
		if len(f) == 0 {
			return false
		}
	}
	return true
}

// expects format "#/#/# but doesn't check for non-standard date values (like day 32)
// since time.Date handles these with overflow. Consumers of the Event type should regulate
// their own date values if possible overflow is undesired
func validDate(date string) bool {
	parts := strings.Split(date, "/")
	if len(parts) != 3 {
		return false
	}
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return false
	}
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return false
	}
	if _, err := strconv.Atoi(parts[2]); err != nil {
		return false
	}
	return true
}
