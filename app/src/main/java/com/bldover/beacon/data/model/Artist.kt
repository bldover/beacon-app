package com.bldover.beacon.data.model

data class Artist(
    var id: String? = null,
    var name: String,
    var genre: String,
    var headliner: Boolean = false
) {
    constructor(artist: RawArtist, headliner: Boolean = false) : this(
        id = artist.id,
        name = artist.name,
        genre = artist.genre,
        headliner = headliner
    )

    fun hasMatch(searchTerm: String): Boolean {
        return name.contains(searchTerm, ignoreCase = true)
    }

    fun isPopulated(): Boolean {
        return name.isNotEmpty() && genre.isNotEmpty()
    }
}