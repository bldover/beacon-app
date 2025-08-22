package com.bldover.beacon.ui.screens.editor.artist

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.artist.Artist
import dagger.hilt.android.lifecycle.HiltViewModel
import timber.log.Timber
import javax.inject.Inject

@HiltViewModel
class ArtistSelectorViewModel @Inject constructor() : ViewModel() {

    private var onSelect: (Artist) -> Unit = {}

    fun launchSelector(
        navController: NavController,
        onSelect: (Artist) -> Unit
    ) {
        Timber.d("launching artist selector")
        this.onSelect = onSelect
        navController.navigate(Screen.SELECT_ARTIST.name)
    }

    fun selectArtist(artist: Artist) {
        Timber.d("artist selector - selected artist $artist")
        onSelect(artist)
    }
}