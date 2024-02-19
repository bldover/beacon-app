package main

import (
	"concert-manager/cli"
	"concert-manager/cli/artist"
	"concert-manager/cli/event"
	"concert-manager/cli/venue"
	"concert-manager/data"
	"concert-manager/db"
	"concert-manager/out"
	"concert-manager/repo"
	"concert-manager/server"
	"concert-manager/svc"
	"context"
	"log"
)

func main() {
	if err := out.Initialize(); err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}
	out.Infoln("Successfully initialized logger")

	firestore, err := db.Setup()
	if err != nil {
	 	out.Fatalf("Failed to set up database: %v", err)
	}
	out.Infoln("Successfully initialized database")

	venueRepo := repo.NewVenueRepo(firestore)
	artistRepo := repo.NewArtistRepo(firestore)
	eventRepo := repo.NewEventRepo(firestore, venueRepo, artistRepo)

	interactor := &svc.EventInteractor{}
	interactor.VenueRepo = venueRepo
	interactor.ArtistRepo = artistRepo
	interactor.EventRepo = eventRepo

	loader := &svc.Loader{EventCreator: interactor}
	go server.StartServer(loader)

	out.Infoln("Starting event initializition...")
	events, err := interactor.ListEvents(context.Background())
	if err != nil {
		out.Fatalf("Failed to initialize events: %v", err)
	}
	out.Infoln("Successfully initialized events")
	pastEvents := []data.Event{}
	futureEvents := []data.Event{}
	for _, e := range *events {
		if data.ValidFutureDate(e.Date) {
			futureEvents = append(futureEvents, e)
		} else {
			pastEvents = append(pastEvents, e)
		}
	}

	out.Infoln("Starting artist initialization...")
	artists, err := interactor.ListArtists(context.Background())
	if err != nil {
		out.Fatalf("Failed to initialize artists: %v", err)
	}
	out.Infoln("Successfully initialized artists")

	out.Infoln("Starting venue initialization...")
	venues, err := interactor.ListVenues(context.Background())
	if err != nil {
		out.Fatalf("Failed to initialize venues: %v", err)
	}
	out.Infoln("Successfully initialized venues")

	mainMenuScreen := cli.NewMainMenu()

	historyViewScreen := event.NewViewerScreen("Concert History")
	historyDeleteScreen := event.NewDeleteScreen()
	historyAddScreen := event.NewAddScreen(false)
	artistEditScreen := artist.NewEditScreen()
	venueEditScreen := venue.NewEditScreen()
	openerRemoveScreen := event.NewOpenerRemovalScreen()

	historyViewScreen.Events = &pastEvents
	historyViewScreen.MainMenu = mainMenuScreen
	historyViewScreen.AddEventScreen = historyAddScreen
	historyViewScreen.DeleteEventScreen = historyDeleteScreen

	historyAddScreen.Events = &pastEvents
	historyAddScreen.EventAdder = interactor
	historyAddScreen.ArtistEditor = artistEditScreen
	historyAddScreen.VenueEditor = venueEditScreen
	historyAddScreen.OpenerRemover = openerRemoveScreen
	historyAddScreen.Viewer = historyViewScreen

	historyDeleteScreen.Events = &pastEvents
	historyDeleteScreen.Viewer = historyViewScreen
	historyDeleteScreen.Deleter = interactor

	futureViewScreen := event.NewViewerScreen("Future Concerts")
	futureDeleteScreen := event.NewDeleteScreen()
	futureAddScreen := event.NewAddScreen(true)

	futureViewScreen.Events = &futureEvents
	futureViewScreen.MainMenu = mainMenuScreen
	futureViewScreen.AddEventScreen = futureAddScreen
	futureViewScreen.DeleteEventScreen = futureDeleteScreen

	futureAddScreen.Events = &futureEvents
	futureAddScreen.EventAdder = interactor
	futureAddScreen.ArtistEditor = artistEditScreen
	futureAddScreen.VenueEditor = venueEditScreen
	futureAddScreen.OpenerRemover = openerRemoveScreen
	futureAddScreen.Viewer = futureViewScreen

	futureDeleteScreen.Events = &futureEvents
	futureDeleteScreen.Viewer = futureViewScreen
	futureDeleteScreen.Deleter = interactor

	artistEditScreen.Artists = artists
	artistEditScreen.ArtistAdder = interactor
	artistEditScreen.AddEventScreen = historyAddScreen

	venueEditScreen.Venues = venues
	venueEditScreen.VenueAdder = interactor
	venueEditScreen.AddEventScreen = historyAddScreen

	openerRemoveScreen.AddEventScreen = historyAddScreen

	mainMenuScreen.Children[1] = historyViewScreen
	mainMenuScreen.Children[2] = futureViewScreen
	out.Infoln("Successfully initialized CLI, starting display...")
	cli.RunCLI(mainMenuScreen)
}
