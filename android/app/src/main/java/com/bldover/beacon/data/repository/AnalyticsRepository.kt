package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.AnalyticsApi
import com.bldover.beacon.data.model.analytics.AnalyticsCategory
import com.bldover.beacon.data.model.analytics.AnalyticsCount
import com.bldover.beacon.data.model.analytics.AnalyticsSummary
import com.bldover.beacon.data.model.event.Event

interface AnalyticsRepository {
    suspend fun getSummary(): AnalyticsSummary
    suspend fun getCounts(category: AnalyticsCategory): List<AnalyticsCount>
    suspend fun getEvents(category: AnalyticsCategory, key: String): List<Event>
}

class AnalyticsRepositoryImpl(private val analyticsApi: AnalyticsApi) : AnalyticsRepository {

    override suspend fun getSummary(): AnalyticsSummary {
        return AnalyticsSummary(analyticsApi.getSummary())
    }

    override suspend fun getCounts(category: AnalyticsCategory): List<AnalyticsCount> {
        val dtos = when (category) {
            AnalyticsCategory.YEARS -> analyticsApi.getYears()
            AnalyticsCategory.MONTHS -> analyticsApi.getMonths()
            AnalyticsCategory.ARTISTS -> analyticsApi.getArtists()
            AnalyticsCategory.VENUES -> analyticsApi.getVenues()
            AnalyticsCategory.GENRES -> analyticsApi.getGenres()
        }
        return dtos.map { AnalyticsCount(it) }
    }

    override suspend fun getEvents(category: AnalyticsCategory, key: String): List<Event> {
        val response = when (category) {
            AnalyticsCategory.YEARS -> analyticsApi.getEventsByYear(key)
            AnalyticsCategory.MONTHS -> analyticsApi.getEventsByMonth(key)
            AnalyticsCategory.ARTISTS -> analyticsApi.getEventsByArtist(key)
            AnalyticsCategory.VENUES -> analyticsApi.getEventsByVenue(key)
            AnalyticsCategory.GENRES -> analyticsApi.getEventsByGenre(key)
        }
        return response.events.map { Event(it) }
    }
}
