package repo

import (
	"concert-manager/data"
	"concert-manager/db"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const venueCollection = "venues"
var venueFields = []string{"Name", "City", "State"}

type VenueRepo struct {
	db *db.Firestore
}

func NewVenueRepo(fs *db.Firestore) *VenueRepo {
	return &VenueRepo{fs}
}

type VenueEntity struct {
	Name    string
	City    string
	State   string
}

type Venue = data.Venue

func (repo *VenueRepo) Add(ctx context.Context, venue Venue) (string, error) {
	existingVenue, err := repo.findDocRef(ctx, venue.Name, venue.City, venue.State)
	if err == nil {
		log.Printf("Skipping adding venue because it already exists %+v, %v", venue, existingVenue.Ref.ID)
		return existingVenue.Ref.ID, nil
	}
	if err != iterator.Done {
		return "", err
	}

	venueEntity := VenueEntity{venue.Name, venue.City, venue.State}
	venues := repo.db.Client.Collection(venueCollection)
	docRef, _, err := venues.Add(ctx, venueEntity)
	if err != nil {
		log.Printf("Failed to add new venue %+v, %v", venue, err)
		return "", err
	}
	log.Printf("Created new venue %+v", docRef.ID)
	return docRef.ID, nil
}

func (repo *VenueRepo) Delete(ctx context.Context, venue Venue) error {
	venueDoc, err := repo.findDocRef(ctx, venue.Name, venue.City, venue.State)
	if err != nil {
		log.Printf("Failed to find existing venue while deleting %+v, %v", venue, err)
		return err
	}
	venueDoc.Ref.Delete(ctx)
	log.Printf("Successfully deleted venue %+v", venue)
	return nil
}

func (repo *VenueRepo) Exists(ctx context.Context, venue Venue) (bool, error) {
	_, err := repo.findDocRef(ctx, venue.Name, venue.City, venue.State)
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *VenueRepo) FindAll(ctx context.Context) (*[]Venue, error) {
	venueDocs, err := repo.db.Client.Collection(venueCollection).
		Select(venueFields...).
		Documents(ctx).
	 	GetAll()
	if err != nil {
		return nil, err
	}

	venues := []Venue{}
	for _, v := range venueDocs {
		venues = append(venues, toVenue(v))
	}
	return &venues, nil
}

func toVenue(doc *firestore.DocumentSnapshot) Venue {
    venueData := doc.Data()
	return Venue{
		Name:    venueData["Name"].(string),
		City:    venueData["City"].(string),
		State:   venueData["State"].(string),
	}
}

func (repo *VenueRepo) findDocRef(ctx context.Context, name string, city string, state string) (*firestore.DocumentSnapshot, error) {
	return repo.db.Client.Collection(venueCollection).
		Select().
		Where("Name", "==", name).
		Where("City", "==", city).
		Where("State", "==", state).
		Documents(ctx).
		Next()
}

func (repo *VenueRepo) findAllDocs(ctx context.Context) (*map[string]Venue, error) {
	venueDocs, err := repo.db.Client.Collection(venueCollection).
		Select(venueFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	venues := make(map[string]Venue)
	for _, v := range venueDocs {
		venues[v.Ref.ID] = toVenue(v)
	}

	return &venues, nil
}
