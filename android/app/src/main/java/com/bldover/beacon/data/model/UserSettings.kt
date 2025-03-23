package com.bldover.beacon.data.model

data class UserSettings(
    val colorScheme: ColorScheme = UserSettingsDefaults.DARK_MODE,
    val startScreen: String = UserSettingsDefaults.START_SCREEN
)

enum class ColorScheme(val scheme: String) {
    SYSTEM_DEFAULT("System Default"),
    DARK_MODE("Dark Mode"),
    LIGHT_MODE("Light Mode");

    companion object {
        fun from(scheme: String?): ColorScheme {
            return entries.find { it.scheme == scheme }
                ?: throw IllegalArgumentException("Invalid color scheme: $scheme")
        }
    }
}

object UserSettingsDefaults {
    val DARK_MODE = ColorScheme.SYSTEM_DEFAULT
    const val START_SCREEN = ""
}