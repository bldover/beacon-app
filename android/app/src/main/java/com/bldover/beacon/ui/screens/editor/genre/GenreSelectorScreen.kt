package com.bldover.beacon.ui.screens.editor.genre

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import timber.log.Timber

@Composable
fun GenreSelectorScreen(
    navController: NavController,
    genreSelectorViewModel: GenreSelectorViewModel,
    genresViewModel: GenresViewModel = hiltViewModel()
) {
    Timber.d("composing GenreSelectorScreen")
    LaunchedEffect(Unit) {
        genresViewModel.resetFilter()
    }
    val genreState by genresViewModel.uiState.collectAsState()

    ScreenFrame(
        topBar = { TitleTopBar(
            title = "Select Genre",
            leadingIcon = { BackButton(navController = navController) }
        ) }
    ) {
        SearchableGenresList(
            genreState = genreState,
            onSearchGenres = genresViewModel::applyFilter,
            onGenreSelected = {
                genreSelectorViewModel.selectGenre(it)
                navController.popBackStack()
            },
            onCustomGenre = {
                genreSelectorViewModel.selectGenre(it)
                navController.popBackStack()
            }
        )
    }
}