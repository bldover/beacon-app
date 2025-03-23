package com.bldover.beacon.data.repository

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import com.bldover.beacon.data.model.ColorScheme
import com.bldover.beacon.data.model.UserSettings
import com.bldover.beacon.data.model.UserSettingsDefaults
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import timber.log.Timber

private object PreferencesKeys {
    val COLOR_SCHEME = stringPreferencesKey("colorScheme")
    val START_SCREEN = stringPreferencesKey("startScreen")
}

class UserSettingsRepository(private val dataStore: DataStore<Preferences>) {

    val userSettingsFlow: Flow<UserSettings> = dataStore.data.map { preferences ->
        val colorScheme = preferences[PreferencesKeys.COLOR_SCHEME]?.let { ColorScheme.valueOf(it) }
            ?: ColorScheme.SYSTEM_DEFAULT
        val startScreen = preferences[PreferencesKeys.START_SCREEN] ?: UserSettingsDefaults.START_SCREEN
        UserSettings(colorScheme, startScreen)
    }

    suspend fun updateColorScheme(colorScheme: ColorScheme) {
        Timber.d("updating color scheme setting: $colorScheme")
        dataStore.edit { preferences ->
            preferences[PreferencesKeys.COLOR_SCHEME] = colorScheme.name
        }
        Timber.d("updated color scheme setting: ${dataStore.data.first()[PreferencesKeys.COLOR_SCHEME]}")
    }

    suspend fun updateStartScreen(startScreen: String) {
        Timber.d("updating start screen setting: $startScreen")
        dataStore.edit { preferences ->
            preferences[PreferencesKeys.START_SCREEN] = startScreen
        }
        Timber.d("updated start screen setting: ${dataStore.data.first()[PreferencesKeys.START_SCREEN]}")
    }
}