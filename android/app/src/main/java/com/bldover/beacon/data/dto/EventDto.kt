package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.event.Event
import com.bldover.beacon.data.model.event.EventDetail
import com.bldover.beacon.data.util.dateFormatter

data class EventDto(
    val id: EventIdDto,
    val mainAct: ArtistDto?,
    val openers: List<ArtistDto>,
    val venue: VenueDto,
    val date: String,
    val purchased: Boolean
) {
    constructor(event: Event): this(
        id = EventIdDto(event.id),
        mainAct = event.artists.find { it.headliner }?.let { ArtistDto(it) },
        openers = event.artists.filter { !it.headliner }.map { ArtistDto(it) },
        venue = VenueDto(event.venue),
        date = event.date.format(dateFormatter),
        purchased = event.purchased
    )

    constructor(eventDetail: EventDetail): this(
        id = EventIdDto(eventDetail.id),
        mainAct = eventDetail.artists.find { it.headliner }?.let { ArtistDto(it) },
        openers = eventDetail.artists.filter { !it.headliner }.map { ArtistDto(it) },
        venue = VenueDto(eventDetail.venue),
        date = eventDetail.date.format(dateFormatter),
        purchased = eventDetail.purchased
    )

    val artists: List<Artist>
        get() = buildList {
            // the backend (for now) returns an "empty" mainAct if not present, TODO: fix and make not nullable
            mainAct?.takeIf { it.name.isNotEmpty() }
                ?.let { add(Artist(artist = it, headliner = true)) }
            addAll(openers.map { Artist(artist = it) })
        }
}