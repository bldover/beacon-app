package com.bldover.beacon.data.model.event

import com.bldover.beacon.data.model.venue.RawVenue
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.artist.RawArtist
import com.bldover.beacon.data.util.dateFormatter

data class RawEvent(
    val id: String,
    val mainAct: RawArtist?,
    val openers: List<RawArtist>,
    val venue: RawVenue,
    val date: String,
    val purchased: Boolean,
    val tmId: String
) {
    constructor(event: Event): this(
        id = event.id ?: "",
        mainAct = event.artists.find { it.headliner }?.let { RawArtist(it) },
        openers = event.artists.filter { !it.headliner }.map { RawArtist(it) },
        venue = RawVenue(event.venue),
        date = event.date.format(dateFormatter),
        purchased = event.purchased,
        tmId = event.ticketmasterId ?: ""
    )

    constructor(eventDetail: EventDetail): this(
        id = eventDetail.id,
        mainAct = eventDetail.artists.find { it.headliner }?.let { RawArtist(it) },
        openers = eventDetail.artists.filter { !it.headliner }.map { RawArtist(it) },
        venue = RawVenue(eventDetail.venue),
        date = eventDetail.date.format(dateFormatter),
        purchased = eventDetail.purchased,
        tmId = eventDetail.ticketmasterId ?: ""
    )

    val artists: List<Artist>
        get() = buildList {
            // the backend (for now) returns an "empty" mainAct if not present, TODO: fix and make not nullable
            mainAct?.takeIf { it.name.isNotEmpty() }
                ?.let { add(Artist(artist = it, headliner = true)) }
            addAll(openers.map { Artist(artist = it) })
        }
}