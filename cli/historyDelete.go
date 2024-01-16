package cli

import (
	"concert-manager/data"
	"context"
	"fmt"
	"slices"
)

type EventDeleter interface {
    Delete(context.Context, data.Event) error
}

type HistoryDelete struct {
	Events   *[]data.Event
	DeleteSvc EventDeleter
	ParentScreen Screen
	startIdx int
	displayCount int
}

func (h *HistoryDelete) AddContext(startIdx int, displayCount int) {
    h.startIdx = startIdx
	h.displayCount = displayCount
}

func (h HistoryDelete) Title() string {
    return "Delete Concert"
}

func (h HistoryDelete) Data() string {
	if len(*h.Events) == 0 {
		return "No concerts found"
	}
	return ""
}

func (h HistoryDelete) Actions() []string {
	actions := []string{}
	for i := h.startIdx; i < h.startIdx + h.displayCount; i++ {
		event := (*h.Events)[i]
		var artist string
		if event.MainAct.Populated() {
			artist = event.MainAct.Name
		} else {
			artist = event.Openers[0].Name
		}

		eventDesc := fmt.Sprintf("%s; %v @ %s", artist, event.Date, event.Venue.Name)
		actions = append(actions, eventDesc)
	}
    return actions
}

func (h HistoryDelete) NextScreen(i int) Screen {
	eventIdx := h.startIdx + i - 1
	h.DeleteSvc.Delete(context.Background(), (*h.Events)[eventIdx])
	slices.Delete(*h.Events, eventIdx, eventIdx + 1)
    return h.ParentScreen
}

func (h HistoryDelete) Parent() Screen {
    return h.ParentScreen
}
