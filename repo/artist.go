package repo

import (
	"concert-manager/data"
	"concert-manager/db"
	"concert-manager/out"
	"context"

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

func (repo *ArtistRepo) Add(ctx context.Context, artist Artist) (string, error) {
	out.Debugf("Attempting to add artist %v", artist)
	existingArtist, err := repo.findDocRef(ctx, artist.Name)
	if err == nil {
		out.Infof("Skipping adding artist because it already exists %+v, %v", artist, existingArtist.Ref.ID)
		return existingArtist.Ref.ID, nil
	}
	if err != iterator.Done {
		out.Errorf("Error occurred while checking if artist %v already exists, %v", artist, err)
		return "", err
	}

	artistEntity := ArtistEntity{artist.Name, artist.Genre}
	artists := repo.db.Client.Collection(artistCollection)
	docRef, _, err := artists.Add(ctx, artistEntity)
	if err != nil {
		out.Errorf("Failed to add new artist %+v, %v", artist, err)
		return "", err
	}
	out.Infof("Created new artist %+v", docRef.ID)
	return docRef.ID, nil
}

func (repo *ArtistRepo) Delete(ctx context.Context, artist Artist) error {
	out.Debugf("Attempting to delete artist %v", artist)
	artistDoc, err := repo.findDocRef(ctx, artist.Name)
	if err != nil {
		out.Errorf("Failed to find existing artist while deleting %+v, %v", artist, err)
		return err
	}
	artistDoc.Ref.Delete(ctx)
	out.Infof("Successfully deleted artist %+v", artist)
	return nil
}

func (repo *ArtistRepo) Exists(ctx context.Context, artist Artist) (bool, error) {
	out.Debugf("Checking for existence of artist %v", artist)
	doc, err := repo.findDocRef(ctx, artist.Name)
	if err == iterator.Done {
		out.Debugf("No existing artist found for %v", artist)
		return false, nil
	}
	if err != nil {
		out.Errorf("Error while checking existence of artist %v, %v", artist, err)
		return false, err
	}
	out.Debugf("Found artist %v with document ID %v", artist, doc.Ref.ID)
	return true, nil
}

func (repo *ArtistRepo) FindAll(ctx context.Context) (*[]Artist, error) {
	out.Debugln("Finding all artists")
	artistDocs, err := repo.db.Client.Collection(artistCollection).
		Select(artistFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		out.Errorf("Error while finding all artists, %v", err)
		return nil, err
	}

	artists := []Artist{}
	for _, a := range artistDocs {
		artists = append(artists, toArtist(a))
	}
	out.Debugf("Found %d artists", len(artists))
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
