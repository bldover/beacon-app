package firestore

import (
	"concert-manager/domain"
	"concert-manager/log"
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const albumCollection = "albums"

var albumFields = []string{"Name", "Artist", "Year", "Signed", "ID"}

type (
	AlbumClient struct {
		Connection *Firestore
	}

	AlbumEntity struct {
		Name   string
		Artist ArtistEntity
		Year   int
		Signed bool
		ID     string
	}
)

func (c *AlbumClient) Add(ctx context.Context, album domain.Album) (string, error) {
	log.Debug("Attempting to add album", album)
	existingAlbum, err := c.findDocRef(ctx, album.ID)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("Error occurred while checking if album %v already exists, %v", album, err)
		return "", err
	}
	if existingAlbum.Exists() {
		log.Debugf("Skipping adding album because it already exists %+v, %v", album, existingAlbum.Ref.ID)
		return existingAlbum.Ref.ID, nil
	}

	artistEntity := toArtistEntity(album.Artist)
	albumEntity := AlbumEntity{
		Name:   album.Name,
		Artist: artistEntity,
		Year:   album.Year,
		Signed: album.Signed,
	}

	albums := c.Connection.Client.Collection(albumCollection)
	docRef, _, err := albums.Add(ctx, albumEntity)
	if err != nil {
		log.Errorf("Failed to add new album %+v, %v", album, err)
		return "", err
	}
	log.Infof("Created new album %+v", docRef.ID)
	return docRef.ID, nil
}

func (c *AlbumClient) Update(ctx context.Context, album domain.Album) error {
	log.Debug("Attempting to update album", album)
	albumDoc, err := c.Connection.Client.Collection(albumCollection).Doc(album.ID).Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("Error while updating album %+v, %v", album, err)
		return err
	}
	if !albumDoc.Exists() {
		log.Errorf("Album does not exist in update for %+v, %v", album, err)
		return err
	}

	artistEntity := toArtistEntity(album.Artist)
	albumEntity := AlbumEntity{
		Name:   album.Name,
		Artist: artistEntity,
		Year:   album.Year,
		Signed: album.Signed,
		ID:     album.ID,
	}

	_, err = albumDoc.Ref.Set(ctx, albumEntity)
	if err != nil {
		log.Errorf("Failed to update album %+v, %v", album, err)
		return err
	}
	log.Info("Successfully updated album", album)
	return nil
}

func (c *AlbumClient) Delete(ctx context.Context, id string) error {
	log.Debug("Attempting to delete album", id)
	albumDoc, err := c.Connection.Client.Collection(albumCollection).Doc(id).Get(ctx)
	if err != nil {
		log.Error("Error while deleting album", id, err)
		return err
	}
	_, err = albumDoc.Ref.Delete(ctx)
	if err != nil {
		log.Error("Failed to delete album", id, err)
		return err
	}
	log.Info("Successfully deleted album", id)
	return nil
}

func (c *AlbumClient) FindAll(ctx context.Context) ([]domain.Album, error) {
	log.Debug("Finding all albums")
	albumDocs, err := c.Connection.Client.Collection(albumCollection).
		Select(albumFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		log.Error("Error while finding all albums,", err)
		return nil, err
	}

	albums := []domain.Album{}
	for _, r := range albumDocs {
		albums = append(albums, toAlbum(r))
	}
	log.Debugf("Found %d albums", len(albums))
	return albums, nil
}

func toAlbum(doc *firestore.DocumentSnapshot) domain.Album {
	data := doc.Data()
	album := domain.Album{
		ID: doc.Ref.ID,
	}

	if name, ok := data["Name"].(string); ok {
		album.Name = name
	}
	if year, ok := data["Year"].(int64); ok {
		album.Year = int(year)
	}
	if signed, ok := data["Signed"].(bool); ok {
		album.Signed = signed
	}
	if artistData, ok := data["Artist"].(map[string]any); ok {
		album.Artist = toArtistFromMap(artistData, "")
	}

	return album
}

func toArtistEntity(artist domain.Artist) ArtistEntity {
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
	return ArtistEntity{Name: artist.Name, Genre: &artist.Genre, Genres: genreEntity, ID: idEntity}
}

func toArtistFromMap(data map[string]any, id string) domain.Artist {
	artist := domain.Artist{
		ID: domain.ID{Primary: id},
		Genres: domain.GenreInfo{
			Spotify:      []string{},
			LastFm:       []string{},
			Ticketmaster: []string{},
			User:         []string{},
		},
	}
	if name, ok := data["Name"].(string); ok {
		artist.Name = name
	}
	if genre, ok := data["Genre"].(string); ok {
		artist.Genre = genre
	}
	if ids, ok := data["ID"].(map[string]any); ok {
		if primaryId, ok := ids["Primary"].(string); ok {
			artist.ID.Primary = primaryId
		}
		if spotifyId, ok := ids["Spotify"].(string); ok {
			artist.ID.Spotify = spotifyId
		}
		if ticketmasterId, ok := ids["Ticketmaster"].(string); ok {
			artist.ID.Ticketmaster = ticketmasterId
		}
		if mbid, ok := ids["MusicBrainz"].(string); ok {
			artist.ID.MusicBrainz = mbid
		}
	}
	if genres, ok := data["Genres"].(map[string]any); ok {
		if spotifyGenres, ok := genres["Spotify"].([]any); ok {
			for _, g := range spotifyGenres {
				if gs, ok := g.(string); ok {
					artist.Genres.Spotify = append(artist.Genres.Spotify, gs)
				}
			}
		}
		if lastFmGenres, ok := genres["LastFm"].([]any); ok {
			for _, g := range lastFmGenres {
				if gs, ok := g.(string); ok {
					artist.Genres.LastFm = append(artist.Genres.LastFm, gs)
				}
			}
		}
		if userGenres, ok := genres["User"].([]any); ok {
			for _, g := range userGenres {
				if gs, ok := g.(string); ok {
					artist.Genres.User = append(artist.Genres.User, gs)
				}
			}
		}
	}
	return artist
}

func (c *AlbumClient) findDocRef(ctx context.Context, id string) (*firestore.DocumentSnapshot, error) {
	if id == "" {
		return &firestore.DocumentSnapshot{}, nil
	}
	return c.Connection.Client.Collection(albumCollection).Doc(id).Get(ctx)
}
