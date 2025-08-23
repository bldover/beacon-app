package com.bldover.beacon.ui.screens.editor.artist

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.util.fromCommaSeparatedString
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import javax.inject.Inject

@HiltViewModel
class ArtistEditorViewModel @Inject constructor() : ViewModel() {

    private val _artistState = MutableStateFlow(Artist())
    val artistState = _artistState.asStateFlow()

    private var onSave: (Artist) -> Unit = {}

    fun launchEditor(
        navController: NavController,
        artist: Artist? = null,
        onSave: (Artist) -> Unit,
    ) {
        this.onSave = onSave
        _artistState.value = artist?.deepCopy() ?: Artist()
        navController.navigate(Screen.EDIT_ARTIST.name)
    }

    fun updateName(name: String) {
        _artistState.value = _artistState.value.copy(name = name)
    }

    fun updateUserGenres(genres: String) {
        val newGenres = _artistState.value.genres.copy(user = fromCommaSeparatedString(genres))
        _artistState.value = _artistState.value.copy(genres = newGenres)
    }

    fun addGenre(genre: String) {
        val currentUserGenres = _artistState.value.genres.user.toMutableList()
        val allCurrentGenres = _artistState.value.genres.getGenres()
        
        // If this genre is not already in the current active genres, add to user genres
        if (!allCurrentGenres.contains(genre)) {
            // Start with current active genres if user genres is empty
            val newUserGenres = if (currentUserGenres.isEmpty()) {
                allCurrentGenres.toMutableList().apply { add(genre) }
            } else {
                currentUserGenres.apply { add(genre) }
            }
            val newGenres = _artistState.value.genres.copy(user = newUserGenres)
            _artistState.value = _artistState.value.copy(genres = newGenres)
        }
    }

    fun removeGenre(genre: String) {
        val currentUserGenres = _artistState.value.genres.user.toMutableList()
        val allCurrentGenres = _artistState.value.genres.getGenres()
        
        // Copy current active genres to user genres if user genres is empty, then remove the genre
        val newUserGenres = if (currentUserGenres.isEmpty()) {
            allCurrentGenres.toMutableList().apply { remove(genre) }
        } else {
            currentUserGenres.apply { remove(genre) }
        }
        val newGenres = _artistState.value.genres.copy(user = newUserGenres)
        _artistState.value = _artistState.value.copy(genres = newGenres)
    }

    fun onSave() {
        onSave(_artistState.value)
    }
}