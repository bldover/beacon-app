package artist

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/ui/terminal/input"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
	"strings"
)

type Editor struct {
	artist         *data.Artist
	tempArtist     data.Artist
	Artists        *[]data.Artist
	AddEventScreen screens.Screen
	actions        []string
}

const (
	search = iota + 1
	setName
	setGenre
	save
	cancel
)

func NewEditScreen() *Editor {
	e := Editor{}
	e.actions = []string{"Search Artists", "Set Name", "Set Genre", "Save Artist", "Cancel"}
	return &e
}

func (e *Editor) AddArtistContext(artist *data.Artist) {
	log.Debugf("in Artist editor - add artist context: %p", artist)
	log.Debugf("in Artist editor - add artist context: %+v", artist)
	e.artist = artist
	e.tempArtist.Name = strings.Clone(artist.Name)
	e.tempArtist.Genre = strings.Clone(artist.Genre)
	log.Debugf("st Artist editor - add artist context: %p", e.artist)
	log.Debugf("st Artist editor - add artist context: %+v", e.artist)
	log.Debug("tp Artist editor - add artist context:", e.tempArtist)
}

func (e Editor) Title() string {
	return "Edit Artist"
}

func (e Editor) DisplayData() {
	output.Displayf("%+v\n", e.tempArtist)
}

func (e Editor) Actions() []string {
	return e.actions
}

func (e *Editor) NextScreen(i int) screens.Screen {
	switch i {
	case search:
		e.handleSearch()
	case setName:
		e.tempArtist.Name = input.PromptAndGetInput("artist name", input.NoValidation)
	case setGenre:
		e.tempArtist.Genre = input.PromptAndGetInput("artist genre", input.NoValidation)
	case save:
		if e.artist.Populated() {
			*e.artist = e.tempArtist
			return e.AddEventScreen
		} else {
			output.Displayf("Failed to save artist: all fields are required")
		}
	case cancel:
		return e.AddEventScreen
	}
	return e
}

func (e *Editor) handleSearch() {
	output.Displayln("Not yet implemented!")
	//TODO: add search functionality
}
