package com.bldover.beacon.data.model

enum class Screen(val title: String, val subScreen: Boolean = false) {
    CONCERT_PLANNER("Planner"),
    PLANNER_DETAIL("Event Detail"),
    CONCERT_HISTORY("Concert History"),
    HISTORY_DETAIL("Event Detail"),
    UPCOMING_EVENTS("Upcoming Events"),
    UPCOMING_DETAIL("Event Detail"),
    UTILITIES("Utilities"),
    USER_SETTINGS("User Settings", true);

    companion object {
        fun fromTitle(title: String): Screen {
            return entries.find { it.title == title }
                ?: throw IllegalArgumentException("No ActiveScreen found for title $title")
        }

        fun fromOrDefault(
            name: String?,
            default: Screen = CONCERT_PLANNER
        ): Screen {
            return entries.find { it.name == name } ?: default
        }

        fun majorScreens(): List<Screen> {
            return listOf(CONCERT_HISTORY, CONCERT_PLANNER, UPCOMING_EVENTS, UTILITIES)
        }
    }
}
