package com.bldover.beacon.data.model

data class RawArtistRank(
    val artist: RawArtist,
    val rank: Float,
    val related: List<String>
)