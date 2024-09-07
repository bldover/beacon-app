package com.bldover.beacon.data.model

data class RawEventRank(
    val event: RawEventDetail,
    val rank: Float,
    val artistRanks: List<RawArtistRank>
)