package firestore

import (
	"concert-manager/domain"
	"concert-manager/log"
	"context"

	"cloud.google.com/go/firestore"
)

const venueCollection = "venues"

var venueFields = []string{"Name", "City", "State"}

type (
	VenueClient struct {
		Connection *Firestore
	}

	VenueEntity struct {
		Name  string
		City  string
		State string
		ID    VenueIDEntity
	}

	VenueIDEntity struct {
		Primary      string
		Ticketmaster string
	}
)

func (c *VenueClient) Add(ctx context.Context, venue domain.Venue) (string, error) {
	log.Debug("Attemping to add venue", venue)
	existingVenue, err := c.findDocRef(ctx, venue.ID.Primary)
	if err != nil {
		log.Errorf("Error occurred while checking if venue %v already exists, %v", venue, err)
		return "", err
	}
	if existingVenue.Exists() {
		log.Debugf("Skipping adding venue because it already exists %+v, %v", venue, existingVenue.Ref.ID)
		return existingVenue.Ref.ID, nil
	}

	idEntity := VenueIDEntity{
		Primary:      venue.ID.Primary,
		Ticketmaster: venue.ID.Ticketmaster,
	}
	venueEntity := VenueEntity{venue.Name, venue.City, venue.State, idEntity}

	venues := c.Connection.Client.Collection(venueCollection)
	docRef, _, err := venues.Add(ctx, venueEntity)
	if err != nil {
		log.Errorf("Failed to add new venue %+v, %v", venue, err)
		return "", err
	}
	log.Infof("Created new venue %+v", docRef.ID)
	return docRef.ID, nil
}

func (c *VenueClient) Update(ctx context.Context, venue domain.Venue) error {
	log.Debug("Attempting to update venue", venue)
	venueDoc, err := c.Connection.Client.Collection(venueCollection).Doc(venue.ID.Primary).Get(ctx)
	if err != nil {
		log.Errorf("Error while updating venue %+v, %v", venue, err)
		return err
	}
	if !venueDoc.Exists() {
		log.Errorf("Venue does not exist in update %+v, %v", venue, err)
		return err
	}

	idEntity := VenueIDEntity{
		Primary:      venue.ID.Primary,
		Ticketmaster: venue.ID.Ticketmaster,
	}
	venueEntity := VenueEntity{venue.Name, venue.City, venue.State, idEntity}

	_, err = venueDoc.Ref.Set(ctx, venueEntity)
	if err != nil {
		log.Errorf("Failed to update venue %+v, %v", venue, err)
		return err
	}
	log.Info("Successfully updated venue", venue)
	return nil
}

func (c *VenueClient) Delete(ctx context.Context, id string) error {
	log.Debug("Attemping to delete venue", id)
	venueDoc, err := c.Connection.Client.Collection(venueCollection).Doc(id).Get(ctx)
	if err != nil {
		log.Error("Error while deleting venue", id, err)
		return err
	}
	_, err = venueDoc.Ref.Delete(ctx)
	if err != nil {
		log.Error("Failed to delete venue", id, err)
		return err
	}
	log.Info("Successfully deleted venue", id)
	return nil
}

func (c *VenueClient) FindAll(ctx context.Context) ([]domain.Venue, error) {
	log.Debug("Finding all venues")
	venueDocs, err := c.Connection.Client.Collection(venueCollection).
		Select(venueFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		log.Error("Error while finding all venues,", err)
		return nil, err
	}

	venues := []domain.Venue{}
	for _, v := range venueDocs {
		venues = append(venues, toVenue(v))
	}
	log.Debugf("Found %d artists", len(venues))
	return venues, nil
}

func toVenue(doc *firestore.DocumentSnapshot) domain.Venue {
	venueData := doc.Data()
	venue := domain.Venue{
		Name:  venueData["Name"].(string),
		City:  venueData["City"].(string),
		State: venueData["State"].(string),
		ID:    domain.ID{Primary: doc.Ref.ID},
	}

	if ids, ok := doc.Data()["ID"].(map[string]any); ok {
		if ticketmasterId, ok := ids["Ticketmaster"].(string); ok {
			venue.ID.Ticketmaster = ticketmasterId
		}
	}

	return venue
}

func (c *VenueClient) findDocRef(ctx context.Context, id string) (*firestore.DocumentSnapshot, error) {
	if id == "" {
		return &firestore.DocumentSnapshot{}, nil
	}
	return c.Connection.Client.Collection(venueCollection).Doc(id).Get(ctx)
}

func (c *VenueClient) findAllDocs(ctx context.Context) (*map[string]domain.Venue, error) {
	venueDocs, err := c.Connection.Client.Collection(venueCollection).
		Select(venueFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	venues := make(map[string]domain.Venue)
	for _, v := range venueDocs {
		venues[v.Ref.ID] = toVenue(v)
	}

	return &venues, nil
}
