package com.bldover.beacon.ui.screens.editor.genre

import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import timber.log.Timber

@Composable
fun GenreSelectorScreen(
    navController: NavController,
    genreSelectorViewModel: GenreSelectorViewModel
) {
    Timber.d("composing GenreSelectorScreen")
    
    val artist = genreSelectorViewModel.getCurrentArtist()
    val allUserGenres by genreSelectorViewModel.allUserGenres.collectAsState()
    val filteredUserGenres by genreSelectorViewModel.filteredUserGenres.collectAsState()
    val isFiltering by genreSelectorViewModel.isFiltering.collectAsState()

    ScreenFrame(
        topBar = { TitleTopBar(
            title = "Select Genre",
            leadingIcon = { BackButton(navController = navController) }
        ) }
    ) {
        SearchableGenresList(
            artist = artist,
            allUserGenres = allUserGenres,
            filteredUserGenres = filteredUserGenres,
            isFiltering = isFiltering,
            onSearchGenres = genreSelectorViewModel::applyFilter,
            onGenreSelected = {
                genreSelectorViewModel.selectGenre(it)
                navController.popBackStack()
            }
        )
    }
}