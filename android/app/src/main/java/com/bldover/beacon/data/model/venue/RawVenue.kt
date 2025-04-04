package com.bldover.beacon.data.model.venue

data class RawVenue(
    var id: String,
    var name: String,
    var city: String,
    var state: String
) {
    constructor(venue: Venue) : this(
        id = venue.id ?: "",
        name = venue.name,
        city = venue.city,
        state = venue.state
    )
}