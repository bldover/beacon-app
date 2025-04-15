package ui

import (
	"concert-manager/db"
	"concert-manager/finder"
	"concert-manager/log"
	"concert-manager/ui/core"
	"concert-manager/ui/screens"
)

func Start(savedCache *db.Cache, upcomingCache *finder.Cache) {
	log.Info("Initializing terminal UI")

	addScreen := screens.NewEventAddScreen()

	artistEditScreen := screens.NewArtistEditScreen()
	artistEditScreen.ArtistCache = savedCache
	artistEditScreen.ReturnScreen = addScreen

	venueEditScreen := screens.NewVenueEditScreen()
	venueEditScreen.VenueCache = savedCache
	venueEditScreen.ReturnScreen = addScreen

	addScreen.ArtistEditor = artistEditScreen
	addScreen.VenueEditor = venueEditScreen
	addScreen.Cache = savedCache

	savedEventSearchResultScreen := screens.NewEventSearchResultScreen()
	savedEventSearchResultScreen.AddEventScreen = addScreen
	savedEventSearchResultScreen.Cache = savedCache

	savedEventViewScreen := screens.NewSavedEventViewScreen()
	savedEventViewScreen.AddEventScreen = addScreen
	savedEventViewScreen.SearchResultScreen = savedEventSearchResultScreen
	savedEventViewScreen.Cache = savedCache

	discoverySearchResultScreen := screens.NewDiscoverySearchResultScreen()
	discoverySearchResultScreen.AddEventScreen = addScreen

	discoveryViewScreen := screens.NewDiscoveryViewScreen()
	discoveryViewScreen.AddEventScreen = addScreen
	discoveryViewScreen.SearchResultScreen = discoverySearchResultScreen
	discoveryViewScreen.Cache = upcomingCache

	recommendedViewScreen := screens.NewRecommendationScreen()
	recommendedViewScreen.AddEventScreen = addScreen
	recommendedViewScreen.RecommendationCache = upcomingCache
	recommendedViewScreen.SavedCache = savedCache

	discoveryMenuScreen := screens.NewDiscoveryMenu()
	discoveryMenuScreen.DiscoveryViewScreen = discoveryViewScreen
	discoveryMenuScreen.RecommendationViewScreen = recommendedViewScreen

	passedEventsScreen := screens.NewPassedEventManager()
	passedEventsScreen.Cache = savedCache
	passedEventsScreen.AddEventScreen = addScreen

	utilityMenuScreen := screens.NewUtilMenu()
	utilityMenuScreen.PassedEventManager = passedEventsScreen

	mainMenuScreen := screens.NewMainMenu()
	mainMenuScreen.Children[1] = savedEventViewScreen
	mainMenuScreen.Children[2] = discoveryMenuScreen
	mainMenuScreen.Children[3] = utilityMenuScreen

	log.Info("Successfully initialized terminal UI")
	core.Run(mainMenuScreen)
}
