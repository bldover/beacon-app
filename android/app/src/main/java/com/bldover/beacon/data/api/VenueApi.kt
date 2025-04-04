package com.bldover.beacon.data.api

import com.bldover.beacon.data.model.venue.Venue
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path

interface VenueApi {

    @GET("v1/venues")
    suspend fun getVenues(): List<Venue>

    @PUT("v1/venues/{id}")
    suspend fun updateVenue(@Path("id") id: String, @Body venue: Venue)

    @POST("v1/venues")
    suspend fun addVenue(@Body venue: Venue): Venue

    @DELETE("v1/venues/{id}")
    suspend fun deleteVenue(@Path("id") id: String)
}