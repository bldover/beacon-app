package finder

import (
	"concert-manager/domain"
	"concert-manager/log"
	"errors"
	"fmt"
	"slices"
)

// This logic helps with a couple things:
//  1. Substitute saved events/venues/artists into found events. This is especially
//     helpful for populating the primary ID so we know which data needs to be inserted
//     when the event is saved.
//  2. Updating IDs for improved future matching. Since the data source is Ticketmaster,
//     we can compare the TM IDs to see if our saved data is a match. However, our saved
//     data could have been manually created (no IDs) or did not have complete data from
//     TM when we first saved it, so we will enrich the saved data with any IDs from TM that
//     we don't already know about. This helps with metadata lookup and future matching here.
func (c *Cache) enrichSavedData(event domain.EventDetails) domain.EventDetails {
	savedArtists := c.SavedDataCache.GetArtists()
	savedVenues := c.SavedDataCache.GetVenues()
	savedEvents := c.SavedDataCache.GetSavedEvents()

	enriched := domain.CloneEventDetail(event)

	eventIdx := slices.IndexFunc(savedEvents, func(o domain.Event) bool {
		return eventMatch(event.Event, o)
	})
	if eventIdx != -1 {
		enriched.Event = c.mergeEvent(savedEvents[eventIdx], event.Event)
		return enriched
	}

	venueIdx := slices.IndexFunc(savedVenues, func(o domain.Venue) bool {
		return venueMatch(event.Event.Venue, o)
	})
	if venueIdx != -1 {
		enriched.Event.Venue = c.mergeVenue(savedVenues[venueIdx], event.Event.Venue)
	}

	if event.Event.MainAct != nil {
		mainActIdx := slices.IndexFunc(savedArtists, func(o domain.Artist) bool {
			return artistMatch(*event.Event.MainAct, o)
		})
		if mainActIdx != -1 {
			artist := c.mergeArtist(savedArtists[mainActIdx], *event.Event.MainAct)
			enriched.Event.MainAct = &artist
		}
	}

	for i, opener := range event.Event.Openers {
		openerIdx := slices.IndexFunc(savedArtists, func(o domain.Artist) bool {
			return artistMatch(opener, o)
		})
		if openerIdx != -1 {
			artist := c.mergeArtist(savedArtists[openerIdx], event.Event.Openers[i])
			enriched.Event.Openers[i] = artist
		}
	}

	return enriched
}

func eventMatch(ext domain.Event, saved domain.Event) bool {
	fieldsMatch := venueMatch(ext.Venue, saved.Venue) && ext.Date == saved.Date && artistMatch(*ext.MainAct, *saved.MainAct)
	return (saved.ID.Ticketmaster != "" && ext.ID.Ticketmaster == saved.ID.Ticketmaster) || (saved.ID.Ticketmaster == "" && fieldsMatch)
}

func venueMatch(ext domain.Venue, saved domain.Venue) bool {
	return (saved.ID.Ticketmaster != "" && ext.ID.Ticketmaster == saved.ID.Ticketmaster) || (saved.ID.Ticketmaster == "" && ext.EqualsFields(saved))
}

func artistMatch(ext domain.Artist, saved domain.Artist) bool {
	tmMatch := saved.ID.Ticketmaster != "" && ext.ID.Ticketmaster == saved.ID.Ticketmaster
	return tmMatch || (saved.ID.Ticketmaster == "" && ext.EqualsFields(saved))
}

func (c *Cache) mergeEvent(source domain.Event, target domain.Event) domain.Event {
	event := domain.CloneEvent(target)
	event.MainAct = source.MainAct
	event.Openers = source.Openers
	event.Venue = source.Venue
	event.Purchased = source.Purchased
	event.ID.Primary = source.ID.Primary

	// due to match logic, either the TM IDs match or the source didn't have an ID
	if source.ID.Ticketmaster == "" && target.ID.Ticketmaster != "" {
		if err := c.SavedDataCache.UpdateSavedEvent(event.ID.Primary, event); err != nil {
			log.Alertf("Failed to update event while merging upcoming results for source: %v, target: %v, err: %v", source, target, err)
			return event
		}
		if err := c.SyncEventUpdate(event.ID.Primary); err != nil {
			log.Alertf("Failed to sync event update while merging upcoming results for source: %v, target: %v, err: %v", source, target, err)
			return event
		}
	}

	return event
}

func (c *Cache) mergeArtist(source domain.Artist, target domain.Artist) domain.Artist {
	artist := domain.CloneArtist(target)
	artist.Name = source.Name
	artist.Genres = source.Genres
	artist.ID.Primary = source.ID.Primary

	if target.ID.Ticketmaster == "" {
		artist.ID.Ticketmaster = source.ID.Ticketmaster
	}
	if target.ID.Spotify == "" {
		artist.ID.Spotify = source.ID.Spotify
	}
	if target.ID.MusicBrainz == "" {
		artist.ID.MusicBrainz = source.ID.MusicBrainz
	}

	tmIDAdded := source.ID.Ticketmaster == "" && target.ID.Ticketmaster != ""
	spotifyIDAdded := source.ID.Spotify == "" && target.ID.Spotify != ""
	mbIDAdded := source.ID.MusicBrainz == "" && target.ID.MusicBrainz != ""
	if tmIDAdded || spotifyIDAdded || mbIDAdded {
		if err := c.SavedDataCache.UpdateArtist(artist.ID.Primary, artist); err != nil {
			log.Errorf("Failed to update artist while merging upcoming results for source: %v, target: %v", source, target)
			return artist
		}
		if err := c.SyncArtistUpdate(artist.ID.Primary); err != nil {
			log.Errorf("Failed to sync artist update while merging upcoming results for source: %v, target: %v", source, target)
			return artist
		}
	}

	return artist
}

func (c *Cache) mergeVenue(source domain.Venue, target domain.Venue) domain.Venue {
	venue := domain.CloneVenue(target)
	venue.Name = source.Name
	venue.City = source.City
	venue.State = source.State
	venue.ID.Primary = source.ID.Primary

	// due to match logic, either the TM IDs match or the source didn't have an ID
	if source.ID.Ticketmaster == "" && target.ID.Ticketmaster != "" {
		if err := c.SavedDataCache.UpdateVenue(venue.ID.Primary, venue); err != nil {
			log.Errorf("Failed to update venue while merging upcoming results for source: %v, target: %v", source, target)
			return venue
		}
		if err := c.SyncVenueUpdate(venue.ID.Primary); err != nil {
			log.Errorf("Failed to sync venue update while merging upcoming results for source: %v, target: %v", source, target)
			return venue
		}
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
