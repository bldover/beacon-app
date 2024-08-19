package com.bldover.beacon.data.api

import com.bldover.beacon.data.model.RawEvent
import com.bldover.beacon.data.model.RawEventDetail
import com.bldover.beacon.data.model.Recommendation
import retrofit2.http.GET
import retrofit2.http.Query

interface EventApi {

    @GET("v1/events/saved")
    suspend fun getSavedEvents(): List<RawEvent>

    @GET("v1/events/upcoming")
    suspend fun getUpcomingEvents(): List<RawEventDetail>

    @GET("v1/events/recommended")
    suspend fun getRecommendations(@Query("threshold") threshold: Boolean): List<Recommendation>
}