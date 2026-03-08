package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.album.Album

data class AlbumDto(
    val id: String,
    val name: String,
    val artist: ArtistDto,
    val year: Int,
    val signed: Boolean
) {
    constructor(album: Album) : this(
        id = album.id ?: "",
        name = album.name,
        artist = ArtistDto(album.artist),
        year = album.year,
        signed = album.signed
    )
}
