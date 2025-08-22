package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.event.EventDetail

data class EventDetailDto(
    val event: EventDto,
    val name: String
) {
    constructor(eventDetail: EventDetail): this(
        event = EventDto(eventDetail),
        name = eventDetail.name
    )
}