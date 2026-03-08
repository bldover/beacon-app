package firestore

import (
	"concert-manager/domain"
	"concert-manager/log"
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const albumCollection = "albums"

var albumFields = []string{"Name", "ArtistRef", "Year", "Signed", "ID"}

type (
	AlbumClient struct {
		Connection   *Firestore
		ArtistClient *ArtistClient
	}

	AlbumEntity struct {
		Name      string
		ArtistRef *firestore.DocumentRef
		Year      int
		Signed    bool
		ID        string
	}
)

func (c *AlbumClient) Add(ctx context.Context, album domain.Album) (string, error) {
	log.Debug("Attempting to add album", album)
	artistDoc, err := c.ArtistClient.findDocRef(ctx, album.Artist.ID.Primary)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("Error finding existing artist %v while adding album %v", album.Artist.Name, album)
		return "", err
	}
	if !artistDoc.Exists() {
		log.Errorf("No existing artist %v while adding album %v", album.Artist.Name, album)
		return "", errors.New("artist does not exist")
	}

	existingAlbum, err := c.findDocRef(ctx, album.ID)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("Error occurred while checking if album %v already exists, %v", album, err)
		return "", err
	}
	if existingAlbum.Exists() {
		log.Debugf("Skipping adding album because it already exists %+v, %v", album, existingAlbum.Ref.ID)
		return existingAlbum.Ref.ID, nil
	}

	albumEntity := AlbumEntity{
		Name:      album.Name,
		ArtistRef: artistDoc.Ref,
		Year:      album.Year,
		Signed:    album.Signed,
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

	artistDoc, err := c.ArtistClient.findDocRef(ctx, album.Artist.ID.Primary)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("Error finding existing artist %v while updating album %v", album.Artist.Name, album)
		return err
	}
	if !artistDoc.Exists() {
		log.Errorf("No existing artist %v while updating album %v", album.Artist.Name, album)
		return errors.New("artist does not exist")
	}

	albumEntity := AlbumEntity{
		Name:      album.Name,
		ArtistRef: artistDoc.Ref,
		Year:      album.Year,
		Signed:    album.Signed,
		ID:        album.ID,
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

	artists, err := c.ArtistClient.findAllDocs(ctx)
	if err != nil {
		log.Error("Error retrieving artists while finding all albums,", err)
		return nil, err
	}

	albums := []domain.Album{}
	for _, doc := range albumDocs {
		data := doc.Data()
		album := domain.Album{ID: doc.Ref.ID}
		if name, ok := data["Name"].(string); ok {
			album.Name = name
		}
		if year, ok := data["Year"].(int64); ok {
			album.Year = int(year)
		}
		if signed, ok := data["Signed"].(bool); ok {
			album.Signed = signed
		}
		if artistRef, ok := data["ArtistRef"].(*firestore.DocumentRef); ok {
			album.Artist = (*artists)[artistRef.ID]
		}
		albums = append(albums, album)
	}
	log.Debugf("Found %d albums", len(albums))
	return albums, nil
}

func (c *AlbumClient) findDocRef(ctx context.Context, id string) (*firestore.DocumentSnapshot, error) {
	if id == "" {
		return &firestore.DocumentSnapshot{}, nil
	}
	return c.Connection.Client.Collection(albumCollection).Doc(id).Get(ctx)
}
