package db

import (
	"concert-manager/data"
	"concert-manager/log"
	"context"

	"cloud.google.com/go/firestore"
)

const artistCollection string = "artists"

var artistFields = []string{"Name", "MbId", "Genre", "Genres"}

type ArtistClient struct {
	Connection *Firestore
}

type ArtistEntity struct {
	Name   string
	MbId   string
	Genre  *string
	Genres GenreInfoEntity
}

type GenreInfoEntity struct {
	LfmGenres  []string
	TmGenres   []string
	UserGenres []string
}

func (c *ArtistClient) Add(ctx context.Context, artist data.Artist) (string, error) {
	log.Debug("Attempting to add artist", artist)
	existingArtist, err := c.findDocRef(ctx, artist.Id)
	if err != nil {
		log.Errorf("Error occurred while checking if artist %v already exists, %v", artist, err)
		return "", err
	}
	if existingArtist.Exists() {
		log.Debugf("Skipping adding artist because it already exists %+v, %v", artist, existingArtist.Ref.ID)
		return existingArtist.Ref.ID, nil
	}

	genreEntity := GenreInfoEntity{artist.Genres.LfmGenres, artist.Genres.TmGenres, artist.Genres.UserGenres}
	artistEntity := ArtistEntity{Name: artist.Name, MbId: artist.MbId, Genres: genreEntity}

	artists := c.Connection.Client.Collection(artistCollection)
	docRef, _, err := artists.Add(ctx, artistEntity)
	if err != nil {
		log.Errorf("Failed to add new artist %+v, %v", artist, err)
		return "", err
	}
	log.Infof("Created new artist %+v", docRef.ID)
	return docRef.ID, nil
}

func (c *ArtistClient) Update(ctx context.Context, artist data.Artist) error {
	log.Debug("Attempting to update artist", artist)
	artistDoc, err := c.Connection.Client.Collection(artistCollection).Doc(artist.Id).Get(ctx)
	if err != nil {
		log.Errorf("Error while updating artist %+v, %v", artist, err)
		return err
	}
	if !artistDoc.Exists() {
		log.Errorf("Artist does not exist in update for %+v, %v", artist, err)
		return err
	}

	genreEntity := GenreInfoEntity{artist.Genres.LfmGenres, artist.Genres.TmGenres, artist.Genres.UserGenres}
	artistEntity := ArtistEntity{Name: artist.Name, MbId: artist.MbId, Genres: genreEntity}

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

func (c *ArtistClient) FindAll(ctx context.Context) ([]data.Artist, error) {
	log.Debug("Finding all artists")
	artistDocs, err := c.Connection.Client.Collection(artistCollection).
		Select(artistFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		log.Error("Error while finding all artists,", err)
		return nil, err
	}

	artists := []data.Artist{}
	for _, a := range artistDocs {
		artists = append(artists, toArtist(a))
	}
	log.Debugf("Found %d artists", len(artists))
	return artists, nil
}

func toArtist(doc *firestore.DocumentSnapshot) data.Artist {
	artist := data.Artist{
		Name:  doc.Data()["Name"].(string),
		Id:    doc.Ref.ID,
		Genres: data.GenreInfo{
			LfmGenres:  []string{},
			UserGenres: []string{},
		},
	}

	if mbid, ok := doc.Data()["MbId"].(string); ok {
		artist.MbId = mbid
	}

	if genre, ok := doc.Data()["Genre"].(string); ok {
		artist.Genre = genre
	}

	if genres, ok := doc.Data()["Genres"].(map[string]interface{}); ok {
		if lfmGenres, ok := genres["LfmGenres"].([]interface{}); ok {
			for _, genre := range lfmGenres {
				if genreStr, ok := genre.(string); ok {
					artist.Genres.LfmGenres = append(artist.Genres.LfmGenres, genreStr)
				}
			}
		}

		if userGenres, ok := genres["UserGenres"].([]interface{}); ok {
			for _, genre := range userGenres {
				if genreStr, ok := genre.(string); ok {
					artist.Genres.UserGenres = append(artist.Genres.UserGenres, genreStr)
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

func (c *ArtistClient) findAllDocs(ctx context.Context) (*map[string]data.Artist, error) {
	artistDocs, err := c.Connection.Client.Collection(artistCollection).
		Select(artistFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	artists := make(map[string]data.Artist)
	for _, a := range artistDocs {
		artists[a.Ref.ID] = toArtist(a)
	}

	return &artists, nil
}
