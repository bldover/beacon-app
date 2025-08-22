package firestore

import (
	"concert-manager/domain"
	"concert-manager/log"
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const artistCollection string = "artists"

var artistFields = []string{"Name", "Genre", "Genres", "ID"}

type (
	ArtistClient struct {
		Connection *Firestore
	}

	ArtistEntity struct {
		Name   string
		Genre  *string
		Genres GenreInfoEntity
		ID     ArtistIDEntity
	}

	GenreInfoEntity struct {
		Spotify      []string
		LastFm       []string
		Ticketmaster []string
		User         []string
	}

	ArtistIDEntity struct {
		Primary      string
		Spotify      string
		Ticketmaster string
		MusicBrainz  string
	}
)

func (c *ArtistClient) Add(ctx context.Context, artist domain.Artist) (string, error) {
	log.Debug("Attempting to add artist", artist)
	existingArtist, err := c.findDocRef(ctx, artist.ID.Primary)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("Error occurred while checking if artist %v already exists, %v", artist, err)
		return "", err
	}
	if existingArtist.Exists() {
		log.Debugf("Skipping adding artist because it already exists %+v, %v", artist, existingArtist.Ref.ID)
		return existingArtist.Ref.ID, nil
	}

	genreEntity := GenreInfoEntity{
		Spotify:      artist.Genres.Spotify,
		LastFm:       artist.Genres.LastFm,
		Ticketmaster: artist.Genres.Ticketmaster,
		User:         artist.Genres.User,
	}

	idEntity := ArtistIDEntity{
		Primary:      artist.ID.Primary,
		Spotify:      artist.ID.Spotify,
		Ticketmaster: artist.ID.Ticketmaster,
		MusicBrainz:  artist.ID.MusicBrainz,
	}
	artistEntity := ArtistEntity{Name: artist.Name, Genre: &artist.Genre, Genres: genreEntity, ID: idEntity}

	artists := c.Connection.Client.Collection(artistCollection)
	docRef, _, err := artists.Add(ctx, artistEntity)
	if err != nil {
		log.Errorf("Failed to add new artist %+v, %v", artist, err)
		return "", err
	}
	log.Infof("Created new artist %+v", docRef.ID)
	return docRef.ID, nil
}

func (c *ArtistClient) Update(ctx context.Context, artist domain.Artist) error {
	log.Debug("Attempting to update artist", artist)
	artistDoc, err := c.Connection.Client.Collection(artistCollection).Doc(artist.ID.Primary).Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("Error while updating artist %+v, %v", artist, err)
		return err
	}
	if !artistDoc.Exists() {
		log.Errorf("Artist does not exist in update for %+v, %v", artist, err)
		return err
	}

	genreEntity := GenreInfoEntity{
		Spotify:      artist.Genres.Spotify,
		LastFm:       artist.Genres.LastFm,
		Ticketmaster: artist.Genres.Ticketmaster,
		User:         artist.Genres.User,
	}
	idEntity := ArtistIDEntity{
		Primary:      artist.ID.Primary,
		Spotify:      artist.ID.Spotify,
		Ticketmaster: artist.ID.Ticketmaster,
		MusicBrainz:  artist.ID.MusicBrainz,
	}
	artistEntity := ArtistEntity{Name: artist.Name, Genre: &artist.Genre, Genres: genreEntity, ID: idEntity}

	_, err = artistDoc.Ref.Set(ctx, artistEntity)
	if err != nil {
		log.Errorf("Failed to update artist %+v, %v", artist, err)
		return err
	}
	log.Info("Successfully updated artist", artist)
	return nil
}

func (c *ArtistClient) Delete(ctx context.Context, id string) error {
	log.Debug("Attempting to delete artist", id)
	artistDoc, err := c.Connection.Client.Collection(artistCollection).Doc(id).Get(ctx)
	if err != nil {
		log.Error("Error while deleting artist", id, err)
		return err
	}
	_, err = artistDoc.Ref.Delete(ctx)
	if err != nil {
		log.Error("Failed to delete artist", id, err)
		return err
	}
	log.Info("Successfully deleted artist", id)
	return nil
}

func (c *ArtistClient) FindAll(ctx context.Context) ([]domain.Artist, error) {
	log.Debug("Finding all artists")
	artistDocs, err := c.Connection.Client.Collection(artistCollection).
		Select(artistFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		log.Error("Error while finding all artists,", err)
		return nil, err
	}

	artists := []domain.Artist{}
	for _, a := range artistDocs {
		artists = append(artists, toArtist(a))
	}
	log.Debugf("Found %d artists", len(artists))
	return artists, nil
}

func toArtist(doc *firestore.DocumentSnapshot) domain.Artist {
	artist := domain.Artist{
		Name: doc.Data()["Name"].(string),
		ID: domain.ID{
			Primary: doc.Ref.ID,
		},
		Genres: domain.GenreInfo{
			Spotify:      []string{},
			LastFm:       []string{},
			Ticketmaster: []string{},
			User:         []string{},
		},
	}

	if genre, ok := doc.Data()["Genre"].(string); ok {
		artist.Genre = genre
	}

	if ids, ok := doc.Data()["ID"].(map[string]any); ok {
		if spotifyId, ok := ids["Spotify"].(string); ok {
			artist.ID.Spotify = spotifyId
		}
		if ticketmasterId, ok := ids["Ticketmaster"].(string); ok {
			artist.ID.Ticketmaster = ticketmasterId
		}
		if mbid, ok := ids["MBID"].(string); ok {
			artist.ID.MusicBrainz = mbid
		}
	}

	if genres, ok := doc.Data()["Genres"].(map[string]any); ok {
		if spotifyGenres, ok := genres["Spotify"].([]any); ok {
			for _, genre := range spotifyGenres {
				if genreStr, ok := genre.(string); ok {
					artist.Genres.Spotify = append(artist.Genres.Spotify, genreStr)
				}
			}
		}

		if lastFmGenres, ok := genres["LastFm"].([]any); ok {
			for _, genre := range lastFmGenres {
				if genreStr, ok := genre.(string); ok {
					artist.Genres.LastFm = append(artist.Genres.LastFm, genreStr)
				}
			}
		}

		if userGenres, ok := genres["User"].([]any); ok {
			for _, genre := range userGenres {
				if genreStr, ok := genre.(string); ok {
					artist.Genres.User = append(artist.Genres.User, genreStr)
				}
			}
		}
	}

	return artist
}

func (c *ArtistClient) findDocRef(ctx context.Context, id string) (*firestore.DocumentSnapshot, error) {
	if id == "" {
		return &firestore.DocumentSnapshot{}, nil
	}
	return c.Connection.Client.Collection(artistCollection).Doc(id).Get(ctx)
}

func (c *ArtistClient) findAllDocs(ctx context.Context) (*map[string]domain.Artist, error) {
	artistDocs, err := c.Connection.Client.Collection(artistCollection).
		Select(artistFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	artists := make(map[string]domain.Artist)
	for _, a := range artistDocs {
		artists[a.Ref.ID] = toArtist(a)
	}

	return &artists, nil
}
