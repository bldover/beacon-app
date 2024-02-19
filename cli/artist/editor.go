package artist

import (
	"concert-manager/cli"
	"concert-manager/data"
	"concert-manager/out"
	"context"
	"strings"
)

type Editor struct {
    artist *data.Artist
	tempArtist data.Artist
	Artists *[]data.Artist
	ArtistAdder ArtistAdder
	AddEventScreen cli.Screen
	actions []string
}

const (
	search = iota + 1
	setName
	setGenre
	save
	cancel
)

type ArtistAdder interface {
    AddArtist(context.Context, data.Artist) error
}

func NewEditScreen() *Editor {
	e := Editor{}
	e.actions = []string{"Search Artists", "Set Name", "Set Genre", "Save Artist", "Cancel"}
    return &e
}

func (e *Editor) AddArtistContext(artist *data.Artist) {
    e.artist = artist
	e.tempArtist.Name = strings.Clone(artist.Name)
	e.tempArtist.Genre = strings.Clone(artist.Genre)
}

func (e Editor) Title() string {
    return "Edit Artist"
}

func (e Editor) DisplayData() {
    out.Displayf("%+v\n", e.tempArtist)
}

func (e Editor) Actions() []string {
    return e.actions
}

func (e *Editor) NextScreen(i int) cli.Screen {
	switch i {
    case search:
		e.handleSearch()
	case setName:
		e.tempArtist.Name = cli.PromptAndGetInput("artist name", cli.NoValidation)
	case setGenre:
		e.tempArtist.Genre = cli.PromptAndGetInput("artist genre", cli.NoValidation)
	case save:
		if err := e.ArtistAdder.AddArtist(context.Background(), e.tempArtist); err != nil {
			out.Displayf("Failed to save artist: %v\n", err)
		} else {
			*e.artist = e.tempArtist
			return e.AddEventScreen
		}
	case cancel:
		return e.AddEventScreen
	}
	return e
}

func (e *Editor) handleSearch() {
	out.Displayln("Not yet implemented!")
	//TODO: add search functionality
}
