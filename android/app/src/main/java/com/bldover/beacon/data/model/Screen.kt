package com.bldover.beacon.data.model

enum class Screen(val title: String) {
    CONCERT_PLANNER("Planner"),
    CONCERT_HISTORY("Concert History"),
    UPCOMING_EVENTS("Upcoming Events"),
    UTILITIES("Utilities"),
    USER_SETTINGS("User Settings"),
    EDIT_EVENT("Edit Event"),
    EDIT_VENUE("Edit Venue"),
    EDIT_ARTIST("Edit Artist"),
    SELECT_VENUE("Select Venue"),
    SELECT_ARTIST("Select Artist"),
    SELECT_GENRE("Select Genre"),
    MANAGE_VENUES("Manage Venues"),
    MANAGE_ARTISTS("Manage Artists");

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
