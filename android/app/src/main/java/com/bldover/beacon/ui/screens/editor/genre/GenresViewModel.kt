package com.bldover.beacon.ui.screens.editor.genre

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.repository.ArtistRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class GenreState {
    data object Loading : GenreState()
    data class Success(
        val genres: List<String>,
        val filtered: List<String>
    ) : GenreState()
    data class Error(val message: String) : GenreState()
}

@HiltViewModel
class GenresViewModel @Inject constructor(
    private val artistRepository: ArtistRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow<GenreState>(GenreState.Loading)
    val uiState: StateFlow<GenreState> = _uiState.asStateFlow()

    init {
        loadGenres()
    }

    fun loadGenres() {
        Timber.i("Loading genres")
        viewModelScope.launch {
            _uiState.value = GenreState.Loading
            try {
                val artists = artistRepository.getArtists()
                val allGenres = artists.flatMap { artist ->
                    val genres = mutableListOf<String>()
                    genres.addAll(artist.genres.spotify)
                    genres.addAll(artist.genres.lastFm)
                    genres.addAll(artist.genres.user)
                    genres
                }.distinct().sorted()
                _uiState.value = GenreState.Success(allGenres, allGenres)
                Timber.i("Loaded ${allGenres.size} distinct genres")
            } catch (e: Exception) {
                Timber.e(e, "Failed to load genres")
                _uiState.value = GenreState.Error(e.message ?: "unknown error")
            }
        }
    }

    fun resetFilter() {
        Timber.d("resetting genres filter")
        if (_uiState.value !is GenreState.Success) {
            return
        }
        val currentState = _uiState.value as GenreState.Success
        _uiState.value = GenreState.Success(
            genres = currentState.genres,
            filtered = currentState.genres
        )
    }

    fun applyFilter(searchTerm: String) {
        when (_uiState.value) {
            is GenreState.Success -> {
                val currentState = _uiState.value as GenreState.Success
                _uiState.value = GenreState.Success(
                    currentState.genres,
                    currentState.genres.filter { 
                        it.contains(searchTerm, ignoreCase = true) 
                    }
                )
            }
            else -> return
        }
    }
}