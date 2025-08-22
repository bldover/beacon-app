package com.bldover.beacon.data.dto

data class ArtistRankDto(
    val rank: Float,
    val related: List<String>? = null
)