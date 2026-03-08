package com.bldover.beacon.ui.screens.editor.album

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.album.Album
import com.bldover.beacon.data.model.album.AlbumFormat
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

    fun addArtist(artist: Artist) {
        val artists = _albumState.value.artists.toMutableList()
        if (!artists.contains(artist)) {
            artists.add(artist)
        }
        _albumState.value = _albumState.value.copy(artists = artists)
    }

    fun removeArtist(artist: Artist) {
        val artists = _albumState.value.artists.toMutableList()
        artists.remove(artist)
        _albumState.value = _albumState.value.copy(artists = artists)
    }

    fun replaceArtist(old: Artist, new: Artist) {
        val artists = _albumState.value.artists.toMutableList()
        val idx = artists.indexOf(old)
        if (idx >= 0) {
            artists[idx] = new
        }
        _albumState.value = _albumState.value.copy(artists = artists)
    }

    fun updateYear(year: Int) {
        _albumState.value = _albumState.value.copy(year = year)
    }

    fun updateSigned(signed: Boolean) {
        _albumState.value = _albumState.value.copy(signed = signed)
    }

    fun updateWishlisted(wishlisted: Boolean) {
        _albumState.value = _albumState.value.copy(wishlisted = wishlisted)
    }

    fun updateVariant(variant: String) {
        _albumState.value = _albumState.value.copy(variant = variant)
    }

    fun updateFormat(format: AlbumFormat) {
        _albumState.value = _albumState.value.copy(format = format)
    }

    fun updateNotes(notes: String) {
        _albumState.value = _albumState.value.copy(notes = notes)
    }

    fun updateCoverImageUri(uri: String) {
        _albumState.value = _albumState.value.copy(coverImageUri = uri)
    }

    fun onSave() {
        onSave(_albumState.value)
    }

    fun onDelete() {
        onDelete(_albumState.value)
    }
}
