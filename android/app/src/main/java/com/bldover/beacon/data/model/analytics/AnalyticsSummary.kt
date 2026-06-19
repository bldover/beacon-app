package com.bldover.beacon.data.model.analytics

import com.bldover.beacon.data.dto.SummaryDto

data class AnalyticsSummary(
    val totalEvents: Int,
    val topYears: List<AnalyticsCount>,
    val topMonths: List<AnalyticsCount>,
    val topArtists: List<AnalyticsCount>,
    val topVenues: List<AnalyticsCount>,
    val topGenres: List<AnalyticsCount>
) {
    constructor(dto: SummaryDto) : this(
        totalEvents = dto.totalEvents,
        topYears = dto.topYears.map { AnalyticsCount(it) },
        topMonths = dto.topMonths.map { AnalyticsCount(it) },
        topArtists = dto.topArtists.map { AnalyticsCount(it) },
        topVenues = dto.topVenues.map { AnalyticsCount(it) },
        topGenres = dto.topGenres.map { AnalyticsCount(it) }
    )
}
