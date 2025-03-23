package com.bldover.beacon.ui.screens.editor.artist

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
fun ArtistSelectorScreen(
    navController: NavController,
    artistSelectorViewModel: ArtistSelectorViewModel,
    artistEditorViewModel: ArtistEditorViewModel,
    artistsViewModel: ArtistsViewModel = hiltViewModel()
) {
    Timber.d("composing ArtistSelectorScreen")
    LaunchedEffect(Unit) {
        artistsViewModel.resetFilter()
    }
    val artistState by artistsViewModel.uiState.collectAsState()

    ScreenFrame(
        topBar = { TitleTopBar(
            title = "Select Artist",
            leadingIcon = { BackButton(navController = navController) }
        ) }
    ) {
        SearchableArtistsList(
            artistState = artistState,
            onSearchArtists = artistsViewModel::applyFilter,
            onArtistSelected = {
                artistSelectorViewModel.selectArtist(it)
                navController.popBackStack()
            },
            onNewArtist = {
                artistEditorViewModel.launchEditor(
                    navController = navController,
                    onSave = {
                        artistSelectorViewModel.selectArtist(it)
                        navController.popBackStack()
                        navController.popBackStack()
                    }
                )
            }
        )
    }
}

