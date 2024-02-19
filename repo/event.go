package repo

import (
	"context"
	"time"

	"concert-manager/data"
	"concert-manager/db"
	"concert-manager/out"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const eventCollection string = "concerts"
var eventFields = []string{"MainActRef", "OpenerRefs", "VenueRef", "Date", "Purchased"}

type EventRepo struct {
	db         *db.Firestore
	venueRepo  *VenueRepo
	artistRepo *ArtistRepo
}

func NewEventRepo(fs *db.Firestore, venueRepo *VenueRepo, artistRepo *ArtistRepo) *EventRepo {
	return &EventRepo{fs, venueRepo, artistRepo}
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
	out.Debugf("Attemping to add event %v", event)
	var (
		mainActDoc *firestore.DocumentSnapshot
		err error
	)
	if event.MainAct.Populated() {
		mainActDoc, err = repo.artistRepo.findDocRef(ctx, event.MainAct.Name)
		if err != nil {
			out.Errorf("Failed to find existing artist %v while creating event %v", event.MainAct.Name, event)
			return "", err
		}
		out.Debugf("Found existing artist %v with document ID %v while adding event",
			event.MainAct.Name, mainActDoc.Ref.ID)
	}
	var mainActRef *firestore.DocumentRef
	if mainActDoc != nil {
		mainActRef = mainActDoc.Ref
	}

	openerRefs := []*firestore.DocumentRef{}
	for _, opener := range event.Openers {
		openerDoc, err := repo.artistRepo.findDocRef(ctx, opener.Name)
		if err != nil {
			out.Errorf("Failed to find existing opening artist %v while creating event %v", opener.Name, event)
			return "", err
		}
		out.Debugf("Found existing artist %v with document ID %v while adding event",
			opener.Name, openerDoc.Ref.ID)
		openerRefs = append(openerRefs, openerDoc.Ref)
	}

	venueDoc, err := repo.venueRepo.findDocRef(ctx, event.Venue.Name, event.Venue.City, event.Venue.State)
	if err != nil {
		out.Errorf("Failed to find existing venue %+v while creating event", event.Venue)
		return "", err
	}
	out.Debugf("Found existing venue %v with document ID %v while adding event", event.Venue, venueDoc.Ref.ID)

	existingEvent, err := repo.findEventDocRef(ctx, event.Date, venueDoc.Ref)
	if err == nil {
		out.Infof("Skipped adding event because it already existed as %+v, %v", event, existingEvent.Ref.ID)
		return existingEvent.Ref.ID, nil
	}
	if err != iterator.Done {
		out.Errorf("Error occurred while checking if artist %v already exists, %v", event, err)
		return "", err
	}

	eventEntity := EventEntity{mainActRef, openerRefs, venueDoc.Ref, data.Timestamp(event.Date), event.Purchased}
	events := repo.db.Client.Collection(eventCollection)
	docRef, _, err := events.Add(ctx, eventEntity)
	if err != nil {
		out.Errorf("Failed to add event %+v, %v", event, err)
		return "", err
	}
	out.Infof("Created new event %+v", docRef.ID)
	return docRef.ID, nil
}

func (repo *EventRepo) Delete(ctx context.Context, event Event) error {
	out.Debugf("Attemting to delete event %v", event)
	venueDoc, err := repo.venueRepo.findDocRef(ctx, event.Venue.Name, event.Venue.City, event.Venue.State)
	if err != nil {
		out.Errorf("Failed to convert venue to ref while removing event %+v", event)
		return err
	}
	out.Debugf("Found existing venue %v with document ID %v while deleting event", event.Venue, venueDoc.Ref.ID)

	eventDoc, err := repo.findEventDocRef(ctx, event.Date, venueDoc.Ref)
	if err != nil {
		out.Errorf("Failed to find existing event while removing %+v", event)
		return err
	}
	out.Debugf("Found existing event document ID %v while deleting event %v", eventDoc.Ref.ID, event)
	eventDoc.Ref.Delete(ctx)
	out.Infof("Successfully deleted event %+v", event)
	return nil
}

func (repo *EventRepo) Exists(ctx context.Context, event Event) (bool, error) {
	out.Debugf("Checking for existence of event %v", event)
	venueDoc, err := repo.venueRepo.findDocRef(ctx, event.Venue.Name, event.Venue.City, event.Venue.State)
	if err == iterator.Done {
		out.Debugf("No existing venue found while checking for event existence %+v", event)
		return false, nil
	}
	if err != nil {
		out.Errorf("Error while checking existence of event venue %v, %v", event.Venue, err)
		return false, err
	}
	out.Debugf("Found existing venue %v with document ID %v while checking event existence", event.Venue, venueDoc.Ref.ID)

	doc, err := repo.findEventDocRef(ctx, event.Date, venueDoc.Ref)
	if err == iterator.Done {
		out.Debugf("No existing event found for %v", event)
		return false, nil
	}
	if err != nil {
		out.Errorf("Error while checking existence of event %v, %v", event, err)
		return false, err
	}
	out.Debugf("Found event %v with document ID %v", event, doc.Ref.ID)
	return true, nil
}

func (repo *EventRepo) FindAll(ctx context.Context) (*[]Event, error) {
	out.Debugln("Finding all events")
	eventDocs, err := repo.db.Client.Collection(eventCollection).
		Select(eventFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		out.Errorf("Error while finding all events, %v", err)
		return nil, err
	}
	out.Debugf("Found %d events", len(eventDocs))

	out.Debugln("Finding all artists while finding all events")
	artists, err := repo.artistRepo.findAllDocs(ctx)
	if err != nil {
		out.Errorf("Error retrieving artists while finding all events", err)
		return nil, err
	}
	out.Debugf("Found %d artists while retrieving all events", len(*artists))

	out.Debugln("Finding all venues while finding all events")
	venues, err := repo.venueRepo.findAllDocs(ctx)
	if err != nil {
		out.Errorf("Error retrieving venues while finding all events", err)
		return nil, err
	}
	out.Debugf("Found %d venues while retrieving all events", len(*venues))

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

	out.Debugf("Returning %d constructed events", len(events))
	return &events, nil
}

func (repo *EventRepo) findEventDocRef(ctx context.Context, date string, venueRef *firestore.DocumentRef) (*firestore.DocumentSnapshot, error) {
	event, err := repo.db.Client.Collection(eventCollection).
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
