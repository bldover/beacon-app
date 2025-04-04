package com.bldover.beacon.data.model.event

data class RawEventDetail(
    val event: RawEvent,
    val name: String,
    val price: String
) {
    constructor(eventDetail: EventDetail): this(
        event = RawEvent(eventDetail),
        name = eventDetail.name,
        price = eventDetail.formattedPrice
    )
}