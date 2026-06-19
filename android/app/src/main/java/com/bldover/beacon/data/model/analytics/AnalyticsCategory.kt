package com.bldover.beacon.data.model.analytics

enum class AnalyticsCategory(val title: String, val routeKey: String) {
    YEARS("Years", "years"),
    MONTHS("Months", "months"),
    ARTISTS("Artists", "artists"),
    VENUES("Venues", "venues"),
    GENRES("Genres", "genres");

    companion object {
        fun fromRouteKey(key: String): AnalyticsCategory? {
            return entries.find { it.routeKey == key }
        }
    }
}
