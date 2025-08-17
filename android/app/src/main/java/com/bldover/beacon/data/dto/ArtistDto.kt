package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.artist.Artist

data class ArtistDto(
    val id: ArtistIdDto,
    val name: String,
    val genres: GenresDto
) {
    constructor(artist: Artist) : this(
        id = ArtistIdDto(artist.id),
        name = artist.name,
        genres = GenresDto(artist.genres)
    )
}