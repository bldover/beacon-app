package com.bldover.beacon.data.model.venue

import com.bldover.beacon.data.model.ordering.Direction
import com.bldover.beacon.data.model.ordering.OrderField
import com.bldover.beacon.data.model.ordering.Ordering

class VenueOrdering(
    option: OrderField = OrderField.NAME,
    order: Direction = Direction.ASCENDING
) : Ordering<Venue>(option, order) {
    override fun compare(o1: Venue, o2: Venue): Int {
        return when (option) {
            OrderField.NAME -> when (order) {
                Direction.ASCENDING -> o1.name.compareTo(o2.name)
                Direction.DESCENDING -> o2.name.compareTo(o1.name)
            }
            OrderField.CITY -> when (order) {
                Direction.ASCENDING -> o1.city.compareTo(o2.city)
                Direction.DESCENDING -> o2.city.compareTo(o1.city)
            }
            else -> 0
        }
    }
}