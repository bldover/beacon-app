package com.bldover.beacon.ui.screens.editor.genre

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.repository.ArtistRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

@HiltViewModel
class GenreSelectorViewModel @Inject constructor(
    private val artistRepository: ArtistRepository
) : ViewModel() {

    private var onSelect: (String) -> Unit = {}
    private var currentArtist: Artist? = null
    
    private val _allUserGenres = MutableStateFlow<List<String>>(emptyList())
    val allUserGenres: StateFlow<List<String>> = _allUserGenres.asStateFlow()
    
    private val _filteredUserGenres = MutableStateFlow<List<String>>(emptyList())
    val filteredUserGenres: StateFlow<List<String>> = _filteredUserGenres.asStateFlow()
    
    private val _isFiltering = MutableStateFlow(false)
    val isFiltering: StateFlow<Boolean> = _isFiltering.asStateFlow()

    fun launchSelector(
        navController: NavController,
        artist: Artist?,
        onSelect: (String) -> Unit
    ) {
        Timber.d("launching genre selector for artist ${artist?.name}")
        this.onSelect = onSelect
        this.currentArtist = artist
        loadAllUserGenres()
        navController.navigate(Screen.SELECT_GENRE.name)
    }

    fun selectGenre(genre: String) {
        Timber.d("genre selector - selected genre $genre")
        onSelect(genre)
    }
    
    fun getCurrentArtist(): Artist? = currentArtist
    
    fun applyFilter(searchTerm: String) {
        val isFiltering = searchTerm.isNotBlank()
        _isFiltering.value = isFiltering
        _filteredUserGenres.value = if (isFiltering) {
            _allUserGenres.value.filter { it.contains(searchTerm, ignoreCase = true) }
        } else {
            _allUserGenres.value
        }
    }
    
    private fun loadAllUserGenres() {
        viewModelScope.launch {
            try {
                val artists = artistRepository.getArtists()
                val userGenres = artists.flatMap { it.genres.user }.distinct().sorted()
                _allUserGenres.value = userGenres
                _filteredUserGenres.value = userGenres
                _isFiltering.value = false
            } catch (e: Exception) {
                Timber.e(e, "Failed to load user genres")
                _allUserGenres.value = emptyList()
                _filteredUserGenres.value = emptyList()
            }
        }
    }
}