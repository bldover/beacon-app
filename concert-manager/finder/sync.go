package finder

import (
	"concert-manager/domain"
	"errors"
	"fmt"
	"slices"
)

// could be optimized but this is simple and fine for the data size
func (c *Cache) enrichSavedData(event domain.EventDetails) domain.EventDetails {
	savedArtists := c.SavedDataCache.GetArtists()
	savedVenues := c.SavedDataCache.GetVenues()
	savedEvents := c.SavedDataCache.GetSavedEvents()

	enriched := domain.CloneEventDetail(event)

	eventIdx := slices.IndexFunc(savedEvents, func(o domain.Event) bool {
		return eventMatch(event.Event, o)
	})
	if eventIdx != -1 {
		enriched.Event = mergeEvent(savedEvents[eventIdx], event.Event)
		return enriched
	}

	venueIdx := slices.IndexFunc(savedVenues, func(o domain.Venue) bool {
		return venueMatch(event.Event.Venue, o)
	})
	if venueIdx != -1 {
		enriched.Event.Venue = mergeVenue(savedVenues[venueIdx], event.Event.Venue)
	}

	if event.Event.MainAct != nil {
		mainActIdx := slices.IndexFunc(savedArtists, func(o domain.Artist) bool {
			return artistMatch(*event.Event.MainAct, o)
		})
		if mainActIdx != -1 {
			artist := mergeArtist(savedArtists[mainActIdx], *event.Event.MainAct)
			enriched.Event.MainAct = &artist
		}
	}

	for i, opener := range event.Event.Openers {
		openerIdx := slices.IndexFunc(savedArtists, func(o domain.Artist) bool {
			return artistMatch(opener, o)
		})
		if openerIdx != -1 {
			artist := mergeArtist(savedArtists[openerIdx], event.Event.Openers[i])
			enriched.Event.Openers[i] = artist
		}
	}

	return enriched
}

func eventMatch(a domain.Event, b domain.Event) bool {
	fieldsMatch := venueMatch(a.Venue, b.Venue) && a.Date == b.Date && artistMatch(*a.MainAct, *b.MainAct)
	return (a.ID.Ticketmaster != "" && a.ID.Ticketmaster == b.ID.Ticketmaster) || fieldsMatch
}

func venueMatch(a domain.Venue, b domain.Venue) bool {
	return (a.ID.Ticketmaster != "" && a.ID.Ticketmaster == b.ID.Ticketmaster) || a.EqualsFields(b)
}

func artistMatch(a domain.Artist, b domain.Artist) bool {
	tmMatch := a.ID.Ticketmaster != "" && a.ID.Ticketmaster == b.ID.Ticketmaster
	spotifyMatch := a.ID.Spotify != "" && a.ID.Spotify == b.ID.Spotify
	mbMatch := a.ID.MusicBrainz != "" && a.ID.MusicBrainz == b.ID.MusicBrainz
	return tmMatch || spotifyMatch || mbMatch || a.EqualsFields(b)
}

func mergeEvent(source domain.Event, target domain.Event) domain.Event {
	event := domain.CloneEvent(target)
	event.MainAct = source.MainAct
	event.Openers = source.Openers
	event.Venue = source.Venue
	event.Purchased = source.Purchased
	event.ID.Primary = source.ID.Primary
	if source.ID.Ticketmaster != "" {
		event.ID.Ticketmaster = source.ID.Ticketmaster
	}
	return event
}

func mergeArtist(source domain.Artist, target domain.Artist) domain.Artist {
	artist := domain.CloneArtist(target)
	artist.Name = source.Name
	artist.Genres = source.Genres
	artist.ID.Primary = source.ID.Primary
	if source.ID.Ticketmaster != "" {
		artist.ID.Ticketmaster = source.ID.Ticketmaster
	}
	if source.ID.Spotify != "" {
		artist.ID.Spotify = source.ID.Spotify
	}
	if source.ID.MusicBrainz != "" {
		artist.ID.MusicBrainz = source.ID.MusicBrainz
	}
	return artist
}

func mergeVenue(source domain.Venue, target domain.Venue) domain.Venue {
	venue := domain.CloneVenue(target)
	venue.Name = source.Name
	venue.City = source.City
	venue.State = source.State
	venue.ID.Primary = source.ID.Primary
	if source.ID.Ticketmaster != "" {
		venue.ID.Ticketmaster = source.ID.Ticketmaster
	}
	return venue
}

func (c *Cache) SyncArtistAdd(id string) error {
	savedArtists := c.SavedDataCache.GetArtists()
	newArtistIdx := slices.IndexFunc(savedArtists, func(o domain.Artist) bool { return id == o.ID.Primary })
	if newArtistIdx < 0 {
		errMsg := fmt.Sprintf("unable to find new cached artist with id %s in SyncArtistAdd", id)
		return errors.New(errMsg)
	}
	newArtist := savedArtists[newArtistIdx]

	for _, eventData := range c.upcomingEvents {
		events := eventData.Events
		for i, event := range events {
			mainAct := event.Event.MainAct
			if mainAct != nil && mainAct.EqualsFields(newArtist) {
				clonedArtist := domain.CloneArtist(newArtist)
				events[i].Event.MainAct = &clonedArtist
			}

			for j, opener := range event.Event.Openers {
				if opener.EqualsFields(newArtist) {
					clonedArtist := domain.CloneArtist(newArtist)
					events[i].Event.Openers[j] = clonedArtist
				}
			}
		}
	}
	return nil
}

