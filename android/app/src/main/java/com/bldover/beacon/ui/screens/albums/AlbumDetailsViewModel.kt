package com.bldover.beacon.ui.screens.albums

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.album.Album
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import javax.inject.Inject

@HiltViewModel
class AlbumDetailsViewModel @Inject constructor() : ViewModel() {

    private val _albumState = MutableStateFlow(Album())
    val albumState = _albumState.asStateFlow()

    fun launchDetails(navController: NavController, album: Album) {
        _albumState.value = album.deepCopy()
        navController.navigate(Screen.ALBUM_DETAILS.name)
    }

    fun updateAlbum(album: Album) {
        _albumState.value = album.deepCopy()
    }
}
