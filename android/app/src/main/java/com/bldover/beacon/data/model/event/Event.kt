package com.bldover.beacon.data.model.event

import com.bldover.beacon.data.model.venue.Venue
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.util.dateFormatter
import java.time.LocalDate

data class Event(
    var id: String? = null,
    var artists: List<Artist>,
    var venue: Venue,
    var date: LocalDate,
    var purchased: Boolean,
    val ticketmasterId: String? = null
) {
    constructor(event: RawEvent) : this(
        id = event.id,
        artists = event.artists,
        venue = Venue(event.venue),
        date = LocalDate.parse(event.date, dateFormatter),
        purchased = event.purchased,
        ticketmasterId = event.tmId.ifBlank { null }
    )

    constructor(eventDetail: EventDetail) : this(
        id = eventDetail.id,
        artists = eventDetail.artists,
        venue = eventDetail.venue,
        date = eventDetail.date,
        purchased = eventDetail.purchased,
        ticketmasterId = eventDetail.ticketmasterId
    )

    fun hasMatch(term: String): Boolean {
        return artists.any { it.name.contains(term, ignoreCase = true) || it.genre.contains(term, ignoreCase = true) }
                || venue.name.contains(term, ignoreCase = true)
    }

    fun isPopulated(): Boolean {
        return artists.isNotEmpty() && artists.all { it.isPopulated() } && venue.isPopulated()
    }
}