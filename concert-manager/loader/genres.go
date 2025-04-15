package loader

import (
	"concert-manager/client/lastfm"
	"concert-manager/data"
	"concert-manager/log"
	"context"
	"errors"
	"slices"
)

type artistCache interface {
	UpdateArtist(string, data.Artist) error
	GetArtists() []data.Artist
}

type artistInfoProvider interface {
    GetArtistInfoByName(string) (lastfm.ArtistInfo, error)
}

type GenreLoader struct {
    Cache artistCache
	InfoProvider artistInfoProvider
}

func (l *GenreLoader) PopulateLastFmData(ctx context.Context, artists []string) (int, error) {
	var targetArtists []data.Artist
	allArtists := l.Cache.GetArtists()

	if len(artists) > 0 {
		for _, artist := range allArtists {
			if slices.Contains(artists, artist.Name) {
				targetArtists = append(targetArtists, artist)
			}
		}
	} else {
		targetArtists = allArtists
	}

	updatedCount := 0

	for _, artist := range targetArtists {
		artistInfo, err := l.InfoProvider.GetArtistInfoByName(artist.Name)
		if err != nil {
			log.Error("Failed to retrieve genres for artist", artist.Name, err)
		}

		if len(artistInfo.Tags.Tag) == 0 {
			log.Info("No tags present for artist", artist.Name)
			continue
		}

		var genreNames []string
		for _, tag := range artistInfo.Tags.Tag {
			genreNames = append(genreNames, tag.Name)
		}

		artist.Genres.LfmGenres = genreNames
		artist.Genres.UserGenres = []string{}
		if err := l.Cache.UpdateArtist(artist.Id, artist); err != nil {
			log.Error("Failed to update genres for artist", artist, err)
		}

		updatedCount++
	}

	if updatedCount != len(targetArtists) {
		return updatedCount, errors.New("failed to load genres for some artists")
	}
	return updatedCount, nil
}

// one time load to transform genre field type
func (l *GenreLoader) PopulateUserGenres(ctx context.Context) (int, error) {
	allArtists := l.Cache.GetArtists()
	var updated int
	for _, artist := range allArtists {
		if artist.Genre != "Unspecified" {
			artist.Genres.UserGenres = append(artist.Genres.UserGenres, artist.Genre)
		}
 		err := l.Cache.UpdateArtist(artist.Id, artist)
		if err != nil {
			log.Error("Failed to update genres for artist", artist, err)
		}
		updated++
	}
	if updated != len(allArtists) {
		return updated, errors.New("failed to update some artist genres")
	}
	return updated, nil
}
