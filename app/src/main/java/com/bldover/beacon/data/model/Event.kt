package com.bldover.beacon.data.model

import java.time.LocalDate
import java.time.format.DateTimeFormatter

private val dateFormatter: DateTimeFormatter = DateTimeFormatter.ofPattern("M/d/yyyy")

data class Event(
    val id: String,
    val artists: List<Artist>,
    val venue: Venue,
    val date: LocalDate,
    val purchased: Boolean
) {
    constructor(event: RawEvent) : this(
        id = event.id,
        artists = event.artists,
        venue = event.venue,
        date = LocalDate.parse(event.date, dateFormatter),
        purchased = event.purchased
    )

    fun hasMatch(term: String): Boolean = artists.any { it.name.contains(term, ignoreCase = true) }
}