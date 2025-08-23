package com.bldover.beacon.ui.screens.editor.genre

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import dagger.hilt.android.lifecycle.HiltViewModel
import timber.log.Timber
import javax.inject.Inject

@HiltViewModel
class GenreSelectorViewModel @Inject constructor() : ViewModel() {

    private var onSelect: (String) -> Unit = {}

    fun launchSelector(
        navController: NavController,
        onSelect: (String) -> Unit
    ) {
        Timber.d("launching genre selector")
        this.onSelect = onSelect
        navController.navigate(Screen.SELECT_GENRE.name)
    }

    fun selectGenre(genre: String) {
        Timber.d("genre selector - selected genre $genre")
        onSelect(genre)
    }
}