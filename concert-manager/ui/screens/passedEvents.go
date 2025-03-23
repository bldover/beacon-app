package screens

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/ui/output"
	"concert-manager/util"
)

type passedEventCache interface {
	GetPassedSavedEvents() []data.Event
	AddSavedEvent(data.Event) (*data.Event, error)
	DeleteSavedEvent(string) error
}

type PassedEventManager struct {
	Cache          passedEventCache
	AddEventScreen *EventAdder
	actions        []string
	passedEvents   []data.Event
	currentEvent   data.Event
	loaded         bool
}

const (
	markAsAttended = iota + 1
	editPassedEvent
	removePassedEvent
	passedEventsToMenu
)

func NewPassedEventManager() *PassedEventManager {
	template := PassedEventManager{}
	template.actions = []string{"Mark As Attended", "Edit", "Remove Event", "Utility Menu"}
	return &template
}

func (m PassedEventManager) Title() string {
	return "Manage Passed Events"
}

func (m *PassedEventManager) DisplayData() {
	if !m.loaded {
		m.passedEvents = m.Cache.GetPassedSavedEvents()
		m.loaded = true
	}

	if len(m.passedEvents) == 0 {
		output.Displayln("No passed events")
		return
	}

	m.currentEvent = m.passedEvents[len(m.passedEvents)-1]
	output.Displayln(util.FormatEvent(m.currentEvent))
}

func (m PassedEventManager) Actions() []string {
	return m.actions
}

func (m *PassedEventManager) NextScreen(i int) Screen {
	switch i {
	case markAsAttended:
		if !m.currentEvent.Populated() {
			output.Displayln("No event to mark")
			return m
		}

		m.currentEvent.Purchased = true
		if err := m.Cache.DeleteSavedEvent(m.currentEvent.Id); err != nil {
			log.Error("Failed to delete passed event:", err)
			output.Displayln("Failed to update event")
		}
		if _, err := m.Cache.AddSavedEvent(m.currentEvent); err != nil {
			log.Error("Failed to add passed event after delete:", err)
			output.Displayln("Failed to update event")
		}

		m.passedEvents = m.passedEvents[:len(m.passedEvents)-1]
	case editPassedEvent:
		if !m.currentEvent.Populated() {
			output.Displayln("No event to edit")
			return m
		}

		m.AddEventScreen.WithBeforeSaveAction(func() error {
			if err := m.Cache.DeleteSavedEvent(m.currentEvent.Id); err != nil {
				return err
			}
			m.passedEvents = m.passedEvents[:len(m.passedEvents)-1]
			m.AddEventScreen.newEvent.Purchased = true
			return nil
		})
		m.currentEvent.Purchased = true
		m.AddEventScreen.newEvent = m.currentEvent
		return m.AddEventScreen
	case removePassedEvent:
		if !m.currentEvent.Populated() {
			output.Displayln("No event to remove")
			return m
		}

		if err := m.Cache.DeleteSavedEvent(m.currentEvent.Id); err != nil {
			log.Error("Failed to delete passed event:", err)
			output.Displayln("Failed to update event")
		}

		m.passedEvents = m.passedEvents[:len(m.passedEvents)-1]
	case passedEventsToMenu:
		m.loaded = false
		return nil
	}
	return m
}
