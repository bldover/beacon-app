package main

import (
	"concert-manager/cache"
	"concert-manager/data"
	"concert-manager/db"
	"concert-manager/db/firestore"
	finderSvc "concert-manager/finder"
	"concert-manager/loader"
	"concert-manager/log"
	"concert-manager/server"
	"concert-manager/ui/terminal"
	"concert-manager/ui/terminal/screens"
	"concert-manager/ui/terminal/screens/artist"
	"concert-manager/ui/terminal/screens/event"
	"concert-manager/ui/terminal/screens/finder"
	"concert-manager/ui/terminal/screens/venue"
)

func main() {
	if err := log.Initialize(); err != nil {
		log.Fatal("Failed to set up logger:", err)
	}

	dbConnection, err := firestore.Setup()
	if err != nil {
		log.Fatal("Failed to set up database:", err)
	}

	venueRepo := firestore.NewVenueRepo(dbConnection)
	artistRepo := firestore.NewArtistRepo(dbConnection)
	eventRepo := firestore.NewEventRepo(dbConnection, venueRepo, artistRepo)
	interactor := db.NewDatabaseRepository(venueRepo, artistRepo, eventRepo)

	finder := finderSvc.NewEventFinder()
	cache := cache.NewLocalCache(interactor, finder)

	loader := loader.NewLoader(cache)
	go server.StartServer(loader)

	log.Info("Starting terminal UI initialization...")
	startTerminalUI(cache)
}

func startTerminalUI(cache *cache.LocalCache) {
	artistEditScreen := artist.NewEditScreen()
	venueEditScreen := venue.NewEditScreen()
	openerRemoveScreen := artist.NewOpenerRemoveScreen()

	addScreen := event.NewAddScreen(artistEditScreen, venueEditScreen, openerRemoveScreen, cache)
	deleteScreen := event.NewDeleteScreen(cache)
	finderSelectorScreen := finder.NewSelectorScreen(addScreen)

	historyViewScreen := event.NewViewScreen("Concert History", data.Past, addScreen, deleteScreen, cache)
	futureViewScreen := event.NewViewScreen("Future Concerts", data.Future, addScreen, deleteScreen, cache)
	finderViewScreen := finder.NewViewScreen("All Upcoming Events", "Atlanta", "GA", finderSelectorScreen, cache)

	finderMenuScreen := finder.NewMenu(finderViewScreen, nil)

	mainMenuScreen := screens.NewMainMenu()
	mainMenuScreen.Children[1] = historyViewScreen
	mainMenuScreen.Children[2] = futureViewScreen
	mainMenuScreen.Children[3] = finderMenuScreen

	log.Info("Successfully initialized terminal UI, starting display...")
	terminal.RunUI(mainMenuScreen)
}
