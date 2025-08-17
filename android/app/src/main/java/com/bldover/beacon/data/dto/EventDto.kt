package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.event.Event
import com.bldover.beacon.data.model.event.EventDetail
import com.bldover.beacon.data.util.dateFormatter

data class EventDto(
    val id: String,
    val mainAct: ArtistDto?,
    val openers: List<ArtistDto>,
    val venue: VenueDto,
    val date: String,
    val purchased: Boolean,
    val tmId: String
) {
    constructor(event: Event): this(
        id = event.id ?: "",
        mainAct = event.artists.find { it.headliner }?.let { ArtistDto(it) },
        openers = event.artists.filter { !it.headliner }.map { ArtistDto(it) },
        venue = VenueDto(event.venue),
        date = event.date.format(dateFormatter),
        purchased = event.purchased,
        tmId = event.ticketmasterId ?: ""
    )

    constructor(eventDetail: EventDetail): this(
        id = eventDetail.id,
        mainAct = eventDetail.artists.find { it.headliner }?.let { ArtistDto(it) },
        openers = eventDetail.artists.filter { !it.headliner }.map { ArtistDto(it) },
        venue = VenueDto(eventDetail.venue),
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