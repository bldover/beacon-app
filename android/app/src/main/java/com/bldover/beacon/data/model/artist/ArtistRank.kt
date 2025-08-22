package com.bldover.beacon.data.model.artist

import com.bldover.beacon.data.dto.ArtistRankDto

data class ArtistRank(
    val name: String,
    val rank: Float,
    val relatedArtists: List<String>? = null
) {
    constructor(name: String, artist: ArtistRankDto) : this(
        name = name,
        rank = artist.rank,
        relatedArtists = artist.related?.ifEmpty { null }
    )
}