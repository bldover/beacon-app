package com.bldover.beacon.ui.screens.editor.artist

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.artist.Artist
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

    fun addGenre(genre: String) {
        val currentUserGenres = _artistState.value.genres.user.toMutableList()
        
        if (!currentUserGenres.contains(genre)) {
            currentUserGenres.add(genre)
            val newGenres = _artistState.value.genres.copy(user = currentUserGenres)
            _artistState.value = _artistState.value.copy(genres = newGenres)
        }
    }

    fun removeGenre(genre: String) {
        val currentUserGenres = _artistState.value.genres.user.toMutableList()
        currentUserGenres.remove(genre)
        val newGenres = _artistState.value.genres.copy(user = currentUserGenres)
        _artistState.value = _artistState.value.copy(genres = newGenres)
    }

    fun onSave() {
        onSave(_artistState.value)
    }
}