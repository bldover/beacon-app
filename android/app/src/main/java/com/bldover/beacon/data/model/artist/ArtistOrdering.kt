package com.bldover.beacon.data.model.artist

import com.bldover.beacon.data.model.ordering.Direction
import com.bldover.beacon.data.model.ordering.OrderField
import com.bldover.beacon.data.model.ordering.Ordering

class ArtistOrdering(
    option: OrderField = OrderField.NAME,
    order: Direction = Direction.ASCENDING
) : Ordering<Artist>(option, order) {
    override fun compare(o1: Artist, o2: Artist): Int {
        return when (option) {
            OrderField.NAME -> when (order) {
                Direction.ASCENDING -> String.CASE_INSENSITIVE_ORDER.compare(o1.name, o2.name)
                Direction.DESCENDING -> String.CASE_INSENSITIVE_ORDER.compare(o2.name, o1.name)
            }
            OrderField.GENRE -> when (order) {
                Direction.ASCENDING -> o1.genre.compareTo(o2.genre)
                Direction.DESCENDING -> o2.genre.compareTo(o1.genre)
            }
            else -> 0
        }
    }
}