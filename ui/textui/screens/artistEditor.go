package screens

import (
	"concert-manager/data"
	"concert-manager/ui/textui/input"
	"concert-manager/ui/textui/output"
	"concert-manager/util"
)

type artistSearch interface {
	FindFuzzyArtistMatchesByName(string) []data.Artist
	WithMaxCount(int)
}

type Editor struct {
	Search       artistSearch
	ReturnScreen Screen
	actions      []string
	artist       *data.Artist
	tempArtist   data.Artist
}

const (
	searchArtist = iota + 1
	setArtistName
	setArtistGenre
	saveArtist
	cancelArtistEdit
)

func NewArtistEditScreen() *Editor {
	e := Editor{}
	e.actions = []string{"Search Artists", "Set Name", "Set Genre", "Save Artist", "Cancel"}
	return &e
}

func (e *Editor) SetArtist(artist *data.Artist) {
	e.artist = artist
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

func (e *Editor) NextScreen(i int) Screen {
	switch i {
	case searchArtist:
		name := input.PromptAndGetInput("artist name to search", input.NoValidation)
		e.Search.WithMaxCount(pageSize)
		matches := e.Search.FindFuzzyArtistMatchesByName(name)
		selectScreen := &Selector[data.Artist]{
			ScreenTitle: "Select Artist",
			Next:        e.ReturnScreen,
			Options:     matches,
			HandleSelect: func(v data.Artist) {
				*e.artist = v
			},
			Formatter: util.FormatArtist,
		}
		return selectScreen
	case setArtistName:
		e.tempArtist.Name = input.PromptAndGetInput("artist name", input.NoValidation)
	case setArtistGenre:
		e.tempArtist.Genre = input.PromptAndGetInput("artist genre", input.NoValidation)
	case saveArtist:
		if e.tempArtist.Populated() {
			*e.artist = e.tempArtist
			return nil
		} else {
			output.Displayln("Failed to save artist: all fields are required")
		}
	case cancelArtistEdit:
		return nil
	}
	return e
}
