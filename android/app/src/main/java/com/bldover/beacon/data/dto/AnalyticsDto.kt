package com.bldover.beacon.data.dto

data class CountDto(
    val key: String,
    val name: String,
    val count: Int
)

data class SummaryDto(
    val totalEvents: Int,
    val topYears: List<CountDto>,
    val topMonths: List<CountDto>,
    val topArtists: List<CountDto>,
    val topVenues: List<CountDto>,
    val topGenres: List<CountDto>
)

data class EventsResponseDto(
    val count: Int,
    val events: List<EventDto>
)
