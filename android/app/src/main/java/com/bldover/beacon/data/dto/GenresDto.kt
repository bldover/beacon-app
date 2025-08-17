package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.artist.Genres

data class GenresDto(
    var spotify: List<String>?,
    var lastFm: List<String>?,
    var ticketmaster: List<String>?,
    var user: List<String>?
) {
    constructor(genres: Genres) : this(
        spotify = genres.spotify,
        lastFm = genres.lastFm,
        ticketmaster = genres.ticketmaster,
        user = genres.user
    )
}