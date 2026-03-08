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

var albumFields = []string{"Name", "ArtistRefs", "Year", "Signed", "Wishlisted", "Variant", "Format", "Notes", "CoverImageUrl", "ID"}

type (
	AlbumClient struct {
		Connection   *Firestore
		ArtistClient *ArtistClient
	}

	AlbumEntity struct {
		Name          string
		ArtistRefs    []*firestore.DocumentRef
		Year          int
		Signed        bool
		Wishlisted    bool
		Variant       string
		Format        string
		Notes         string
		CoverImageUrl string
		ID            string
	}
)

func (c *AlbumClient) Add(ctx context.Context, album domain.Album) (string, error) {
	log.Debug("Attempting to add album", album)
	artistRefs := make([]*firestore.DocumentRef, 0, len(album.Artists))
	for _, artist := range album.Artists {
		artistDoc, err := c.ArtistClient.findDocRef(ctx, artist.ID.Primary)
		if err != nil && status.Code(err) != codes.NotFound {
			log.Errorf("Error finding existing artist %v while adding album %v", artist.Name, album)
			return "", err
		}
		if !artistDoc.Exists() {
			log.Errorf("No existing artist %v while adding album %v", artist.Name, album)
			return "", errors.New("artist does not exist")
		}
		artistRefs = append(artistRefs, artistDoc.Ref)
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
		Name:          album.Name,
		ArtistRefs:    artistRefs,
		Year:          album.Year,
		Signed:        album.Signed,
		Wishlisted:    album.Wishlisted,
		Variant:       album.Variant,
		Format:        album.Format,
		Notes:         album.Notes,
		CoverImageUrl: album.CoverImageUrl,
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

	artistRefs := make([]*firestore.DocumentRef, 0, len(album.Artists))
	for _, artist := range album.Artists {
		artistDoc, err := c.ArtistClient.findDocRef(ctx, artist.ID.Primary)
		if err != nil && status.Code(err) != codes.NotFound {
			log.Errorf("Error finding existing artist %v while updating album %v", artist.Name, album)
			return err
		}
		if !artistDoc.Exists() {
			log.Errorf("No existing artist %v while updating album %v", artist.Name, album)
			return errors.New("artist does not exist")
		}
		artistRefs = append(artistRefs, artistDoc.Ref)
	}

	albumEntity := AlbumEntity{
		Name:          album.Name,
		ArtistRefs:    artistRefs,
		Year:          album.Year,
		Signed:        album.Signed,
		Wishlisted:    album.Wishlisted,
		Variant:       album.Variant,
		Format:        album.Format,
		Notes:         album.Notes,
		CoverImageUrl: album.CoverImageUrl,
		ID:            album.ID,
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
		if wishlisted, ok := data["Wishlisted"].(bool); ok {
			album.Wishlisted = wishlisted
		}
		if variant, ok := data["Variant"].(string); ok {
			album.Variant = variant
		}
		if format, ok := data["Format"].(string); ok {
			album.Format = format
		}
		if notes, ok := data["Notes"].(string); ok {
			album.Notes = notes
		}
		if coverImageUrl, ok := data["CoverImageUrl"].(string); ok {
			album.CoverImageUrl = coverImageUrl
		}
		if artistRefs, ok := data["ArtistRefs"].([]interface{}); ok {
			for _, ref := range artistRefs {
				if docRef, ok := ref.(*firestore.DocumentRef); ok {
					if artist, exists := (*artists)[docRef.ID]; exists {
						album.Artists = append(album.Artists, artist)
					}
				}
			}
		}
		if album.Artists == nil {
			album.Artists = []domain.Artist{}
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
