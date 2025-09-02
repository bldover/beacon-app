package com.bldover.beacon.ui.screens.utility

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.repository.GenreRepository
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
class GenreManagerViewModel @Inject constructor(
    private val genreRepository: GenreRepository
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
                val genreResponse = genreRepository.getGenres()
                val sortedUserGenres = genreResponse.user.sorted()
                _uiState.value = GenreState.Success(sortedUserGenres, sortedUserGenres)
                Timber.i("Loaded ${sortedUserGenres.size} user genres")
            } catch (e: Exception) {
                Timber.e(e, "Failed to load genres")
                _uiState.value = GenreState.Error(e.message ?: "unknown error")
            }
        }
    }

    fun resetFilter() {
        Timber.d("resetting genres filter")
        if (_uiState.value !is GenreState.Success) {
            Timber.d("resetting genres filter - not a success state")
            return
        }
        _uiState.value = GenreState.Success(
            genres = (_uiState.value as GenreState.Success).genres,
            filtered = (_uiState.value as GenreState.Success).genres
        )
        Timber.d("resetting genres filter - done")
    }

    fun applyFilter(searchTerm: String) {
        when (_uiState.value) {
            is GenreState.Success -> {
                val allGenres = (_uiState.value as GenreState.Success).genres
                _uiState.value = GenreState.Success(
                    allGenres,
                    if (searchTerm.isBlank()) {
                        allGenres
                    } else {
                        allGenres.filter { it.contains(searchTerm, ignoreCase = true) }
                    }
                )
            }
            else -> return
        }
    }
}