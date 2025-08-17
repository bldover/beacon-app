package firestore

import (
	"context"
	"errors"
	"time"

	"concert-manager/domain"
	"concert-manager/log"
	"concert-manager/util"

	"cloud.google.com/go/firestore"
)

const eventCollection string = "events"

var eventFields = []string{"MainActRef", "OpenerRefs", "VenueRef", "Date", "Purchased"}

type (
	EventClient struct {
		Connection   *Firestore
		VenueClient  *VenueClient
		ArtistClient *ArtistClient
	}

	EventEntity struct {
		MainActRef *firestore.DocumentRef
		OpenerRefs []*firestore.DocumentRef
		VenueRef   *firestore.DocumentRef
		Date       time.Time
		Purchased  bool
		ID         EventIDEntity
	}

	EventIDEntity struct {
		Primary      string
		Ticketmaster string
	}
)

func (c *EventClient) Add(ctx context.Context, event domain.Event) (string, error) {
	log.Debug("Attemping to add event", event)
	var mainActDoc *firestore.DocumentSnapshot
	var err error
	if event.MainAct.Populated() {
		mainActDoc, err = c.ArtistClient.findDocRef(ctx, event.MainAct.ID.Primary)
		if err != nil {
			log.Errorf("Error finding existing artist %v while creating event %v", event.MainAct.Name, event)
			return "", err
		}
		if !mainActDoc.Exists() {
			log.Errorf("No existing artist %v while creating event %v", event.MainAct.Name, event)
			return "", errors.New("main artist does not exist")
		}
		log.Debugf("Found existing artist %v with document ID %v while adding event",
			event.MainAct.Name, mainActDoc.Ref.ID)
	}
	var mainActRef *firestore.DocumentRef
	if mainActDoc != nil {
		mainActRef = mainActDoc.Ref
	}

	openerRefs := []*firestore.DocumentRef{}
	for _, opener := range event.Openers {
		openerDoc, err := c.ArtistClient.findDocRef(ctx, opener.ID.Primary)
		if err != nil {
			log.Errorf("Error finding existing opening artist %v while creating event %v", opener.Name, event)
			return "", err
		}
		if !openerDoc.Exists() {
			log.Errorf("No existing opening artist %v while creating event %v", opener.Name, event)
			return "", errors.New("opering artist does not exist")
		}
		log.Debugf("Found existing artist %v with document ID %v while adding event",
			opener.Name, openerDoc.Ref.ID)
		openerRefs = append(openerRefs, openerDoc.Ref)
	}

	venueDoc, err := c.VenueClient.findDocRef(ctx, event.Venue.ID.Primary)
	if err != nil {
		log.Errorf("Error finding existing venue %+v while creating event", event.Venue)
		return "", err
	}
	if !venueDoc.Exists() {
		log.Errorf("No existing venue %+v while creating event", event.Venue)
		return "", errors.New("venue does not exist")
	}
	log.Debugf("Found existing venue %v with document ID %v while adding event", event.Venue, venueDoc.Ref.ID)

	existingEvent, err := c.findEventDocRef(ctx, event.ID.Primary)
	if err != nil {
		log.Errorf("Error occurred while checking if event %v already exists, %v", event, err)
		return "", err
	}
	if existingEvent.Exists() {
		log.Debugf("Skipped adding event because it already existed as %+v", event)
		return existingEvent.Ref.ID, nil
	}

	idEntity := EventIDEntity{
		Primary:      event.ID.Primary,
		Ticketmaster: event.ID.Ticketmaster,
	}
	eventEntity := EventEntity{mainActRef, openerRefs, venueDoc.Ref, util.Timestamp(event.Date), event.Purchased, idEntity}

	events := c.Connection.Client.Collection(eventCollection)
	docRef, _, err := events.Add(ctx, eventEntity)
	if err != nil {
		log.Errorf("Failed to add event %+v, %v", event, err)
		return "", err
	}
	log.Infof("Created new event %+v", docRef.ID)
	return docRef.ID, nil
}

func (c *EventClient) Delete(ctx context.Context, id string) error {
	log.Debug("Attemting to delete event", id)
	eventDoc, err := c.Connection.Client.Collection(eventCollection).Doc(id).Get(ctx)
	if err != nil {
		log.Errorf("Error while deleting event %s", id)
		return err
	}
	_, err = eventDoc.Ref.Delete(ctx)
	if err != nil {
		log.Error("Failed to delete event", id, err)
		return err
	}
	log.Infof("Successfully deleted event %+v", id)
	return nil
}

func (c *EventClient) FindAll(ctx context.Context) ([]domain.Event, error) {
	log.Debug("Finding all events")
	eventDocs, err := c.Connection.Client.Collection(eventCollection).
		Select(eventFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		log.Error("Error while finding all events,", err)
		return nil, err
	}
	log.Debugf("Found %d events", len(eventDocs))

	log.Debug("Finding all artists while finding all events")
	artists, err := c.ArtistClient.findAllDocs(ctx)
	if err != nil {
		log.Error("Error retrieving artists while finding all events,", err)
		return nil, err
	}
	log.Debugf("Found %d artists while retrieving all events", len(*artists))

	log.Debug("Finding all venues while finding all events")
	venues, err := c.VenueClient.findAllDocs(ctx)
	if err != nil {
		log.Error("Error retrieving venues while finding all events,", err)
		return nil, err
	}
	log.Debugf("Found %d venues while retrieving all events", len(*venues))

	// TODO: This logic could use better error handling for when the firestore event is invalid
	// Currently, the whole app panics if the event data is invalid or the artist or venue is missing
	// It would be better to log an error and ignore invalid events
	events := []domain.Event{}
	for _, e := range eventDocs {
		eventData := e.Data()

		var mainAct domain.Artist
		if mainActRef, ok := eventData["MainActRef"].(*firestore.DocumentRef); ok {
			mainAct = (*artists)[mainActRef.ID]
		}
		venueRef := eventData["VenueRef"].(*firestore.DocumentRef)
		venue := (*venues)[venueRef.ID]

		openers := []domain.Artist{}
		if openerRefs, ok := eventData["OpenerRefs"].([]interface{}); ok {
			for _, openerRef := range openerRefs {
				openers = append(openers, (*artists)[openerRef.(*firestore.DocumentRef).ID])
			}
		}

		event := domain.Event{
			MainAct:   &mainAct,
			Openers:   openers,
			Venue:     venue,
			Date:      util.Date(eventData["Date"].(time.Time)),
			Purchased: eventData["Purchased"].(bool),
			ID:        domain.ID{Primary: e.Ref.ID},
		}

		if ids, ok := e.Data()["ID"].(map[string]any); ok {
			if ticketmasterId, ok := ids["Ticketmaster"].(string); ok {
				event.ID.Ticketmaster = ticketmasterId
			}
		}

		events = append(events, event)
	}

	log.Debugf("Returning %d constructed events", len(events))
	return events, nil
}

func (c *EventClient) findEventDocRef(ctx context.Context, id string) (*firestore.DocumentSnapshot, error) {
	if id == "" {
		return &firestore.DocumentSnapshot{}, nil
	}
	return c.Connection.Client.Collection(eventCollection).Doc(id).Get(ctx)
}
