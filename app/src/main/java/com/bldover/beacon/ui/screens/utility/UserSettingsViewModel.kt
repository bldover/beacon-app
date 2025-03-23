package com.bldover.beacon.ui.screens.utility

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.ColorScheme
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.UserSettings
import com.bldover.beacon.data.repository.UserSettingsRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class SettingsState {
    data object Loading : SettingsState()
    data class Success(val data: UserSettings) : SettingsState()
}

@HiltViewModel
class UserSettingsViewModel @Inject constructor(
    private val userSettingsRepository: UserSettingsRepository
) : ViewModel() {

    val userSettings = userSettingsRepository.userSettingsFlow
        .map { SettingsState.Success(it) }
        .stateIn(
            scope = viewModelScope,
            started = SharingStarted.Eagerly,
            initialValue = SettingsState.Loading
        )

    fun updateColorScheme(
        colorScheme: ColorScheme,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                userSettingsRepository.updateColorScheme(colorScheme)
                onSuccess()
            } catch (e: Exception) {
                Timber.e(e, "Failed to update color scheme $colorScheme")
                onError("Error updating color scheme, try again later")
            }
        }
    }

    fun updateStartScreen(
        startScreen: Screen,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                userSettingsRepository.updateStartScreen(startScreen.name)
                onSuccess()
            } catch (e: Exception) {
                Timber.e(e, "Failed to update start screen $startScreen")
                onError("Error updating start screen, try again later")
            }
        }
    }
}