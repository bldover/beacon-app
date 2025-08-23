package com.bldover.beacon.ui.screens.editor.artist

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.Direction
import com.bldover.beacon.data.model.OrderField
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.artist.ArtistOrdering
import com.bldover.beacon.data.repository.ArtistRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class ArtistState {
    data object Loading : ArtistState()
    data class Success(
        val artists: List<Artist>,
        val filtered: List<Artist>
    ) : ArtistState()
    data class Error(val message: String) : ArtistState()
}

@HiltViewModel
class ArtistsViewModel @Inject constructor(
    private val artistRepository: ArtistRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow<ArtistState>(ArtistState.Loading)
    val uiState: StateFlow<ArtistState> = _uiState.asStateFlow()
    private val _ordering = MutableStateFlow(ArtistOrdering(OrderField.NAME, Direction.ASCENDING))
    val ordering: StateFlow<ArtistOrdering> = _ordering.asStateFlow()

    init {
        loadArtists()
    }

    fun loadArtists() {
        Timber.i("Loading artists")
        viewModelScope.launch {
            _uiState.value = ArtistState.Loading
            try {
                val artists = artistRepository.getArtists().sortedBy(Artist::name)
                _uiState.value = ArtistState.Success(artists, artists)
                Timber.i("Loaded ${artists.size} artists")
                applyOrdering(_ordering.value)
            } catch (e: Exception) {
                Timber.e(e,"Failed to load artists")
                _uiState.value = ArtistState.Error(e.message ?: "unknown error")
            }
        }
    }

    fun resetFilter() {
        Timber.d("resetting artists filter")
        if (_uiState.value !is ArtistState.Success) {
            Timber.d("resetting artists filter - not a success state")
            return
        }
        _uiState.value = ArtistState.Success(
            artists = (_uiState.value as ArtistState.Success).artists,
            filtered = (_uiState.value as ArtistState.Success).artists
        )
        Timber.d("resetting artists filter - done")
    }

    fun applyFilter(searchTerm: String) {
        when (_uiState.value) {
            is ArtistState.Success -> {
                val allArtists = (_uiState.value as ArtistState.Success).artists
                _uiState.value = ArtistState.Success(
                    allArtists,
                    allArtists.filter { it.hasMatch(searchTerm) }
                )
            }
            else -> return
        }
    }

    fun applyOrdering(ordering: ArtistOrdering) {
        Timber.i("Sorting artists events by $ordering")
        when (_uiState.value) {
            is ArtistState.Success -> {
                val state = (_uiState.value as ArtistState.Success)
                _uiState.value = ArtistState.Success(
                    state.artists.sortedWith(Comparator(ordering::compare)),
                    state.filtered.sortedWith(Comparator(ordering::compare))
                )
                _ordering.value = ordering
            }
            else -> {
                Timber.w("Sorting artists by $ordering - not in success state")
                return
            }
        }
    }

    fun addArtist(
        artist: Artist,
        onSuccess: (Artist) -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                val createdArtist = artistRepository.addArtist(artist)
                onSuccess(createdArtist)
                loadArtists()
            } catch (e: Exception) {
                Timber.e(e,"Failed to add artist $artist")
                onError("Error saving artist ${artist.name}, try again later")
            }
        }
    }

    fun updateArtist(
        artist: Artist,
        onSuccess: (Artist) -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                val updatedArtist = artistRepository.updateArtist(artist)
                onSuccess(updatedArtist)
                loadArtists()
            } catch (e: Exception) {
                Timber.e(e,"Failed to update artist $artist")
                onError("Error saving artist ${artist.name}, try again later")
            }
        }
    }

    fun upsertArtist(
        artist: Artist,
        onSuccess: (Artist) -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        if (artist.id.primary == null) {
            addArtist(artist, onSuccess, onError)
        } else {
            updateArtist(artist, onSuccess, onError)
        }
    }

    fun deleteArtist(
        artist: Artist,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                artistRepository.deleteArtist(artist)
                onSuccess()
                loadArtists()
            } catch (e: Exception) {
                Timber.e(e, "Failed to delete artist $artist")
                onError("Error deleting artist ${artist.name}, try again later")
            }
        }
    }
}