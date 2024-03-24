package firestore

import (
	"context"
	"time"

	"concert-manager/data"
	"concert-manager/log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const eventCollection string = "concerts"
var eventFields = []string{"MainActRef", "OpenerRefs", "VenueRef", "Date", "Purchased"}

type EventRepo struct {
	Connection         *Firestore
	VenueRepo  *VenueRepo
	ArtistRepo *ArtistRepo
}

type EventEntity struct {
	MainActRef *firestore.DocumentRef
	OpenerRefs []*firestore.DocumentRef
	VenueRef   *firestore.DocumentRef
	Date       time.Time
	Purchased  bool
}

type Event = data.Event

func (repo *EventRepo) Add(ctx context.Context, event Event) (string, error) {
	log.Debug("Attemping to add event", event)
	var mainActDoc *firestore.DocumentSnapshot
	var err error
	if event.MainAct.Populated() {
		mainActDoc, err = repo.ArtistRepo.findDocRef(ctx, event.MainAct.Name)
		if err != nil {
			log.Errorf("Failed to find existing artist %v while creating event %v", event.MainAct.Name, event)
			return "", err
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
		openerDoc, err := repo.ArtistRepo.findDocRef(ctx, opener.Name)
		if err != nil {
			log.Errorf("Failed to find existing opening artist %v while creating event %v", opener.Name, event)
			return "", err
		}
		log.Debugf("Found existing artist %v with document ID %v while adding event",
			opener.Name, openerDoc.Ref.ID)
		openerRefs = append(openerRefs, openerDoc.Ref)
	}

	venueDoc, err := repo.VenueRepo.findDocRef(ctx, event.Venue.Name, event.Venue.City, event.Venue.State)
	if err != nil {
		log.Errorf("Failed to find existing venue %+v while creating event", event.Venue)
		return "", err
	}
	log.Debugf("Found existing venue %v with document ID %v while adding event", event.Venue, venueDoc.Ref.ID)

	existingEvent, err := repo.findEventDocRef(ctx, event.Date, venueDoc.Ref)
	if err == nil {
		log.Infof("Skipped adding event because it already existed as %+v, %v", event, existingEvent.Ref.ID)
		return existingEvent.Ref.ID, nil
	}
	if err != iterator.Done {
		log.Errorf("Error occurred while checking if artist %v already exists, %v", event, err)
		return "", err
	}

	eventEntity := EventEntity{mainActRef, openerRefs, venueDoc.Ref, data.Timestamp(event.Date), event.Purchased}
	events := repo.Connection.Client.Collection(eventCollection)
	docRef, _, err := events.Add(ctx, eventEntity)
	if err != nil {
		log.Errorf("Failed to add event %+v, %v", event, err)
		return "", err
	}
	log.Infof("Created new event %+v", docRef.ID)
	return docRef.ID, nil
}

func (repo *EventRepo) Delete(ctx context.Context, event Event) error {
	log.Debug("Attemting to delete event", event)
	venueDoc, err := repo.VenueRepo.findDocRef(ctx, event.Venue.Name, event.Venue.City, event.Venue.State)
	if err != nil {
		log.Errorf("Failed to convert venue to ref while removing event %+v", event)
		return err
	}
	log.Debugf("Found existing venue %v with document ID %v while deleting event", event.Venue, venueDoc.Ref.ID)

	eventDoc, err := repo.findEventDocRef(ctx, event.Date, venueDoc.Ref)
	if err != nil {
		log.Errorf("Failed to find existing event while removing %+v", event)
		return err
	}
	log.Debugf("Found existing event document ID %v while deleting event %v", eventDoc.Ref.ID, event)
	eventDoc.Ref.Delete(ctx)
	log.Infof("Successfully deleted event %+v", event)
	return nil
}

func (repo *EventRepo) Exists(ctx context.Context, event Event) (bool, error) {
	log.Debug("Checking for existence of event", event)
	venueDoc, err := repo.VenueRepo.findDocRef(ctx, event.Venue.Name, event.Venue.City, event.Venue.State)
	if err == iterator.Done {
		log.Debugf("No existing venue found while checking for event existence %+v", event)
		return false, nil
	}
	if err != nil {
		log.Errorf("Error while checking existence of event venue %v, %v", event.Venue, err)
		return false, err
	}
	log.Debugf("Found existing venue %v with document ID %v while checking event existence", event.Venue, venueDoc.Ref.ID)

	doc, err := repo.findEventDocRef(ctx, event.Date, venueDoc.Ref)
	if err == iterator.Done {
		log.Debug("No existing event found for", event)
		return false, nil
	}
	if err != nil {
		log.Errorf("Error while checking existence of event %v, %v", event, err)
		return false, err
	}
	log.Debugf("Found event %v with document ID %v", event, doc.Ref.ID)
	return true, nil
}

func (repo *EventRepo) FindAll(ctx context.Context) ([]Event, error) {
	log.Debug("Finding all events")
	eventDocs, err := repo.Connection.Client.Collection(eventCollection).
		Select(eventFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		log.Error("Error while finding all events,", err)
		return nil, err
	}
	log.Debugf("Found %d events", len(eventDocs))

	log.Debug("Finding all artists while finding all events")
	artists, err := repo.ArtistRepo.findAllDocs(ctx)
	if err != nil {
		log.Error("Error retrieving artists while finding all events,", err)
		return nil, err
	}
	log.Debugf("Found %d artists while retrieving all events", len(*artists))

	log.Debug("Finding all venues while finding all events")
	venues, err := repo.VenueRepo.findAllDocs(ctx)
	if err != nil {
		log.Error("Error retrieving venues while finding all events,", err)
		return nil, err
	}
	log.Debugf("Found %d venues while retrieving all events", len(*venues))

	// TODO: This logic could use better error handling for when the firestore event is invalid
	// Currently, the whole app panics if the event data is invalid or the artist or venue is missing
	// It would be better to log an error and ignore invalid events
	events := []Event{}
	for _, e := range eventDocs {
		eventData := e.Data()

		var mainAct Artist
		if mainActRef, ok := eventData["MainActRef"].(*firestore.DocumentRef); ok {
			mainAct = (*artists)[mainActRef.ID]
		}
		venueRef := eventData["VenueRef"].(*firestore.DocumentRef)
		venue := (*venues)[venueRef.ID]

		openers := []Artist{}
		if openerRefs, ok := eventData["OpenerRefs"].([]interface{}); ok {
 			for _, openerRef := range openerRefs {
				openers = append(openers, (*artists)[openerRef.(*firestore.DocumentRef).ID])
			}
		}
		event := Event{
			MainAct: mainAct,
			Openers: openers,
			Venue: venue,
			Date: data.Date(eventData["Date"].(time.Time)),
			Purchased: eventData["Purchased"].(bool),
		}

		events = append(events, event)
	}

	log.Debugf("Returning %d constructed events", len(events))
	return events, nil
}

func (repo *EventRepo) findEventDocRef(ctx context.Context, date string, venueRef *firestore.DocumentRef) (*firestore.DocumentSnapshot, error) {
	event, err := repo.Connection.Client.Collection(eventCollection).
		Select().
		Where("Date", "==", data.Timestamp(date)).
		Where("VenueRef", "==", venueRef).
		Documents(ctx).
		Next()
	if err != nil {
		return nil, err
	}
	return event, nil
}
