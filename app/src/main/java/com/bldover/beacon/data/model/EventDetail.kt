package com.bldover.beacon.data.model

import com.bldover.beacon.data.util.dateFormatter
import java.time.LocalDate

data class EventDetail(
    val id: String,
    val name: String,
    val artists: List<Artist>,
    val venue: Venue,
    val date: LocalDate,
    val purchased: Boolean,
    val price: Float?,
    val ticketmasterId: String? = null
) {
    constructor(event: RawEventDetail) : this(
        id = event.event.id,
        name = event.name,
        artists = event.event.artists,
        venue = event.event.venue,
        date = LocalDate.parse(event.event.date, dateFormatter),
        purchased = event.event.purchased,
        price = event.price.toFloatOrNull(),
        ticketmasterId = event.event.tmId.ifBlank { null }
    )

    fun asEvent(): Event {
        return Event(
            id = id,
            artists = artists,
            venue = venue,
            date = date,
            purchased = purchased,
            ticketmasterId = ticketmasterId
        )
    }

    fun hasMatch(term: String): Boolean = artists.any { it.name.contains(term, ignoreCase = true) }

    val formattedPrice: String
        get() = price?.let { String.format("%.2f", it) } ?: "Unknown"
}