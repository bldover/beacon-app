package com.bldover.beacon.data.model

data class RawEvent(
    val id: String,
    val mainAct: RawArtist?,
    val openers: List<RawArtist>,
    val venue: Venue,
    val date: String,
    val purchased: Boolean
) {
    val artists: List<Artist>
        get() = buildList {
            // the backend (for now) returns an "empty" mainAct if not present, TODO: fix and make not nullable
            mainAct?.takeIf { it.id.isNotEmpty() }
                ?.let { add(Artist(artist = it, headliner = true)) }
            addAll(openers.map { Artist(artist = it) })
        }
}