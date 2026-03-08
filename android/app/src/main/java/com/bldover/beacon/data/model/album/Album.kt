package com.bldover.beacon.data.model.album

import com.bldover.beacon.data.dto.AlbumDto
import com.bldover.beacon.data.model.artist.Artist

data class Album(
    var id: String? = null,
    var name: String = "",
    var artists: List<Artist> = emptyList(),
    var year: Int = 9999,
    var signed: Boolean = false,
    var wishlisted: Boolean = false,
    var variant: String = "",
    var format: AlbumFormat = AlbumFormat.LP,
    var notes: String = "",
    var coverImageUri: String? = null
) {
    constructor(dto: AlbumDto) : this(
        id = dto.id.ifBlank { null },
        name = dto.name,
        artists = dto.artists.map { Artist(it) },
        year = dto.year,
        signed = dto.signed,
        wishlisted = dto.wishlisted,
        variant = dto.variant,
        format = AlbumFormat.fromString(dto.format),
        notes = dto.notes,
        coverImageUri = dto.coverImageUri
    )

    fun isPopulated(): Boolean {
        return name.isNotBlank()
            && artists.isNotEmpty() && artists.all { it.isPopulated() }
            && year in 1000..9999
            && variant.isNotBlank()
    }

    fun hasMatch(searchTerm: String): Boolean {
        return name.contains(searchTerm, ignoreCase = true)
            || artists.any { it.name.contains(searchTerm, ignoreCase = true) }
            || variant.contains(searchTerm, ignoreCase = true)
    }

    fun deepCopy(): Album {
        return Album(
            id = id,
            name = name,
            artists = artists.map { it.deepCopy() },
            year = year,
            signed = signed,
            wishlisted = wishlisted,
            variant = variant,
            format = format,
            notes = notes,
            coverImageUri = coverImageUri
        )
    }
}
