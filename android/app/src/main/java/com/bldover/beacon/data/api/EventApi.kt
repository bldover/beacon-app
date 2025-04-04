package com.bldover.beacon.data.api

import com.bldover.beacon.data.model.event.RawEvent
import com.bldover.beacon.data.model.event.RawEventDetail
import com.bldover.beacon.data.model.event.RawEventRank
import com.bldover.beacon.data.model.RecommendationThreshold
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path
import retrofit2.http.Query

interface EventApi {

    @GET("v1/events/saved")
    suspend fun getSavedEvents(): List<RawEvent>

    @GET("v1/events/saved")
    suspend fun getEvent(@Query("id") id: String): List<RawEvent>

    @POST("v1/events/saved")
    suspend fun addEvent(@Body event: RawEvent)

    @DELETE("v1/events/saved/{id}")
    suspend fun deleteEvent(@Path("id") id: String)

    @GET("v1/events/upcoming")
    suspend fun getUpcomingEvents(): List<RawEventDetail>

    @GET("v1/events/recommended")
    suspend fun getRecommendations(@Query("threshold") threshold: RecommendationThreshold): List<RawEventRank>
}