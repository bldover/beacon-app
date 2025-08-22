package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.venue.VenueId

data class VenueIdDto(
    var primary: String,
    var ticketmaster: String
) {
    constructor(id: VenueId?) : this(
        primary = id?.primary ?: "",
        ticketmaster = id?.ticketmaster ?: ""
    )
}