package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.artist.Id

data class ArtistIdDto(
    var primary: String,
    var ticketmaster: String,
    var spotify: String,
    var musicbrainz: String
) {
    constructor(id: Id) : this(
        primary = id.primary ?: "",
        ticketmaster = id.ticketmaster ?: "",
        spotify = id.spotify ?: "",
        musicbrainz = id.musicbrainz ?: ""
    )
}