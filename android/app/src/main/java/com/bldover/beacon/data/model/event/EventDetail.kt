package com.bldover.beacon.data.model.event

import com.bldover.beacon.data.model.venue.Venue
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.artist.ArtistRank
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
    val ticketmasterId: String? = null,
    val rank: Float? = null,
    val artistRanks: List<ArtistRank>? = null
) {
    constructor(event: RawEventDetail) : this(
        id = event.event.id,
        name = event.name,
        artists = event.event.artists,
        venue = Venue(event.event.venue),
        date = LocalDate.parse(event.event.date, dateFormatter),
        purchased = event.event.purchased,
        price = event.price.toFloatOrNull(),
        ticketmasterId = event.event.tmId.ifBlank { null }
    )

    constructor(event: RawEventRank) : this(
        id = event.event.event.id,
        name = event.event.name,
        artists = event.event.event.artists,
        venue = Venue(event.event.event.venue),
        date = LocalDate.parse(event.event.event.date, dateFormatter),
        purchased = event.event.event.purchased,
        price = event.event.price.toFloatOrNull(),
        ticketmasterId = event.event.event.tmId.ifBlank { null },
        rank = event.rank,
        artistRanks = event.artistRanks.map { ArtistRank(it) }
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

    fun hasMatch(term: String): Boolean {
        return artists.any { it.name.contains(term, ignoreCase = true) || it.genre.contains(term, ignoreCase = true) }
                || name.contains(term, ignoreCase = true)
                || venue.name.contains(term, ignoreCase = true)
    }

    val formattedPrice: String
        get() = price?.let { String.format("%.2f", it) } ?: "Unknown"
}