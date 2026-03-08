package com.bldover.beacon.ui.screens.editor.album

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.album.Album
import com.bldover.beacon.data.repository.AlbumRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class AlbumState {
    data object Loading : AlbumState()
    data class Success(
        val albums: List<Album>,
        val filtered: List<Album>
    ) : AlbumState()
    data class Error(val message: String) : AlbumState()
}

@HiltViewModel
class AlbumsViewModel @Inject constructor(
    private val albumRepository: AlbumRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow<AlbumState>(AlbumState.Loading)
    val uiState: StateFlow<AlbumState> = _uiState.asStateFlow()

    init {
        loadAlbums()
    }

    fun loadAlbums() {
        Timber.i("Loading albums")
        viewModelScope.launch {
            _uiState.value = AlbumState.Loading
            try {
                val albums = albumRepository.getAlbums()
                    .sortedWith(compareBy({ it.artists.firstOrNull()?.name ?: "" }, { it.year }))
                _uiState.value = AlbumState.Success(albums, albums)
                Timber.i("Loaded ${albums.size} albums")
            } catch (e: Exception) {
                Timber.e(e, "Failed to load albums")
                _uiState.value = AlbumState.Error(e.message ?: "unknown error")
            }
        }
    }

    fun resetFilter() {
        if (_uiState.value !is AlbumState.Success) return
        val state = _uiState.value as AlbumState.Success
        _uiState.value = AlbumState.Success(albums = state.albums, filtered = state.albums)
    }

    fun applyFilter(searchTerm: String) {
        if (_uiState.value !is AlbumState.Success) return
        val allAlbums = (_uiState.value as AlbumState.Success).albums
        _uiState.value = AlbumState.Success(
            allAlbums,
            allAlbums.filter { it.hasMatch(searchTerm) }
        )
    }

    fun addAlbum(
        album: Album,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            if (!album.isPopulated()) {
                onError("Album is missing required fields")
                return@launch
            }
            try {
                albumRepository.addAlbum(album)
                onSuccess()
                loadAlbums()
            } catch (e: Exception) {
                Timber.e(e, "Failed to add album $album")
                onError("Error saving album ${album.name}, try again later")
            }
        }
    }

    fun updateAlbum(
        album: Album,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            if (!album.isPopulated()) {
                onError("Album is missing required fields")
                return@launch
            }
            try {
                albumRepository.updateAlbum(album)
                onSuccess()
                loadAlbums()
            } catch (e: Exception) {
                Timber.e(e, "Failed to update album $album")
                onError("Error saving album ${album.name}, try again later")
            }
        }
    }

    fun deleteAlbum(
        album: Album,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                albumRepository.deleteAlbum(album)
                onSuccess()
                loadAlbums()
            } catch (e: Exception) {
                Timber.e(e, "Failed to delete album $album")
                onError("Error deleting album ${album.name}, try again later")
            }
        }
    }
}
