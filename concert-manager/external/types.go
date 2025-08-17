package external

type ArtistInfo struct {
	Name   string
	Genres []string
	Id     string
}

type Artist struct {
	Id   string
	Name string
}

type RankedArtist struct {
	Id   string
	Name string
	Rank float64
}

type Track struct {
	Title   string
	Artists []Artist
}

type TimeRange string

const LongTerm = "long_term"
const MediumTerm = "medium_term"
const ShortTerm = "short_term"

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}
