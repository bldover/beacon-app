package com.bldover.beacon.data.model.event

import com.bldover.beacon.data.model.Direction
import com.bldover.beacon.data.model.OrderField
import com.bldover.beacon.data.model.Ordering

class EventOrdering(
    option: OrderField = OrderField.DATE,
    order: Direction = Direction.ASCENDING
) : Ordering<Event>(option, order) {
    override fun compare(o1: Event, o2: Event): Int {
        return when (option) {
            OrderField.VENUE -> when (order) {
                Direction.ASCENDING -> o1.venue.name.compareTo(o2.venue.name)
                Direction.DESCENDING -> o2.venue.name.compareTo(o1.venue.name)
            }
            OrderField.DATE -> when (order) {
                Direction.ASCENDING -> o1.date.compareTo(o2.date)
                Direction.DESCENDING -> o2.date.compareTo(o1.date)
            }
            else -> 0
        }
    }
}