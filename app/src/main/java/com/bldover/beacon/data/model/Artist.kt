package com.bldover.beacon.data.model

data class Artist(
    val id: String,
    val name: String,
    val genre: String,
    val headliner: Boolean = false
) {
    constructor(artist: RawArtist, headliner: Boolean = false) : this(
        id = artist.id,
        name = artist.name,
        genre = artist.genre,
        headliner = headliner
    )
}