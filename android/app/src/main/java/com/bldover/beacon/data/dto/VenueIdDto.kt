package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.venue.Id

data class VenueIdDto(
    var primary: String,
    var ticketmaster: String
) {
    constructor(id: Id?) : this(
        primary = id?.primary ?: "",
        ticketmaster = id?.ticketmaster ?: ""
    )
}