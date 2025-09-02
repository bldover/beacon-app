package finder

import (
	"concert-manager/domain"
	"concert-manager/external"
	"concert-manager/log"
	"strings"
)

type MetadataFinder struct {
	Spotify metadataProvider
	LastFm  metadataProvider
}

type artistCache interface {
	GetArtists() []domain.Artist
}

type metadataProvider interface {
	ArtistInfoById(id string) (external.ArtistInfo, error)
	SearchByName(name string) (external.ArtistInfo, error)
}

type ArtistInfo struct {
	Name      string
	Genres    []string
	SpotifyID string
}

type eventPosition struct {
	eventIndex  int
	isMainAct   bool
	openerIndex int
}

func (f MetadataFinder) PopulateMetadata(events []domain.EventDetails) []domain.EventDetails {
	log.Infof("Populating metadata for %v events", len(events))
	result := make([]domain.EventDetails, len(events))
	for i, event := range events {
		result[i] = domain.CloneEventDetail(event)
	}

	artistToEvents := f.buildArtistsToUpdateMap(result)
	for artist, eventPositions := range artistToEvents {
		if len(artist.Genres.Spotify) == 0 {
			f.updateSpotifyMetadata(artist, result, eventPositions)
		}
		if len(artist.Genres.LastFm) == 0 {
			f.updateLastFmMetadata(artist, result, eventPositions)
		}
	}

	log.Info("Metadata loaded")
	return result
}

func (f MetadataFinder) buildArtistsToUpdateMap(events []domain.EventDetails) map[*domain.Artist][]eventPosition {
	artistToEvents := make(map[*domain.Artist][]eventPosition)
	for i, event := range events {
		for j, artist := range event.Event.ArtistsMut() {
			loc := eventPosition{eventIndex: i, isMainAct: j == 0}
			artistToEvents[artist] = append(artistToEvents[artist], loc)
		}
	}
	return artistToEvents
}

func (f MetadataFinder) updateSpotifyMetadata(artist *domain.Artist, events []domain.EventDetails, locations []eventPosition) {
	var artistInfo external.ArtistInfo
	var err error
	if artist.ID.Spotify != "" {
		artistInfo, err = f.Spotify.ArtistInfoById(artist.ID.Spotify)
		if err != nil {
			if _, ok := err.(external.NotFoundError); ok {
				log.Errorf("Unable to find Spotify artist by ID: %s", err.Error())
				artist.ID.Spotify = ""
			} else {
				log.Errorf("Failed to fetch artist genre from Spotify: %v", err)
				return
			}
		} else {
			artist.ID.Spotify = artistInfo.Id
			artist.Genres.Spotify = toLower(artistInfo.Genres)
			return
		}
	}

	if artist.Name != "" {
		artistInfo, err = f.Spotify.SearchByName(artist.Name)
		if err != nil {
			log.Errorf("Failed to fetch artist genres by name from Spotify: %v", err)
			return
		}
		artist.ID.Spotify = artistInfo.Id
		artist.Genres.Spotify = toLower(artistInfo.Genres)
	}
}

func (f MetadataFinder) updateLastFmMetadata(artist *domain.Artist, events []domain.EventDetails, locations []eventPosition) {
	var artistInfo external.ArtistInfo
	var err error

	if artist.ID.MusicBrainz != "" {
		artistInfo, err = f.LastFm.ArtistInfoById(artist.ID.MusicBrainz)
		if err != nil {
			if _, ok := err.(external.NotFoundError); ok {
				log.Errorf("Unable to find LastFm artist by ID: %s", err.Error())
				artist.ID.MusicBrainz = ""
			} else {
				log.Errorf("Failed to fetch artist genre from LastFm: %v", err)
				return
			}
		} else {
			artist.ID.MusicBrainz = artistInfo.Id
			artist.Genres.LastFm = toLower(artistInfo.Genres)
			return
		}
	}

	if artist.Name != "" {
		artistInfo, err = f.LastFm.SearchByName(artist.Name)
		if err != nil {
			log.Errorf("Failed to fetch artist genres by name from LastFm: %v", err)
			return
		}
		artist.ID.MusicBrainz = artistInfo.Id
		artist.Genres.LastFm = toLower(artistInfo.Genres)
	}
}

func (f MetadataFinder) ReloadMetadata(artists []domain.Artist) ([]domain.Artist, error) {
	result := domain.CloneArtists(artists)

	spotifyIdsToFetch := map[string]int{}
	spotifyArtistsToFetchByName := map[string][]int{}
	lastfmIdsToFetch := map[string]int{}
	lastfmArtistsToFetchByName := map[string][]int{}

	for i, artist := range result {
		if artist.ID.Spotify != "" {
			spotifyIdsToFetch[artist.ID.Spotify] = i
		} else {
			spotifyArtistsToFetchByName[artist.Name] = append(spotifyArtistsToFetchByName[artist.Name], i)
		}

		if artist.ID.MusicBrainz != "" {
			lastfmIdsToFetch[artist.ID.MusicBrainz] = i
		} else {
			lastfmArtistsToFetchByName[artist.Name] = append(lastfmArtistsToFetchByName[artist.Name], i)
		}
	}

	for name, indices := range spotifyArtistsToFetchByName {
		artistInfo, err := f.Spotify.SearchByName(name)
		if err != nil {
			log.Errorf("Failed to fetch Spotify metadata for artist %s: %v", name, err)
			continue
		}

		for _, idx := range indices {
			result[idx].ID.Spotify = artistInfo.Id
			result[idx].Genres.Spotify = toLower(artistInfo.Genres)
		}
	}

	for id, idx := range spotifyIdsToFetch {
		artistInfo, err := f.Spotify.ArtistInfoById(id)
		if err != nil {
			log.Errorf("Failed to fetch Spotify metadata for ID %s: %v", id, err)
			continue
		}
		result[idx].Genres.Spotify = toLower(artistInfo.Genres)
	}

	for name, indices := range lastfmArtistsToFetchByName {
		artistInfo, err := f.LastFm.SearchByName(name)
		if err != nil {
			log.Errorf("Failed to fetch LastFm metadata for artist %s: %v", name, err)
			continue
		}

		for _, idx := range indices {
			result[idx].ID.MusicBrainz = artistInfo.Id
			result[idx].Genres.LastFm = toLower(artistInfo.Genres)
		}
	}

	for id, idx := range lastfmIdsToFetch {
		artistInfo, err := f.LastFm.ArtistInfoById(id)
		if err != nil {
			log.Errorf("Failed to fetch LastFm metadata for ID %s: %v", id, err)
			continue
		}
		result[idx].Genres.LastFm = toLower(artistInfo.Genres)
	}

	return result, nil
}

func toLower(genres []string) []string {
	lower := []string{}
	for _, genre := range genres {
		lower = append(lower, strings.ToLower(genre))
	}
	return lower
}
