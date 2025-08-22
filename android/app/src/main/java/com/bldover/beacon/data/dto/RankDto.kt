package com.bldover.beacon.data.dto

data class RankDto(
    val rank: Float,
    val recommendation: String,
    val artistRanks: Map<String, ArtistRankDto>,
)
