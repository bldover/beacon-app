package com.bldover.beacon.ui.screens.utility

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.screens.editor.artist.ArtistEditorViewModel
import com.bldover.beacon.ui.screens.editor.artist.ArtistsViewModel
import com.bldover.beacon.ui.screens.editor.artist.SearchableArtistsList
import timber.log.Timber

@Composable
fun ManageArtistsScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    artistEditorViewModel: ArtistEditorViewModel,
    artistsViewModel: ArtistsViewModel = hiltViewModel()
) {
    Timber.d("composing ManageArtistsScreen")
    LaunchedEffect(Unit) {
        artistsViewModel.resetFilter()
    }
    val artistState by artistsViewModel.uiState.collectAsState()

    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Manage Artists",
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        SearchableArtistsList(
            artistState = artistState,
            showAllGenres = true,
            orderingState = artistsViewModel.ordering.collectAsState().value,
            onSearchArtists = artistsViewModel::applyFilter,
            onOrderingChange = artistsViewModel::applyOrdering,
            onArtistSelected = {
                artistEditorViewModel.launchEditor(
                    navController = navController,
                    artist = it,
                    onSave = { updated ->
                        artistsViewModel.updateArtist(
                            artist = updated,
                            onSuccess = {
                                navController.popBackStack()
                            },
                            onError = { err ->
                                Timber.e(err)
                                snackbarState.showSnackbar("Failed to save artist")
                            }
                        )
                    }
                )
            },
            onNewArtist = {
                artistEditorViewModel.launchEditor(
                    navController = navController,
                    onSave = {
                        artistsViewModel.addArtist(
                            artist = it,
                            onSuccess = {
                                navController.popBackStack()
                            },
                            onError = { err ->
                                Timber.e(err)
                                snackbarState.showSnackbar("Failed to save artist")
                            }
                        )
                    }
                )
            }
        )
    }
}