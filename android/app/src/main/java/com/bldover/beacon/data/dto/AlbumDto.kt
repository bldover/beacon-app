package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.album.Album

data class AlbumDto(
    val id: String,
    val name: String,
    val artists: List<ArtistDto>,
    val year: Int,
    val signed: Boolean,
    val wishlisted: Boolean,
    val variant: String,
    val format: String,
    val notes: String,
    val coverImageUri: String?
) {
    constructor(album: Album) : this(
        id = album.id ?: "",
        name = album.name,
        artists = album.artists.map { ArtistDto(it) },
        year = album.year,
        signed = album.signed,
        wishlisted = album.wishlisted,
        variant = album.variant,
        format = album.format.displayName,
        notes = album.notes,
        coverImageUri = album.coverImageUri
    )
}
