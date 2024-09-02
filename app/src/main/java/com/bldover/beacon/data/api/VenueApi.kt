package com.bldover.beacon.data.api

import com.bldover.beacon.data.model.Venue
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Query

interface VenueApi {

    @GET("v1/venues")
    suspend fun getVenues(): List<Venue>
    @POST("v1/venues")
    suspend fun addVenue(@Body venue: Venue): Venue

    @DELETE("v1/venues")
    suspend fun deleteVenue(@Query("id") id: String)
}