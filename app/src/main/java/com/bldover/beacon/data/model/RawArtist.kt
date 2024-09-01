package com.bldover.beacon.data.model

data class RawArtist(
    val id: String,
    val name: String,
    val genre: String
) {
    constructor(artist: Artist) : this(
        id = artist.id ?: "",
        name = artist.name,
        genre = artist.genre
    )
}