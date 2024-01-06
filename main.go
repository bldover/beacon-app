package main

import (
	"concert-manager/db"
	"concert-manager/repo"
	"concert-manager/server"
	"concert-manager/svc"
	"log"
)

func main() {
	firestore, err := db.Setup()
	if err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}

	venueRepo := repo.NewVenueRepo(firestore)
	artistRepo := repo.NewArtistRepo(firestore)
	eventRepo := repo.NewEventRepo(firestore, venueRepo, artistRepo)

	interactor := &svc.EventInteractor{}
	interactor.VenueRepo = venueRepo
	interactor.ArtistRepo = artistRepo
	interactor.EventRepo = eventRepo

	loader := &svc.Loader{EventCreator: interactor}
	server.StartServer(loader)

	// ADD integration test
	// ctx := context.Background()
	// venue := data.Venue{Name: "name", City: "city", State: "state"}
	// artist := data.Artist{Name: "name", Genre: "genre"}
	// artist2 := data.Artist{Name: "name2", Genre: "genre2"}
	// event := data.Event{MainAct: artist, Openers: []data.Artist{artist}, Venue: venue, Date: "1/1/2000"}
	// event2 := data.Event{Openers: []data.Artist{artist}, Venue: venue, Date: "1/1/2001"}
	// event3 := data.Event{MainAct: artist, Venue: venue, Date: "1/2/2002"}
	// if err := interactor.AddVenue(ctx, venue); err != nil {
	// 	log.Printf("failed to add venue %v", err)
	// }
	// if err := interactor.AddArtist(ctx, artist); err != nil {
	// 	log.Printf("failed to add artist %v", err)
	// }
	// if err := interactor.AddArtist(ctx, artist2); err != nil {
	// 	log.Printf("failed to add artist2 %v", err)
	// }
	// if err := interactor.AddEvent(ctx, event); err != nil {
	// 	log.Printf("failed to add event %v", err)
	// }
	// if err := interactor.AddEvent(ctx, event2); err != nil {
	// 	log.Printf("failed to add event2 %v", err)
	// }
	// if err := interactor.AddEvent(ctx, event3); err != nil {
	// 	log.Printf("failed to add event3 %v", err)
	// }

	// venues, err := interactor.ListVenues(ctx)
	// if err != nil {
	// 	log.Printf("failed to read venues %v", err)
	// }
	// log.Printf("Found venues %+v", venues)
	// artists, err := interactor.ListArtists(ctx)
	// if err != nil {
	// 	log.Printf("failed to read artists %v", err)
	// }
	// log.Printf("Found artists %+v", artists)
	// events, err := interactor.ListEvents(ctx)
	// if err != nil {
	// 	log.Printf("failed to read events %v", err)
	// }
	// log.Printf("Found events %+v", events)

	// // DELETE integration test
	// if err := interactor.DeleteEvent(ctx, event); err != nil {
	// 	log.Printf("failed to delete event %v", err)
	// }
	// if err := interactor.DeleteEvent(ctx, event2); err != nil {
	// 	log.Printf("failed to delete event %v", err)
	// }
	// if err := interactor.DeleteEvent(ctx, event3); err != nil {
	// 	log.Printf("failed to delete event %v", err)
	// }
	// if err := interactor.DeleteVenue(ctx, venue); err != nil {
	// 	log.Printf("failed to delete venue %v", err)
	// }
	// if err := interactor.DeleteArtist(ctx, artist); err != nil {
	// 	log.Printf("failed to delete artist 2 %v", err)
	// }
	// if err := interactor.DeleteArtist(ctx, artist2); err != nil {
	// 	log.Printf("failed to delete artist 3 %v", err)
	// }
}
