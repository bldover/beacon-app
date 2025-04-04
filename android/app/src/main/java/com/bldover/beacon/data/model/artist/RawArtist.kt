package com.bldover.beacon.data.model.artist

data class RawArtist(
    val id: String,
    val name: String,
    val genre: String
) {
    constructor(artist: Artist) : this(
        id = artist.id ?: "",
        name = artist.name,
        genre = if (artist.genreSet) artist.genre else GENRE_UNSPECIFIED
    )
}