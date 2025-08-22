package com.bldover.beacon.data.model.event

import com.bldover.beacon.data.dto.ArtistIdDto

data class EventId (
    var primary: String? = null,
    var ticketmaster: String? = null
) {
    constructor(idDto: ArtistIdDto) : this(
        primary = idDto.primary.ifBlank { null },
        ticketmaster = idDto.ticketmaster.ifBlank { null }
    )
}