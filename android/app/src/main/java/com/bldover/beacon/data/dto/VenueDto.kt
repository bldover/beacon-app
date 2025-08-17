package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.venue.Venue

data class VenueDto(
    var id: VenueIdDto,
    var name: String,
    var city: String,
    var state: String
) {
    constructor(venue: Venue) : this(
        id = VenueIdDto(venue.id),
        name = venue.name,
        city = venue.city,
        state = venue.state
    )
}