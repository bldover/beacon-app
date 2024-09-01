package com.bldover.beacon.data.model

import com.bldover.beacon.data.util.dateFormatter
import java.time.LocalDate

data class Event(
    var id: String? = null,
    var artists: List<Artist>,
    var venue: Venue,
    var date: LocalDate,
    var purchased: Boolean
) {
    constructor(event: RawEvent) : this(
        id = event.id,
        artists = event.artists,
        venue = event.venue,
        date = LocalDate.parse(event.date, dateFormatter),
        purchased = event.purchased
    )

    constructor(eventDetail: EventDetail) : this(
        id = eventDetail.id,
        artists = eventDetail.artists,
        venue = eventDetail.venue,
        date = eventDetail.date,
        purchased = eventDetail.purchased
    )

    fun hasMatch(term: String): Boolean = artists.any { it.name.contains(term, ignoreCase = true) }

    fun isPopulated(): Boolean {
        return artists.isNotEmpty() && artists.all { it.isPopulated() } && venue.isPopulated()
    }
}