package artist

import (
	"concert-manager/data"
	"concert-manager/ui/terminal/input"
	"concert-manager/ui/terminal/output"
	"concert-manager/ui/terminal/screens"
)

type artistAddCache interface {
	AddArtist(data.Artist) error
}

type Editor struct {
	actions      []string
	artist       *data.Artist
	tempArtist   data.Artist
	returnScreen screens.Screen
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

func (e *Editor) AddContext(returnScreen screens.Screen, props ...any) {
	e.returnScreen = returnScreen
	e.artist = props[0].(*data.Artist)
	e.tempArtist.Name = e.artist.Name
	e.tempArtist.Genre = e.artist.Genre
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
		if e.tempArtist.Populated() {
			*e.artist = e.tempArtist
			return e.returnScreen
		} else {
			output.Displayln("Failed to save artist: all fields are required")
		}
	case cancel:
		return e.returnScreen
	}
	return e
}

func (e *Editor) handleSearch() {
	output.Displayln("Not yet implemented!")
	//TODO: add search functionality
}
