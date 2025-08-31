package com.bldover.beacon.data.model.event

import com.bldover.beacon.data.dto.EventIdDto

data class EventId (
    var primary: String? = null,
    var ticketmaster: String? = null
) {
    constructor(id: EventIdDto) : this(
        primary = id.primary.ifBlank { null },
        ticketmaster = id.ticketmaster.ifBlank { null }
    )
}