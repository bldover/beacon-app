package com.bldover.beacon.ui.screens.editor.artist

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.AddNewCard
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.editor.ArtistNameDialogEditCard
import com.bldover.beacon.ui.components.editor.SaveCancelButtons
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
        LazyColumn(
            modifier = Modifier.fillMaxWidth(),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            item {
                ArtistNameDialogEditCard(
                    artist = artist,
                    onValueChange = artistEditorViewModel::updateName,
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
            
            item {
                Row(
                    horizontalArrangement = Arrangement.End,
                    modifier = Modifier.fillMaxWidth()
                ) {
                    SaveCancelButtons(
                        onSave = {
                            artistEditorViewModel.onSave()
                        },
                        onCancel = { navController.popBackStack() }
                    )
                }
            }
        }
    }
}