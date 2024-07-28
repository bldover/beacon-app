package com.bldover.beacon.ui.screens.history

import java.time.LocalDate

data class Artist(val name: String, val genre: String)
data class Venue(val name: String, val city: String, val state: String)

data class Event(
    val artists: List<Artist> = emptyList(),
    val venue: Venue? = null,
    val date: LocalDate? = null,
    val purchased: Boolean = false
)

fun getDummyEvent(): Event {
    return Event(
        artists = listOf(
            Artist(name = "artist name", genre = "artist genre"),
            Artist(name = "artist name 2", genre = "artist genre 2")
        ),
        venue = Venue(
            name = "The Masquerade - Hell",
            city = "Atlanta",
            state = "GA"
        ),
        date = LocalDate.now()
    )
}