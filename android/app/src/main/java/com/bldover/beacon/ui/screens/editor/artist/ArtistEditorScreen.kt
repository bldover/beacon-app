package com.bldover.beacon.ui.screens.editor.artist

import androidx.compose.foundation.lazy.items
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.AddNewCard
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TextEntryCard
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.editor.SaveableEditFieldsColumn
import com.bldover.beacon.ui.components.editor.SwipeableGenreCard
import com.bldover.beacon.ui.screens.editor.genre.GenreSelectorViewModel

@Composable
fun ArtistEditorScreen(
    navController: NavController,
    artistEditorViewModel: ArtistEditorViewModel,
    genreSelectorViewModel: GenreSelectorViewModel = hiltViewModel()
) {
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Edit Artist",
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        val artist by artistEditorViewModel.artistState.collectAsState()

        SaveableEditFieldsColumn(
            onSave = { artistEditorViewModel.onSave() },
            onCancel = { navController.popBackStack() }
        ) {
            item {
                TextEntryCard(
                    label = "Name",
                    value = artist.name,
                    dialogTitle = "Edit Name",
                    dialogLabel = "Artist Name",
                    onValueChange = artistEditorViewModel::updateName
                )
            }
            
            val currentGenres = artist.genres.user
            items(items = currentGenres) { genre ->
                SwipeableGenreCard(
                    genre = genre,
                    onSwipe = artistEditorViewModel::removeGenre
                )
            }
            
            item {
                AddNewCard(
                    label = "Add Genre",
                    onClick = {
                        genreSelectorViewModel.launchSelector(navController, artist) {
                            artistEditorViewModel.addGenre(it)
                        }
                    }
                )
            }
        }
    }
}