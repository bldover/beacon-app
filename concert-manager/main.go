package main

import (
	"concert-manager/client/lastfm"
	"concert-manager/client/spotify"
	"concert-manager/db"
	"concert-manager/finder"
	"concert-manager/finder/ticketmaster"
	"concert-manager/loader"
	"concert-manager/log"
	"concert-manager/server"
	"concert-manager/ui"
	"os"
	"slices"
)

func main() {
	if err := log.Initialize(); err != nil {
		log.Fatal("Failed to set up logger:", err)
	}

	dbConnection, err := db.Setup()
	if err != nil {
		log.Fatal("Failed to set up database:", err)
	}

	if slices.Contains(os.Args, "--test") {
		log.Info("Starting in test mode")
		spotify.TEST_MODE = true
	}

	venueClient := &db.VenueClient{Connection: dbConnection}
	artistClient := &db.ArtistClient{Connection: dbConnection}
	eventClient := &db.EventClient{
		Connection: dbConnection,
		VenueClient:  venueClient,
		ArtistClient: artistClient,
	}
	interactor := &db.EventRepository{
		VenueRepo:  venueClient,
		ArtistRepo: artistClient,
		EventRepo:  eventClient,
	}

	savedCache := &db.Cache{}
	savedCache.Database = interactor
	savedCache.LoadCaches()

	ticketmaster := ticketmaster.Ticketmaster{}
	eventFinder := finder.NewEventFinder()
	eventFinder.Ticketmaster = ticketmaster

	eventRanker := &finder.EventRanker{
		MusicSvc: spotify.NewClient(),
		AnalyticsSvc: lastfm.NewClient(),
	}

	upcomingCache := finder.NewUpcomingEventCache()
	upcomingCache.Finder = eventFinder
	upcomingCache.Ranker = eventRanker

	eventLoader := &loader.EventLoader{Cache: savedCache}
	genreLoader := &loader.GenreLoader{Cache: savedCache, InfoProvider: lastfm.NewClient()}

	server := server.Server{}
	server.EventLoader = eventLoader
	server.ArtistInfoLoader = genreLoader
	server.SavedEventCache = savedCache
	server.ArtistCache = savedCache
	server.VenueCache = savedCache
	server.UpcomingEventsCache = upcomingCache
	server.RecommendationCache = upcomingCache

	if slices.Contains(os.Args, "--tui") {
		go server.StartServer()
		ui.Start(savedCache, upcomingCache)
	} else {
		server.StartServer()
	}
}
