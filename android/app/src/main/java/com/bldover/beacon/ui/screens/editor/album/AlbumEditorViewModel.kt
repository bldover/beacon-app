package com.bldover.beacon.ui.screens.editor.album

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.album.Album
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import java.time.LocalDate
import javax.inject.Inject

@HiltViewModel
class AlbumEditorViewModel @Inject constructor() : ViewModel() {

    private val _albumState = MutableStateFlow(Album())
    val albumState = _albumState.asStateFlow()

    private var onSave: (Album) -> Unit = {}
    private var onDelete: (Album) -> Unit = {}
    var showDelete: Boolean = false
        private set

    fun launchEditor(
        navController: NavController,
        album: Album? = null,
        onSave: (Album) -> Unit,
        onDelete: ((Album) -> Unit)? = null
    ) {
        this.onSave = onSave
        this.onDelete = onDelete ?: {}
        this.showDelete = album != null
        _albumState.value = album?.deepCopy() ?: Album(year = LocalDate.now().year)
        navController.navigate(Screen.EDIT_ALBUM.name)
    }

    fun updateName(name: String) {
        _albumState.value = _albumState.value.copy(name = name)
    }

    fun updateArtist(artist: Artist) {
        _albumState.value = _albumState.value.copy(artist = artist)
    }

    fun updateYear(year: Int) {
        _albumState.value = _albumState.value.copy(year = year)
    }

    fun updateSigned(signed: Boolean) {
        _albumState.value = _albumState.value.copy(signed = signed)
    }

    fun onSave() {
        onSave(_albumState.value)
    }

    fun onDelete() {
        onDelete(_albumState.value)
    }
}
