package api

type (
	Track struct {
		TmId string
		Title string
		Artists []Artist
		Offset int
		Rank float64
	}
	Artist struct {
		TmId string
		Name string
		Offset int
		Rank float64
	}
)
