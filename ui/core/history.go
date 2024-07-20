package core

import (
	"concert-manager/log"
	"concert-manager/ui/screens"
	"slices"
	"strings"
)

type screenHistory []screens.Screen

func (s screenHistory) String() string {
	screenNames := []string{}
	for _, screen := range s {
		screenNames = append(screenNames, screen.Title())
	}
	return strings.Join(screenNames, "->")
}

type history struct {
	history screenHistory
}

func (h *history) update(screen screens.Screen) {
	if i := slices.Index(h.history, screen); i != -1 {
		h.history = h.history[:i+1]
	} else {
		h.history = append(h.history, screen)
	}
	log.Debugf("History updated: %v", h.history)
}

func (h *history) getPrevious() screens.Screen {
	h.history = h.history[:len(h.history)-1]
	screen := h.history[len(h.history)-1]
	// hack, couples this with select title implementation
	for strings.Contains(screen.Title(), "Select") {
		h.history = h.history[:len(h.history)-1]
		screen = h.history[len(h.history)-1]
	}
	log.Debugf("History previous: %s, new history: %v", screen.Title(), h.history)
	return screen
}
