package com.bldover.beacon.data.dto

data class EventRankDto(
    val event: EventDetailDto,
    val rank: Float,
    val artistRanks: List<ArtistRankDto>
)