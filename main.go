package main

import (
	"concert-manager/data"
	"concert-manager/db"
	"concert-manager/db/firestore"
	"concert-manager/loader"
	"concert-manager/log"
	"concert-manager/server"
	"concert-manager/ui/terminal"
	"concert-manager/ui/terminal/screens"
	"context"
)

func main() {
	if err := log.Initialize(); err != nil {
		log.Fatal("Failed to set up logger:", err)
	}
	log.Info("Successfully initialized logger")

	dbConnection, err := firestore.Setup()
	if err != nil {
		log.Fatal("Failed to set up database:", err)
	}
	log.Info("Successfully initialized database")

	venueRepo := firestore.NewVenueRepo(dbConnection)
	artistRepo := firestore.NewArtistRepo(dbConnection)
	eventRepo := firestore.NewEventRepo(dbConnection, venueRepo, artistRepo)

	interactor := &db.DatabaseRepository{}
	interactor.VenueRepo = venueRepo
	interactor.ArtistRepo = artistRepo
	interactor.EventRepo = eventRepo

	loader := &loader.Loader{EventCreator: interactor}
	go server.StartServer(loader)

	log.Info("Starting event initialization...")
	events, err := interactor.ListEvents(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize events:", err)
	}
	log.Info("Successfully initialized events")
	pastEvents := []data.Event{}
	futureEvents := []data.Event{}
	for _, e := range *events {
		if data.ValidFutureDate(e.Date) {
			futureEvents = append(futureEvents, e)
		} else {
			pastEvents = append(pastEvents, e)
		}
	}

	log.Info("Starting artist initialization...")
	artists, err := interactor.ListArtists(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize artists:", err)
	}
	log.Info("Successfully initialized artists")

	log.Info("Starting venue initialization...")
	venues, err := interactor.ListVenues(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize venues:", err)
	}
	log.Info("Successfully initialized venues")

	log.Info("Starting terminal UI initialization...")
	mainMenu := setupTerminalUI(interactor, &pastEvents, &futureEvents, artists, venues)
	log.Info("Successfully initialized terminal UI, starting display...")
	terminal.RunUI(mainMenu)
}

func setupTerminalUI(dbRepo *db.DatabaseRepository, pastEvents *[]data.Event, futureEvents *[]data.Event,
	artists *[]data.Artist, venues *[]data.Venue) *screens.MainMenu {

	mainMenuScreen := screens.NewMainMenu()

	historyViewScreen := screens.NewEventViewScreen("Concert History")
	historyDeleteScreen := screens.NewEventDeleteScreen()
	historyAddScreen := screens.NewEventAddScreen(false)
	artistEditScreen := screens.NewArtistEditScreen()
	venueEditScreen := screens.NewVenueEditScreen()
	openerRemoveScreen := screens.NewOpenerRemoveScreen()

	historyViewScreen.Events = pastEvents
	historyViewScreen.MainMenu = mainMenuScreen
	historyViewScreen.AddEventScreen = historyAddScreen
	historyViewScreen.DeleteEventScreen = historyDeleteScreen

	historyAddScreen.Events = pastEvents
	historyAddScreen.Database = dbRepo
	historyAddScreen.ArtistEditor = artistEditScreen
	historyAddScreen.VenueEditor = venueEditScreen
	historyAddScreen.OpenerRemover = openerRemoveScreen
	historyAddScreen.Viewer = historyViewScreen

	historyDeleteScreen.Events = pastEvents
	historyDeleteScreen.Viewer = historyViewScreen
	historyDeleteScreen.Database = dbRepo

	futureViewScreen := screens.NewEventViewScreen("Future Concerts")
	futureDeleteScreen := screens.NewEventDeleteScreen()
	futureAddScreen := screens.NewEventAddScreen(true)

	futureViewScreen.Events = futureEvents
	futureViewScreen.MainMenu = mainMenuScreen
	futureViewScreen.AddEventScreen = futureAddScreen
	futureViewScreen.DeleteEventScreen = futureDeleteScreen

	futureAddScreen.Events = futureEvents
	futureAddScreen.Database = dbRepo
	futureAddScreen.ArtistEditor = artistEditScreen
	futureAddScreen.VenueEditor = venueEditScreen
	futureAddScreen.OpenerRemover = openerRemoveScreen
	futureAddScreen.Viewer = futureViewScreen

	futureDeleteScreen.Events = futureEvents
	futureDeleteScreen.Viewer = futureViewScreen
	futureDeleteScreen.Database = dbRepo

	artistEditScreen.Artists = artists
	artistEditScreen.ArtistAdder = dbRepo
	artistEditScreen.AddEventScreen = historyAddScreen

	venueEditScreen.Venues = venues
	venueEditScreen.VenueAdder = dbRepo
	venueEditScreen.AddEventScreen = historyAddScreen

	openerRemoveScreen.AddEventScreen = historyAddScreen

	mainMenuScreen.Children[1] = historyViewScreen
	mainMenuScreen.Children[2] = futureViewScreen

	return mainMenuScreen
}
