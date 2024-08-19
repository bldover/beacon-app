package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.VenueApi
import com.bldover.beacon.data.model.Venue

interface VenueRepository {
    suspend fun getVenues(): List<Venue>
}

class VenueRepositoryImpl(private val venueApi: VenueApi) : VenueRepository {

    override suspend fun getVenues(): List<Venue> {
        return venueApi.getVenues()
    }
}