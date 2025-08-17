package com.bldover.beacon.data.model.artist

import com.bldover.beacon.data.dto.ArtistIdDto

data class Id (
    var primary: String? = null,
    var ticketmaster: String? = null,
    var spotify: String? = null,
    var musicbrainz: String? = null
) {
    constructor(idDto: ArtistIdDto) : this(
        primary = idDto.primary.ifBlank { null },
        ticketmaster = idDto.ticketmaster.ifBlank { null },
        spotify = idDto.spotify.ifBlank { null },
        musicbrainz = idDto.musicbrainz.ifBlank { null }
    )
}