package com.bldover.beacon.data.model.artist

import com.bldover.beacon.data.dto.ArtistDto

data class Artist(
    var id: ArtistId = ArtistId(),
    var name: String = "",
    var genres: Genres = Genres(),
    var headliner: Boolean = false
) {
    constructor(
        artist: ArtistDto,
        headliner: Boolean = false,
    ) : this(
        id = ArtistId(artist.id),
        name = artist.name,
        genres = Genres(artist.genres),
        headliner = headliner
    )

    fun hasMatch(searchTerm: String): Boolean {
        return name.contains(searchTerm, ignoreCase = true)
                || genres.hasGenre(searchTerm)
    }

    fun isPopulated(): Boolean {
        return name.isNotEmpty()
    }

    fun deepCopy(): Artist {
        return Artist(
            id = id.copy(),
            name = name,
            genres = genres.deepCopy(),
            headliner = headliner
        )
    }
}