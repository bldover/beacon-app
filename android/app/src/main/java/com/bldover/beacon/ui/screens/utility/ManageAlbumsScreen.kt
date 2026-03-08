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
import com.bldover.beacon.ui.screens.editor.album.AlbumEditorViewModel
import com.bldover.beacon.ui.screens.editor.album.AlbumsViewModel
import com.bldover.beacon.ui.screens.editor.album.SearchableAlbumsList
import timber.log.Timber

@Composable
fun ManageAlbumsScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    albumEditorViewModel: AlbumEditorViewModel,
    albumsViewModel: AlbumsViewModel = hiltViewModel()
) {
    Timber.d("composing ManageAlbumsScreen")
    LaunchedEffect(Unit) {
        albumsViewModel.resetFilter()
    }
    val albumState by albumsViewModel.uiState.collectAsState()

    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Manage Albums",
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        SearchableAlbumsList(
            albumState = albumState,
            onSearchAlbums = albumsViewModel::applyFilter,
            onAlbumSelected = { album ->
                albumEditorViewModel.launchEditor(
                    navController = navController,
                    album = album,
                    onSave = { updated ->
                        albumsViewModel.updateAlbum(
                            album = updated,
                            onSuccess = { navController.popBackStack() },
                            onError = { msg -> snackbarState.showSnackbar(msg) }
                        )
                    },
                    onDelete = { toDelete ->
                        albumsViewModel.deleteAlbum(
                            album = toDelete,
                            onSuccess = { navController.popBackStack() },
                            onError = { msg -> snackbarState.showSnackbar(msg) }
                        )
                    }
                )
            },
            onNewAlbum = {
                albumEditorViewModel.launchEditor(
                    navController = navController,
                    onSave = { newAlbum ->
                        albumsViewModel.addAlbum(
                            album = newAlbum,
                            onSuccess = { navController.popBackStack() },
                            onError = { msg -> snackbarState.showSnackbar(msg) }
                        )
                    }
                )
            }
        )
    }
}
