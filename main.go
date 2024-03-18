package main

import (
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
	"context"
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

	interactor := &db.DatabaseRepository{}
	interactor.VenueRepo = venueRepo
	interactor.ArtistRepo = artistRepo
	interactor.EventRepo = eventRepo

	loader := &loader.Loader{EventCreator: interactor}
	go server.StartServer(loader)

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

	artists, err := interactor.ListArtists(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize artists:", err)
	}
	log.Info("Successfully initialized artists")

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

	// TODO: maybe migrate some of these to use context for some fields to
	// allow better reuse of the structs
	historyViewScreen := event.NewViewScreen("Concert History")
	historyDeleteScreen := event.NewDeleteScreen()
	historyAddScreen := event.NewAddScreen(false)
	historyArtistEditScreen := artist.NewEditScreen()
	historyVenueEditScreen := venue.NewEditScreen()
	historyOpenerRemoveScreen := artist.NewOpenerRemoveScreen()

	historyViewScreen.Events = pastEvents
	historyViewScreen.MainMenu = mainMenuScreen
	historyViewScreen.AddEventScreen = historyAddScreen
	historyViewScreen.DeleteEventScreen = historyDeleteScreen

	historyAddScreen.Events = pastEvents
	historyAddScreen.Database = dbRepo
	historyAddScreen.ArtistEditor = historyArtistEditScreen
	historyAddScreen.VenueEditor = historyVenueEditScreen
	historyAddScreen.OpenerRemover = historyOpenerRemoveScreen
	historyAddScreen.Viewer = historyViewScreen

	historyDeleteScreen.Events = pastEvents
	historyDeleteScreen.Viewer = historyViewScreen
	historyDeleteScreen.Database = dbRepo

	historyArtistEditScreen.Artists = artists
	historyArtistEditScreen.AddEventScreen = historyAddScreen

	historyVenueEditScreen.Venues = venues
	historyVenueEditScreen.AddEventScreen = historyAddScreen

	historyOpenerRemoveScreen.AddEventScreen = historyAddScreen

	futureViewScreen := event.NewViewScreen("Future Concerts")
	futureDeleteScreen := event.NewDeleteScreen()
	futureAddScreen := event.NewAddScreen(true)
	futureArtistEditScreen := artist.NewEditScreen()
	futureVenueEditScreen := venue.NewEditScreen()
	futureOpenerRemoveScreen := artist.NewOpenerRemoveScreen()

	futureViewScreen.Events = futureEvents
	futureViewScreen.MainMenu = mainMenuScreen
	futureViewScreen.AddEventScreen = futureAddScreen
	futureViewScreen.DeleteEventScreen = futureDeleteScreen

	futureAddScreen.Events = futureEvents
	futureAddScreen.Database = dbRepo
	futureAddScreen.ArtistEditor = futureArtistEditScreen
	futureAddScreen.VenueEditor = futureVenueEditScreen
	futureAddScreen.OpenerRemover = futureOpenerRemoveScreen
	futureAddScreen.Viewer = futureViewScreen

	futureDeleteScreen.Events = futureEvents
	futureDeleteScreen.Viewer = futureViewScreen
	futureDeleteScreen.Database = dbRepo

	futureArtistEditScreen.Artists = artists
	futureArtistEditScreen.AddEventScreen = futureAddScreen

	futureVenueEditScreen.Venues = venues
	futureVenueEditScreen.AddEventScreen = futureAddScreen

	futureOpenerRemoveScreen.AddEventScreen = futureAddScreen

	finderMenuScreen := finder.NewMenu()
	finderViewScreen := finder.NewViewerScreen("All Upcoming Events", "Atlanta", "GA")
	finderSelectorScreen := finder.NewSelectorScreen()
	finderAddScreen := futureAddScreen
	finderArtistEditScreen := artist.NewEditScreen()
	finderVenueEditScreen := venue.NewEditScreen()
	finderOpenerRemoveScreen := artist.NewOpenerRemoveScreen()

	finderMenuScreen.MainMenu = mainMenuScreen
	finderMenuScreen.UpcomingEventViewer = finderViewScreen

	finderViewScreen.AddEventSelectorScreen = finderSelectorScreen
	finderViewScreen.FinderMenu = finderMenuScreen
	finderViewScreen.EventFinder = finderSvc.NewEventFinder()

	finderAddScreen.Events = futureEvents
	finderAddScreen.Database = dbRepo
	finderAddScreen.ArtistEditor = finderArtistEditScreen
	finderAddScreen.VenueEditor = finderVenueEditScreen
	finderAddScreen.OpenerRemover = finderOpenerRemoveScreen
	finderAddScreen.Viewer = finderViewScreen

	finderArtistEditScreen.Artists = artists
	finderArtistEditScreen.AddEventScreen = finderAddScreen

	finderVenueEditScreen.Venues = venues
	finderVenueEditScreen.AddEventScreen = finderAddScreen

	finderOpenerRemoveScreen.AddEventScreen = finderAddScreen

	finderSelectorScreen.EventAddScreen = finderAddScreen
	finderSelectorScreen.ViewerScreen = finderViewScreen

	mainMenuScreen.Children[1] = historyViewScreen
	mainMenuScreen.Children[2] = futureViewScreen
	mainMenuScreen.Children[3] = finderMenuScreen

	return mainMenuScreen
}
