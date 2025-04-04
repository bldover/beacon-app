package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.VenueApi
import com.bldover.beacon.data.model.venue.Venue

interface VenueRepository {
    suspend fun getVenues(): List<Venue>
    suspend fun addVenue(venue: Venue): Venue
    suspend fun updateVenue(venue: Venue): Venue
    suspend fun deleteVenue(venue: Venue)
}

class VenueRepositoryImpl(private val venueApi: VenueApi) : VenueRepository {

    override suspend fun getVenues(): List<Venue> {
        return venueApi.getVenues()
    }

    override suspend fun addVenue(venue: Venue): Venue {
        return venueApi.addVenue(venue)
    }

    override suspend fun updateVenue(venue: Venue): Venue {
        venueApi.updateVenue(venue.id!!, venue)
        return venue
    }

    override suspend fun deleteVenue(venue: Venue) {
        venueApi.deleteVenue(venue.id!!)
    }
}