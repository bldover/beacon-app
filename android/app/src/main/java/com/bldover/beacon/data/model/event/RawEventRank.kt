package com.bldover.beacon.data.model.event

import com.bldover.beacon.data.model.artist.RawArtistRank

data class RawEventRank(
    val event: RawEventDetail,
    val rank: Float,
    val artistRanks: List<RawArtistRank>
)