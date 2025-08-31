package loader

import (
	"concert-manager/domain"
	"concert-manager/log"
	"context"
	"errors"
	"slices"
	"strings"
)

type artistCache interface {
	UpdateArtist(string, domain.Artist) error
	GetArtists() []domain.Artist
}

type metadataProvider interface {
	ReloadMetadata([]domain.Artist) ([]domain.Artist, error)
}

type GenreLoader struct {
	Cache            artistCache
	MetadataProvider metadataProvider
}

func (l *GenreLoader) ReloadGenres(ctx context.Context, artists []string) (int, error) {
	log.Info("Reloading genres for artists:", artists)
	var targetArtists []domain.Artist
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
	log.Debug("Reload genres target artists:", targetArtists)

	updatedCount := 0

	updatedArtists, err := l.MetadataProvider.ReloadMetadata(targetArtists)
	if err != nil {
		return 0, err
	}

	for _, artist := range updatedArtists {
		if len(artist.Genres.Spotify) == 0 {
			log.Info("No Spotify genres present for artist", artist.Name)
		}

		// remove once Genre is fully migrated to UserGenres
		isUserGenreValid := len(artist.Genres.User) == 0 && artist.Genre != "" && artist.Genre != "Unspecified"
		if isUserGenreValid {
			artist.Genres.User = []string{strings.ToLower(artist.Genre)}
		}

		if err := l.Cache.UpdateArtist(artist.ID.Primary, artist); err != nil {
			log.Error("Failed to update genres for artist", artist, err)
		}

		updatedCount++
	}

	if updatedCount != len(targetArtists) {
		return updatedCount, errors.New("failed to load genres for some artists")
	}

	log.Debugf("Finished reloading genres for %v/%v artists", updatedCount, len(targetArtists))
	return updatedCount, nil
}
