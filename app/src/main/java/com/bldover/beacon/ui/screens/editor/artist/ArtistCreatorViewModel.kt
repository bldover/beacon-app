package com.bldover.beacon.ui.screens.editor.artist

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Artist
import com.bldover.beacon.data.model.Screen
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import javax.inject.Inject

@HiltViewModel
class ArtistCreatorViewModel @Inject constructor() : ViewModel() {

    private val _artistState = MutableStateFlow(Artist(name = "", genre = ""))
    val artistState = _artistState.asStateFlow()

    private var onSave: (Artist) -> Unit = {}

    fun launchCreator(
        navController: NavController,
        onSave: (Artist) -> Unit,
        name: String = "",
    ) {
        this.onSave = onSave
        _artistState.value = Artist(name = name, genre = "")
        navController.navigate(Screen.CREATE_ARTIST.name)
    }

    fun updateName(name: String) {
        _artistState.value = _artistState.value.copy(name = name)
    }

    fun updateGenre(genre: String) {
        _artistState.value = _artistState.value.copy(genre = genre)
    }

    fun onSave() {
        onSave(_artistState.value)
    }
}