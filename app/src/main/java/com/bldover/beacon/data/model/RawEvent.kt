package com.bldover.beacon.data.model

import com.bldover.beacon.data.util.dateFormatter

data class RawEvent(
    val id: String,
    val mainAct: RawArtist?,
    val openers: List<RawArtist>,
    val venue: Venue,
    val date: String,
    val purchased: Boolean
) {
    constructor(event: Event): this(
        id = event.id ?: "",
        mainAct = event.artists.find { it.headliner }?.let { RawArtist(it) },
        openers = event.artists.filter { !it.headliner }.map { RawArtist(it) },
        venue = event.venue,
        date = event.date.format(dateFormatter),
        purchased = event.purchased
    )

    constructor(eventDetail: EventDetail): this(
        id = eventDetail.id,
        mainAct = eventDetail.artists.find { it.headliner }?.let { RawArtist(it) },
        openers = eventDetail.artists.filter { !it.headliner }.map { RawArtist(it) },
        venue = eventDetail.venue,
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