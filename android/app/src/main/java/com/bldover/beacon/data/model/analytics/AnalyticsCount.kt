package com.bldover.beacon.data.model.analytics

import com.bldover.beacon.data.dto.CountDto

data class AnalyticsCount(
    val key: String,
    val name: String,
    val count: Int
) {
    constructor(dto: CountDto) : this(
        key = dto.key,
        name = dto.name,
        count = dto.count
    )
}
