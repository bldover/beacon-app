package textui

import (
	"concert-manager/cache"
	"concert-manager/log"
	"concert-manager/ui/textui/core"
	"concert-manager/ui/textui/screens"
	"concert-manager/util"
)

func Start(cache *cache.LocalCache) {
	search := util.NewSearch()
	search.Cache = cache

	addScreen := screens.NewEventAddScreen()

	artistEditScreen := screens.NewArtistEditScreen()
	artistEditScreen.Search = search
	artistEditScreen.ReturnScreen = addScreen

	venueEditScreen := screens.NewVenueEditScreen()
	venueEditScreen.Search = search
	venueEditScreen.ReturnScreen = addScreen

	addScreen.ArtistEditor = artistEditScreen
	addScreen.VenueEditor = venueEditScreen
	addScreen.Cache = cache

	savedEventSearchResultScreen := screens.NewEventSearchResultScreen()
	savedEventSearchResultScreen.AddEventScreen = addScreen
	savedEventSearchResultScreen.Cache = cache

	savedEventViewScreen := screens.NewSavedEventViewScreen()
	savedEventViewScreen.AddEventScreen = addScreen
	savedEventViewScreen.Search = search
	savedEventViewScreen.SearchResultScreen = savedEventSearchResultScreen
	savedEventViewScreen.Cache = cache

	discoverySearchResultScreen := screens.NewDiscoverySearchResultScreen()
	discoverySearchResultScreen.AddEventScreen = addScreen

	discoveryViewScreen := screens.NewDiscoveryViewScreen()
	discoveryViewScreen.City = "Atlanta"
	discoveryViewScreen.State = "GA"
	discoveryViewScreen.AddEventScreen = addScreen
	discoveryViewScreen.Search = search
	discoveryViewScreen.SearchResultScreen = discoverySearchResultScreen
	discoveryViewScreen.Cache = cache

	discoveryMenuScreen := screens.NewDiscoveryMenu()
	discoveryMenuScreen.DiscoveryViewScreen = discoveryViewScreen

	passedEventsScreen := screens.NewPassedEventManager()
	passedEventsScreen.Cache = cache
	passedEventsScreen.AddEventScreen = addScreen

	utilityMenuScreen := screens.NewUtilMenu()
	utilityMenuScreen.PassedEventManager = passedEventsScreen

	mainMenuScreen := screens.NewMainMenu()
	mainMenuScreen.Children[1] = savedEventViewScreen
	mainMenuScreen.Children[2] = discoveryMenuScreen
	mainMenuScreen.Children[3] = utilityMenuScreen

	log.Info("Successfully initialized terminal UI, starting display...")
	core.Run(mainMenuScreen)
}
