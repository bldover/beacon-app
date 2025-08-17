package com.bldover.beacon.data.model.venue

import com.bldover.beacon.data.dto.VenueIdDto

data class Id (
    var primary: String? = null,
    var ticketmaster: String? = null
) {
    constructor(idDto: VenueIdDto) : this(
        primary = idDto.primary.ifBlank { null },
        ticketmaster = idDto.ticketmaster.ifBlank { null }
    )
}