func (c *Cache) SyncArtistUpdate(id string) error {
	savedArtists := c.SavedDataCache.GetArtists()
	updatedArtistIdx := slices.IndexFunc(savedArtists, func(o domain.Artist) bool { return id == o.ID.Primary })
	if updatedArtistIdx < 0 {
		errMsg := fmt.Sprintf("unable to find updated cached artist with id %s in SyncArtistUpdate", id)
		return errors.New(errMsg)
	}
	updatedArtist := savedArtists[updatedArtistIdx]

	for _, eventData := range c.upcomingEvents {
		events := eventData.Events
		for i, event := range events {
			mainAct := event.Event.MainAct
			if mainAct != nil && (mainAct.EqualsFields(updatedArtist) || mainAct.Equals(updatedArtist)) {
				clonedArtist := domain.CloneArtist(updatedArtist)
				events[i].Event.MainAct = &clonedArtist
			}

			for j, opener := range event.Event.Openers {
				if opener.EqualsFields(updatedArtist) || opener.Equals(updatedArtist) {
					clonedArtist := domain.CloneArtist(updatedArtist)
					events[i].Event.Openers[j] = clonedArtist
				}
			}
		}
	}
	return nil
}

func (c *Cache) SyncArtistDelete(id string) {
	for _, eventData := range c.upcomingEvents {
		events := eventData.Events
		for i, event := range events {
			mainAct := event.Event.MainAct
			if mainAct != nil && mainAct.ID.Primary == id {
				events[i].Event.MainAct.ID.Primary = ""
			}

			for j, opener := range event.Event.Openers {
				if opener.ID.Primary == id {
					events[i].Event.Openers[j].ID.Primary = ""
				}
			}
		}
	}
}

func (c *Cache) SyncVenueAdd(id string) error {
	savedVenues := c.SavedDataCache.GetVenues()
	newVenueIdx := slices.IndexFunc(savedVenues, func(o domain.Venue) bool { return id == o.ID.Primary })
	if newVenueIdx < 0 {
		errMsg := fmt.Sprintf("unable to find new cached venue with id %s in SyncVenueAdd", id)
		return errors.New(errMsg)
	}
	newVenue := savedVenues[newVenueIdx]

	for _, eventData := range c.upcomingEvents {
		events := eventData.Events
		for i, event := range events {
			if event.Event.Venue.Equals(newVenue) {
				events[i].Event.Venue = newVenue
			}
		}
	}
	return nil
}

func (c *Cache) SyncVenueUpdate(id string) error {
	savedVenues := c.SavedDataCache.GetVenues()
	updatedVenueIdx := slices.IndexFunc(savedVenues, func(o domain.Venue) bool { return id == o.ID.Primary })
	if updatedVenueIdx < 0 {
		errMsg := fmt.Sprintf("unable to find updated cached venue with id %s in SyncVenueUpdate", id)
		return errors.New(errMsg)
	}
	updatedVenue := savedVenues[updatedVenueIdx]

	for _, eventData := range c.upcomingEvents {
		events := eventData.Events
		for i, event := range events {
			venue := event.Event.Venue
			if venue.EqualsFields(updatedVenue) || venue.Equals(updatedVenue) {
				events[i].Event.Venue = updatedVenue
			}
		}
	}
	return nil
}

func (c *Cache) SyncVenueDelete(id string) {
	for _, eventData := range c.upcomingEvents {
		events := eventData.Events
		for i, event := range events {
			if event.Event.Venue.ID.Primary == id {
				events[i].Event.Venue.ID.Primary = ""
			}
		}
	}
}

func (c *Cache) SyncEventAdd(id string) error {
	savedEvents := c.SavedDataCache.GetSavedEvents()
	newEventIdx := slices.IndexFunc(savedEvents, func(o domain.Event) bool { return id == o.ID.Primary })
	if newEventIdx < 0 {
		errMsg := fmt.Sprintf("unable to find new cached event with id %s in SyncEventAdd", id)
		return errors.New(errMsg)
	}
	newEvent := savedEvents[newEventIdx]

	for _, eventData := range c.upcomingEvents {
		events := eventData.Events
		for i, event := range events {
			if event.Event.EqualsFields(newEvent) {
				events[i].Event = newEvent
			}
		}
	}
	return nil
}

func (c *Cache) SyncEventUpdate(id string) error {
	savedEvents := c.SavedDataCache.GetSavedEvents()
	updatedEventIdx := slices.IndexFunc(savedEvents, func(o domain.Event) bool { return id == o.ID.Primary })
	if updatedEventIdx < 0 {
		errMsg := fmt.Sprintf("unable to find updated cached event with id %s in SyncEventUpdate", id)
		return errors.New(errMsg)
	}
	updatedEvent := savedEvents[updatedEventIdx]

	for _, eventData := range c.upcomingEvents {
		events := eventData.Events
		for i, event := range events {
			if event.Event.EqualsFields(updatedEvent) || event.Event.Equals(updatedEvent) {
				events[i].Event = updatedEvent
			}
		}
	}
	return nil
}

func (c *Cache) SyncEventDelete(id string) {
	for _, eventData := range c.upcomingEvents {
		events := eventData.Events
		for i, event := range events {
			if event.Event.ID.Primary == id {
				events[i].Event.ID.Primary = ""
			}
		}
	}
}
