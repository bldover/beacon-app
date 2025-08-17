package spotify

import "concert-manager/external"

type track struct {
	Id      string   `json:"id"`
	Title   string   `json:"name"`
	Artists []artist `json:"artists"`
}

type artist struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

type artistInfoResponse struct {
	Artists []artist `json:"artists"`
}

type artistSearchResponse struct {
	Artists struct {
		Items []artist `json:"items"`
	} `json:"artists"`
}

type savedTrackResponse struct {
	Next        string `json:"next"`
	Total       int    `json:"total"`
	SavedTracks []struct {
		Track track `json:"track"`
	} `json:"items"`
}

type topTrackResponse struct {
	Next      string  `json:"next"`
	Total     int     `json:"total"`
	Offset    int     `json:"offset"`
	TopTracks []track `json:"items"`
}

type topArtistResponse struct {
	Next       string   `json:"next"`
	Total      int      `json:"total"`
	Offset     int      `json:"offset"`
	TopArtists []artist `json:"items"`
}

func mapSpotifyTracks(rawTracks []track) []external.Track {
	tracks := []external.Track{}
	for _, rawTrack := range rawTracks {
		tracks = append(tracks, mapSpotifyTrack(rawTrack))
	}
	return tracks
}

func mapSpotifyTrack(rawTrack track) external.Track {
	return external.Track{
		Title:   rawTrack.Title,
		Artists: mapSpotifyArtists(rawTrack.Artists),
	}
}

func mapSpotifyArtists(rawArtists []artist) []external.Artist {
	artists := []external.Artist{}
	for _, rawArtist := range rawArtists {
		artists = append(artists, mapSpotifyArtist(rawArtist))
	}
	return artists
}

func mapSpotifyArtist(rawArtist artist) external.Artist {
	return external.Artist{
		Id:   rawArtist.Id,
		Name: rawArtist.Name,
	}
}
