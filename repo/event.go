package repo

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"concert-manager/data"
	"concert-manager/db"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const eventCollection string = "concertsHist"
var eventFields = []string{"MainActRef", "OpenerRefs", "VenueRef", "Date"}

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
}

type Event = data.Event

func (repo *EventRepo) Add(ctx context.Context, event Event) (string, error) {
	var (
		mainActDoc *firestore.DocumentSnapshot
		err error
	)
	if event.MainAct.Populated() {
		mainActDoc, err = repo.artistRepo.findDocRef(ctx, event.MainAct.Name)
		if err != nil {
			log.Printf("Failed to find existing artist %+v while creating event", event.MainAct)
			return "", err
		}
	}
	var mainActRef *firestore.DocumentRef
	if mainActDoc != nil {
		mainActRef = mainActDoc.Ref
	}

	openerRefs := []*firestore.DocumentRef{}
	for _, opener := range event.Openers {
		openerDoc, err := repo.artistRepo.findDocRef(ctx, opener.Name)
		if err != nil {
			log.Printf("Failed to find existing opening artist %+v while creating event", opener)
			return "", err
		}
		openerRefs = append(openerRefs, openerDoc.Ref)
	}

	venueDoc, err := repo.venueRepo.findDocRef(ctx, event.Venue.Name, event.Venue.City, event.Venue.State)
	if err != nil {
		log.Printf("Failed to find existing venue %+v while creating event", event.Venue)
		return "", err
	}

	existingEvent, err := repo.findEventDocRef(ctx, event.Date, venueDoc.Ref)
	if err == nil {
		log.Printf("Skipped adding event because it already existed as %+v, %v", event, existingEvent.Ref.ID)
		return existingEvent.Ref.ID, nil
	}
	if err != iterator.Done {
		return "", err
	}

	eventEntity := EventEntity{mainActRef, openerRefs, venueDoc.Ref, toTimestamp(event.Date)}
	events := repo.db.Client.Collection(eventCollection)
	docRef, _, err := events.Add(ctx, eventEntity)
	if err != nil {
		log.Printf("Failed to add event %+v, %v", event, err)
		return "", err
	}
	log.Printf("Created new event %+v", docRef.ID)
	return docRef.ID, nil
}

func (repo *EventRepo) Delete(ctx context.Context, event Event) error {
	venueDoc, err := repo.venueRepo.findDocRef(ctx, event.Venue.Name, event.Venue.City, event.Venue.State)
	if err != nil {
		log.Printf("Failed to convert venue to ref while removing event %+v", event)
		return err
	}

	eventDoc, err := repo.findEventDocRef(ctx, event.Date, venueDoc.Ref)
	if err != nil {
		log.Printf("Failed to find existing event while removing %+v", event)
		return err
	}
	eventDoc.Ref.Delete(ctx)
	log.Printf("Successfully deleted event %+v", event)
	return nil
}

func (repo *EventRepo) Exists(ctx context.Context, event Event) (bool, error) {
	venueDoc, err := repo.venueRepo.findDocRef(ctx, event.Venue.Name, event.Venue.City, event.Venue.State)
	if err == iterator.Done {
		log.Printf("Failed to find existing venue while checking for event existence %+v", event)
		return false, nil
	}
	if err != nil {
		return false, err
	}

	_, err = repo.findEventDocRef(ctx, event.Date, venueDoc.Ref)
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *EventRepo) FindAll(ctx context.Context) (*[]Event, error) {
	eventDocs, err := repo.db.Client.Collection(eventCollection).
		Select(eventFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	artists, err := repo.artistRepo.findAllDocs(ctx)
	if err != nil {
		return nil, err
	}
	venues, err := repo.venueRepo.findAllDocs(ctx)
	if err != nil {
		return nil, err
	}

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
			Date: toDate(eventData["Date"].(time.Time)),
		}

		events = append(events, event)
	}

	return &events, nil
}

func (repo *EventRepo) findEventDocRef(ctx context.Context, date string, venueRef *firestore.DocumentRef) (*firestore.DocumentSnapshot, error) {
	event, err := repo.db.Client.Collection(eventCollection).
		Select().
		Where("Date", "==", toTimestamp(date)).
		Where("VenueRef", "==", venueRef).
		Documents(ctx).
		Next()
	if err != nil {
		return nil, err
	}
	return event, nil
}

// format is "mm/dd/yyyy", with leading zeros optional
// expected that the date string has been previously validated to not error when converted to ints
func toTimestamp(date string) time.Time {
	parts := strings.Split(date, "/")
    month, _ := strconv.Atoi(parts[0])
	day, _ := strconv.Atoi(parts[1])
	year, _ := strconv.Atoi(parts[2])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func toDate(ts time.Time) string {
	day, month, year := ts.Day(), ts.Month(), ts.Year()
	return fmt.Sprintf("%d/%d/%d", month, day, year)
}
