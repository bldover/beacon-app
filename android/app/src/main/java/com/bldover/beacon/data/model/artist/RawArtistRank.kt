package com.bldover.beacon.data.model.artist

data class RawArtistRank(
    val artist: RawArtist,
    val rank: Float,
    val related: List<String>
)