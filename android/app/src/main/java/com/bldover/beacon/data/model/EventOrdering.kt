package com.bldover.beacon.data.model

data class EventOrdering(
    val option: OrderType = OrderType.DATE,
    val order: Order = Order.DESCENDING
) : Comparator<Event> {
    override fun compare(o1: Event, o2: Event): Int {
        return when (option) {
            OrderType.VENUE -> when (order) {
                Order.ASCENDING -> o1.venue.name.compareTo(o2.venue.name)
                Order.DESCENDING -> o2.venue.name.compareTo(o1.venue.name)
            }
            OrderType.DATE -> when (order) {
                Order.ASCENDING -> o1.date.compareTo(o2.date)
                Order.DESCENDING -> o2.date.compareTo(o1.date)
            }
        }
    }
}

enum class OrderType(val label: String) {
    DATE("Date"),
    VENUE("Venue");
}

enum class Order {
    ASCENDING,
    DESCENDING
}