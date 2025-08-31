package com.bldover.beacon.data.model.event

import com.bldover.beacon.data.dto.EventDetailDto
import com.bldover.beacon.data.dto.EventRankDto
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.artist.ArtistRank
import com.bldover.beacon.data.model.venue.Venue
import com.bldover.beacon.data.util.dateFormatter
import java.time.LocalDate

data class EventDetail(
    val id: EventId,
    val name: String,
    val artists: List<Artist>,
    val venue: Venue,
    val date: LocalDate,
    val purchased: Boolean,
    val rank: Float? = null,
    val artistRanks: List<ArtistRank>? = null
) {
    constructor(event: EventDetailDto) : this(
        id = EventId(event.event.id),
        name = event.name,
        artists = event.event.artists,
        venue = Venue(event.event.venue),
        date = LocalDate.parse(event.event.date, dateFormatter),
        purchased = event.event.purchased
    )

    constructor(event: EventRankDto) : this(
        id = EventId(event.event.id),
        name = event.name,
        artists = event.event.artists,
        venue = Venue(event.event.venue),
        date = LocalDate.parse(event.event.date, dateFormatter),
        purchased = event.event.purchased,
        rank = event.ranks.rank,
        artistRanks = event.ranks.artistRanks.map { ArtistRank(it.key, it.value) }
    )

    fun asEvent(): Event {
        return Event(
            id = id,
            artists = artists,
            venue = venue,
            date = date,
            purchased = purchased
        )
    }

    fun hasMatch(term: String): Boolean {
        return artists.any { it.name.contains(term, ignoreCase = true) || it.genres.hasGenre(term) }
                || name.contains(term, ignoreCase = true)
                || venue.name.contains(term, ignoreCase = true)
    }
}