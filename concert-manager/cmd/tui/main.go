package main

import (
	"concert-manager/db"
	"concert-manager/db/firestore"
	"concert-manager/external/lastfm"
	"concert-manager/external/spotify"
	"concert-manager/external/ticketmaster"
	"concert-manager/finder"
	"concert-manager/log"
	"concert-manager/ranker"
	"concert-manager/tui"
	"os"
	"slices"
)

func main() {
	if err := log.Initialize(); err != nil {
		log.Fatal("Failed to set up logger:", err)
	}

	dbConnection, err := firestore.Setup()
	if err != nil {
		log.Fatal("Failed to set up database:", err)
	}

	if slices.Contains(os.Args, "--test") {
		log.Info("Starting in test mode")
		spotify.TEST_MODE = true
	}

	venueClient := &firestore.VenueClient{Connection: dbConnection}
	artistClient := &firestore.ArtistClient{Connection: dbConnection}
	eventClient := &firestore.EventClient{
		Connection:   dbConnection,
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

	spotifyClient := spotify.NewClient()
	lastFmClient := lastfm.NewClient()

	artistRanksCache := &ranker.ArtistRankCache{
		MusicSvc:       spotifyClient,
		ArtistProvider: lastFmClient,
	}
	err = artistRanksCache.InitializeFromFile()
	if err != nil {
		log.Fatal("Failed to initialize artist ranks cache:", err)
	}

	eventRanker := &ranker.EventRanker{
		Cache: artistRanksCache,
	}

	artistInfoFinder := finder.MetadataFinder{
		Spotify: spotifyClient,
		LastFm:  lastFmClient,
	}

	upcomingCache := finder.NewUpcomingEventCache()
	upcomingCache.Finder = eventFinder
	upcomingCache.Ranker = eventRanker
	upcomingCache.SavedDataCache = savedCache
	upcomingCache.MetadataFinder = artistInfoFinder
	err = upcomingCache.InitializeFromFile()
	if err != nil {
		log.Fatal("Failed to initialize upcoming events cache:", err)
	}

	tui.Start(savedCache, upcomingCache)
}
