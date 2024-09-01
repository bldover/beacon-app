package com.bldover.beacon.data.model

data class Venue(
    var id: String? = null,
    var name: String,
    var city: String,
    var state: String
) {
    fun hasMatch(searchTerm: String): Boolean {
        return name.contains(searchTerm, ignoreCase = true)
    }

    fun isPopulated(): Boolean {
        return name.isNotEmpty() && city.isNotEmpty() && state.isNotEmpty()
    }
}