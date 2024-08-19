package com.bldover.beacon.data.api

import com.bldover.beacon.data.model.Venue
import retrofit2.http.GET

interface VenueApi {

    @GET("v1/venues")
    suspend fun getVenues(): List<Venue>
}