package com.bldover.beacon.data.dto

data class ArtistRankDto(
    val artist: ArtistDto,
    val rank: Float,
    val related: List<String>
)