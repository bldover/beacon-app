package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.event.EventId

data class EventIdDto(
    var primary: String,
    var ticketmaster: String
) {
    constructor(id: EventId?) : this(
        primary = id?.primary ?: "",
        ticketmaster = id?.ticketmaster ?: ""
    )
}