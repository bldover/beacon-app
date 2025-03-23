package com.bldover.beacon.data.model

data class Artist(
    var id: String? = null,
    var name: String,
    var genre: String,
    var genreSet: Boolean = true,
    var headliner: Boolean = false
) {
    constructor(
        artist: RawArtist,
        headliner: Boolean = false,
        genreSet: Boolean = true
    ) : this(
        id = artist.id.ifBlank { null },
        name = artist.name,
        genre = artist.genre,
        genreSet = !(artist.genre.isBlank() || artist.genre == GENRE_UNSPECIFIED),
        headliner = headliner
    )

    fun hasMatch(searchTerm: String): Boolean {
        return name.contains(searchTerm, ignoreCase = true)
    }

    fun isPopulated(): Boolean {
        return name.isNotEmpty() && (genre.isNotEmpty() || !genreSet)
    }
}

const val GENRE_UNSPECIFIED = "Unspecified"