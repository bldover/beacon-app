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

    fun updateColorScheme(colorScheme: ColorScheme) {
        viewModelScope.launch {
            userSettingsRepository.updateColorScheme(colorScheme)
        }
    }

    fun updateStartScreen(startScreen: Screen) {
        viewModelScope.launch {
            userSettingsRepository.updateStartScreen(startScreen.name)
        }
    }
}