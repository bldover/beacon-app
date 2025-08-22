package com.bldover.beacon.data.dto

data class EventRankDto(
    var name: String,
    val event: EventDto,
    val ranks: RankDto
)