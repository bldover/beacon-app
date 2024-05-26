package main

import (
	"concert-manager/cache"
	"concert-manager/db"
	"concert-manager/db/firestore"
	"concert-manager/finder"
	"concert-manager/loader"
	"concert-manager/log"
	"concert-manager/server"
	"concert-manager/ui/textui"
)

func main() {
	if err := log.Initialize(); err != nil {
		log.Fatal("Failed to set up logger:", err)
	}

	dbConnection, err := firestore.Setup()
	if err != nil {
		log.Fatal("Failed to set up database:", err)
	}

	venueRepo := &firestore.VenueRepo{Connection: dbConnection}
	artistRepo := &firestore.ArtistRepo{Connection: dbConnection}
	eventRepo := &firestore.EventRepo{
		Connection: dbConnection,
		VenueRepo:  venueRepo,
		ArtistRepo: artistRepo,
	}
	interactor := &db.DatabaseRepository{
		VenueRepo:  venueRepo,
		ArtistRepo: artistRepo,
		EventRepo:  eventRepo,
	}

	eventFinder := finder.NewEventFinder()
	cache := cache.NewLocalCache()
	cache.Database = interactor
	cache.Finder = eventFinder
	cache.Sync()

	loader := &loader.Loader{Cache: cache}
	go server.StartServer(loader)

	log.Info("Starting terminal UI initialization...")
	textui.Start(cache)
}
