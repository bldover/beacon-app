package com.bldover.beacon.data.model.venue

import com.bldover.beacon.data.dto.VenueDto

data class Venue(
    var id: Id = Id(),
    var name: String = "",
    var city: String = "",
    var state: String = ""
) {
    constructor(venue: VenueDto) : this(
        id = Id(venue.id),
        name = venue.name,
        city = venue.city,
        state = venue.state
    )

    fun hasMatch(searchTerm: String): Boolean {
        return name.contains(searchTerm, ignoreCase = true)
                || city.contains(searchTerm, ignoreCase = true)
    }

    fun isPopulated(): Boolean {
        return name.isNotEmpty() && city.isNotEmpty() && state.isNotEmpty()
    }

    fun deepCopy(): Venue {
        return Venue(
            id = id.copy(),
            name = name,
            city = city,
            state = state
        )
    }
}