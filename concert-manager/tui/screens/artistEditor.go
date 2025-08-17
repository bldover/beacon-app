package screens

import (
	"concert-manager/domain"
	"concert-manager/search"
	"concert-manager/tui/input"
	"concert-manager/tui/output"
	"strings"
)

type artistCache interface {
	GetArtists() []domain.Artist
}

type Editor struct {
	ArtistCache  artistCache
	ReturnScreen Screen
	actions      []string
	artist       *domain.Artist
	tempArtist   domain.Artist
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

func (e *Editor) SetArtist(artist *domain.Artist) {
	e.artist = artist
	e.tempArtist.Name = e.artist.Name
	e.tempArtist.Genres = e.artist.Genres
}

func (e Editor) Title() string {
	return "Edit Artist"
}

func (e Editor) DisplayData() {
	output.Displayf("%+v\n\n", output.FormatArtistExpanded(e.tempArtist))
}

func (e Editor) Actions() []string {
	return e.actions
}

func (e *Editor) NextScreen(i int) Screen {
	switch i {
	case searchArtist:
		name := input.PromptAndGetInput("artist name to search", input.NoValidation)
		matches := search.SearchArtists(name, e.ArtistCache.GetArtists(), pageSize, search.LenientTolerance)
		selectScreen := &Selector[domain.Artist]{
			ScreenTitle: "Select Artist",
			Next:        e.ReturnScreen,
			Options:     matches,
			HandleSelect: func(v domain.Artist) {
				*e.artist = v
			},
			Formatter: output.FormatArtists,
		}
		return selectScreen
	case setArtistName:
		e.tempArtist.Name = input.PromptAndGetInput("artist name", input.NoValidation)
	case setArtistGenre:
		input := input.PromptAndGetInput("artist genre", input.NoValidation)
		e.tempArtist.Genres.User = strings.Split(input, ",")
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
