package main

import (
	"concert-manager/cache"
	"concert-manager/db"
	"concert-manager/db/firestore"
	"concert-manager/finder"
	"concert-manager/loader"
	"concert-manager/log"
	"concert-manager/ranker"
	"concert-manager/server"
	"concert-manager/spotify"
	"concert-manager/ui"
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

	savedCache := &cache.SavedEventCache{}
	savedCache.Database = interactor
	savedCache.Sync()

	eventFinder := finder.NewEventFinder()
	artistRanker := ranker.ArtistRanker{MusicSvc: spotify.NewClient()}
	eventRanker := &ranker.EventRanker{ArtistRanker: artistRanker}

	upcomingCache := cache.NewUpcomingEventCache()
	upcomingCache.Finder = eventFinder
	upcomingCache.Ranker = eventRanker

	loader := &loader.Loader{Cache: savedCache}
	go server.StartServer(loader)

	ui.StartUI("", savedCache, upcomingCache)
}
