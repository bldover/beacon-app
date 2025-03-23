package com.bldover.beacon.data.model

data class Venue(
    var id: String? = null,
    var name: String,
    var city: String,
    var state: String
) {
    constructor(venue: RawVenue) : this(
        id = venue.id.ifBlank { null },
        name = venue.name,
        city = venue.city,
        state = venue.state
    )

    fun hasMatch(searchTerm: String): Boolean {
        return name.contains(searchTerm, ignoreCase = true)
    }

    fun isPopulated(): Boolean {
        return name.isNotEmpty() && city.isNotEmpty() && state.isNotEmpty()
    }
}