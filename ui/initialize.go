package ui

import (
	"concert-manager/cache"
	"concert-manager/ui/textui"
)

func StartUI(args string, savedCache *cache.SavedEventCache, upcomingCache *cache.UpcomingEventCache) {
	// parse args to determine which UI(s) to start
	textui.Start(savedCache, upcomingCache)
}
