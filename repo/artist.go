package repo

import (
	"concert-manager/data"
	"concert-manager/db"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const artistCollection string = "artists"
var artistFields = []string{"Name", "Genre"}

type ArtistRepo struct {
	db *db.Firestore
}

func NewArtistRepo(fs *db.Firestore) *ArtistRepo {
	return &ArtistRepo{fs}
}

type ArtistEntity struct {
	Name  string
	Genre string
}

type Artist = data.Artist

func (repo *ArtistRepo) Add(ctx context.Context, artist Artist) (id string, err error) {
	existingArtist, err := repo.findDocRef(ctx, artist.Name)
	if err == nil {
		log.Printf("Skipping adding artist because it already exists %+v, %v", artist, existingArtist.Ref.ID)
		return existingArtist.Ref.ID, nil
	}
	if err != iterator.Done {
		return
	}

	artistEntity := ArtistEntity{artist.Name, artist.Genre}
	artists := repo.db.Client.Collection(artistCollection)
	docRef, _, err := artists.Add(ctx, artistEntity)
	if err != nil {
		log.Printf("Failed to add new artist %+v, %v", artist, err)
		return
	}
	log.Printf("Created new artist %+v", docRef.ID)
	return docRef.ID, nil
}

func (repo *ArtistRepo) Delete(ctx context.Context, artist Artist) error {
	artistDoc, err := repo.findDocRef(ctx, artist.Name)
	if err != nil {
		log.Printf("Failed to find existing artist while deleting %+v, %v", artist, err)
		return err
	}
	artistDoc.Ref.Delete(ctx)
	log.Printf("Successfully deleted artist %+v", artist)
	return nil
}

func (repo *ArtistRepo) Exists(ctx context.Context, artist Artist) (bool, error) {
	_, err := repo.findDocRef(ctx, artist.Name)
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *ArtistRepo) FindAll(ctx context.Context) (*[]Artist, error) {
	artistDocs, err := repo.db.Client.Collection(artistCollection).
		Select(artistFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	artists := []Artist{}
	for _, a := range artistDocs {
		artists = append(artists, toArtist(a))
	}
	return &artists, nil
}

func toArtist(doc *firestore.DocumentSnapshot) Artist {
    artistData := doc.Data()
	return Artist{
		Name:  artistData["Name"].(string),
		Genre: artistData["Genre"].(string),
	}
}

func (repo *ArtistRepo) findDocRef(ctx context.Context, name string) (*firestore.DocumentSnapshot, error) {
	artist, err := repo.db.Client.Collection(artistCollection).
		Select().
		Where("Name", "==", name).
		Documents(ctx).
		Next()
	if err != nil {
		return nil, err
	}
	return artist, nil
}

func (repo *ArtistRepo) findAllDocs(ctx context.Context) (*map[string]Artist, error) {
	artistDocs, err := repo.db.Client.Collection(artistCollection).
		Select(artistFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	artists := make(map[string]Artist)
	for _, a := range artistDocs {
		artists[a.Ref.ID] = toArtist(a)
	}

	return &artists, nil
}
