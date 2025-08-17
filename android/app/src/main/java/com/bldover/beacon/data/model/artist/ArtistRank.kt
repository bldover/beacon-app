package com.bldover.beacon.data.model.artist

import com.bldover.beacon.data.dto.ArtistRankDto

data class ArtistRank(
    val artist: Artist,
    val rank: Float,
    val relatedArtists: List<String>? = null
) {
    constructor(artist: ArtistRankDto) : this(
        artist = Artist(artist.artist),
        rank = artist.rank,
        relatedArtists = artist.related.ifEmpty { null }
    )
}