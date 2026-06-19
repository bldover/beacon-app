package com.bldover.beacon.data.api

import com.bldover.beacon.data.dto.CountDto
import com.bldover.beacon.data.dto.EventsResponseDto
import com.bldover.beacon.data.dto.SummaryDto
import retrofit2.http.GET
import retrofit2.http.Path

interface AnalyticsApi {

    @GET("v1/analytics/summary")
    suspend fun getSummary(): SummaryDto

    @GET("v1/analytics/years")
    suspend fun getYears(): List<CountDto>

    @GET("v1/analytics/years/{key}")
    suspend fun getEventsByYear(@Path("key") key: String): EventsResponseDto

    @GET("v1/analytics/months")
    suspend fun getMonths(): List<CountDto>

    @GET("v1/analytics/months/{key}")
    suspend fun getEventsByMonth(@Path("key") key: String): EventsResponseDto

    @GET("v1/analytics/artists")
    suspend fun getArtists(): List<CountDto>

    @GET("v1/analytics/artists/{key}")
    suspend fun getEventsByArtist(@Path("key") key: String): EventsResponseDto

    @GET("v1/analytics/venues")
    suspend fun getVenues(): List<CountDto>

    @GET("v1/analytics/venues/{key}")
    suspend fun getEventsByVenue(@Path("key") key: String): EventsResponseDto

    @GET("v1/analytics/genres")
    suspend fun getGenres(): List<CountDto>

    @GET("v1/analytics/genres/{key}")
    suspend fun getEventsByGenre(@Path("key") key: String): EventsResponseDto
}
