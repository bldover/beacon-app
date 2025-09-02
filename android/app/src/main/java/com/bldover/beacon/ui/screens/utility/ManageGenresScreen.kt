package com.bldover.beacon.ui.screens.utility

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.screens.editor.genre.SearchableGenresList
import timber.log.Timber

@Composable
fun ManageGenresScreen(
    navController: NavController,
    genreManagerViewModel: GenreManagerViewModel = hiltViewModel()
) {
    Timber.d("composing ManageGenresScreen")
    LaunchedEffect(Unit) {
        genreManagerViewModel.resetFilter()
    }
    val genreState by genreManagerViewModel.uiState.collectAsState()

    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Manage Genres",
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        SearchableGenresList(
            genreState = genreState,
            onSearchGenres = genreManagerViewModel::applyFilter
        )
    }
}