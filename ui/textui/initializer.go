package textui

import (
	"concert-manager/cache"
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/ui/textui/screens"
	"concert-manager/ui/textui/screens/artist"
	"concert-manager/ui/textui/screens/common"
	"concert-manager/ui/textui/screens/event"
	"concert-manager/ui/textui/screens/finder"
	"concert-manager/ui/textui/screens/venue"
	"concert-manager/util/format"
	"concert-manager/util/search"
)

func Start(cache *cache.LocalCache) {
	search := &search.Search{Cache: cache}

	addScreen := event.NewAddScreen()
	artistEditScreen := artist.NewEditScreen()

	artistSelectScreen := &common.Selector[data.Artist, data.Artist]{
		ScreenTitle: "Select Artist",
		Formatter: format.FormatArtist,
		OutputTransformer: common.IdentityTransform[data.Artist, data.Artist],
		Next: addScreen,
	}
	artistEditScreen.Search = search
	artistEditScreen.SelectScreen = artistSelectScreen

	venueEditScreen := venue.NewEditScreen()
	venueSelectScreen := &common.Selector[data.Venue, data.Venue]{
		ScreenTitle: "Select Venue",
		Formatter: format.FormatVenue,
		OutputTransformer: common.IdentityTransform[data.Venue, data.Venue],
		Next: addScreen,
	}
	venueEditScreen.Search = search
	venueEditScreen.SelectScreen = venueSelectScreen

	openerRemoveScreen := artist.NewOpenerRemoveScreen()

	addScreen.ArtistEditor = artistEditScreen
	addScreen.VenueEditor = venueEditScreen
	addScreen.OpenerRemover = openerRemoveScreen
	addScreen.Cache = cache

	deleteScreen := &event.Deleter{Cache: cache}

	historyViewScreen := event.NewViewScreen()
	historyViewScreen.ScreenTitle = "Concert History"
	historyViewScreen.AddEventScreen = addScreen
	historyViewScreen.DeleteEventScreen = deleteScreen
	historyViewScreen.ViewType = data.Past
	historyViewScreen.Cache = cache

	futureViewScreen := event.NewViewScreen()
	futureViewScreen.ScreenTitle = "Future Concerts"
	futureViewScreen.AddEventScreen = addScreen
	futureViewScreen.DeleteEventScreen = deleteScreen
	futureViewScreen.ViewType = data.Future
	futureViewScreen.Cache = cache

	eventDetailSelectScreen := &common.Selector[data.EventDetails, data.Event]{
		ScreenTitle: "Select Concert",
		Formatter: format.FormatEventDetailsShort,
		OutputTransformer: func(detail data.EventDetails) data.Event { return detail.Event },
		Next: addScreen,
	}
	finderViewScreen := finder.NewViewScreen()
	finderViewScreen.ScreenTitle = "All Upcoming Events"
	finderViewScreen.City = "Atlanta"
	finderViewScreen.State = "GA"
	finderViewScreen.Cache = cache
	finderViewScreen.AddEventSelectScreen = eventDetailSelectScreen

	finderMenuScreen := finder.NewMenu()
	finderMenuScreen.UpcomingEventViewer = finderViewScreen

	mainMenuScreen := screens.NewMainMenu()
	mainMenuScreen.Children[1] = historyViewScreen
	mainMenuScreen.Children[2] = futureViewScreen
	mainMenuScreen.Children[3] = finderMenuScreen

	log.Info("Successfully initialized terminal UI, starting display...")
	run(mainMenuScreen)
}
