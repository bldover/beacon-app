package com.bldover.beacon.ui.screens.utility

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.SpotifyAuthStatus
import com.bldover.beacon.data.repository.SpotifyRepository
import com.bldover.beacon.data.spotify.SpotifyAuthResult
import com.bldover.beacon.data.spotify.SpotifyAuthResultBus
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharedFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class SpotifyAuthStatusState {
    data object Loading : SpotifyAuthStatusState()
    data class Loaded(val status: SpotifyAuthStatus) : SpotifyAuthStatusState()
    data class Error(val message: String) : SpotifyAuthStatusState()
}

@HiltViewModel
class SpotifyAuthViewModel @Inject constructor(
    private val spotifyRepository: SpotifyRepository
) : ViewModel() {

    val authResults: SharedFlow<SpotifyAuthResult> = SpotifyAuthResultBus.results

    private val _authStatus = MutableStateFlow<SpotifyAuthStatusState>(SpotifyAuthStatusState.Loading)
    val authStatus: StateFlow<SpotifyAuthStatusState> = _authStatus.asStateFlow()

    private val _initialPromptResolved = MutableStateFlow(false)
    val initialPromptResolved: StateFlow<Boolean> = _initialPromptResolved.asStateFlow()

    init {
        loadAuthStatus()
        viewModelScope.launch {
            authResults.collect { loadAuthStatus() }
        }
    }

    fun loadAuthStatus() {
        viewModelScope.launch {
            try {
                _authStatus.value = SpotifyAuthStatusState.Loaded(spotifyRepository.getAuthStatus())
            } catch (e: Exception) {
                Timber.e(e, "Failed to load Spotify auth status")
                _authStatus.value = SpotifyAuthStatusState.Error(e.message ?: "unknown error")
            }
        }
    }

    fun startReauth(
        onAuthUrl: (String) -> Unit,
        onError: (String) -> Unit
    ) {
        viewModelScope.launch {
            try {
                val url = spotifyRepository.startReauth()
                Timber.d("Received Spotify auth URL from backend")
                onAuthUrl(url)
            } catch (e: Exception) {
                Timber.e(e, "Failed to start Spotify reauth")
                onError(e.message ?: "unknown error")
            }
        }
    }

    fun markInitialPromptResolved() {
        _initialPromptResolved.value = true
    }
}
