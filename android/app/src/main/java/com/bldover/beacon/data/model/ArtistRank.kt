package com.bldover.beacon.data.model

data class ArtistRank(
    val artist: Artist,
    val rank: Float,
    val relatedArtists: List<String>? = null
) {
    constructor(artist: RawArtistRank) : this(
        artist = Artist(artist.artist),
        rank = artist.rank,
        relatedArtists = artist.related.ifEmpty { null }
    )
}