package com.bldover.beacon.data.model.analytics

import com.bldover.beacon.data.model.Direction
import com.bldover.beacon.data.model.OrderField
import com.bldover.beacon.data.model.Ordering

class AnalyticsOrdering(
    option: OrderField = OrderField.COUNT,
    order: Direction = Direction.DESCENDING
) : Ordering<AnalyticsCount>(option, order) {
    override fun compare(o1: AnalyticsCount, o2: AnalyticsCount): Int {
        return when (option) {
            OrderField.COUNT -> {
                val byCount = when (order) {
                    Direction.ASCENDING -> o1.count.compareTo(o2.count)
                    Direction.DESCENDING -> o2.count.compareTo(o1.count)
                }
                if (byCount != 0) byCount else o1.name.compareTo(o2.name, ignoreCase = true)
            }
            OrderField.NAME -> when (order) {
                Direction.ASCENDING -> o1.name.compareTo(o2.name, ignoreCase = true)
                Direction.DESCENDING -> o2.name.compareTo(o1.name, ignoreCase = true)
            }
            else -> 0
        }
    }
}
