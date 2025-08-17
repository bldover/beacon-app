package com.bldover.beacon.ui.screens.editor.artist

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.util.toCommaSeparatedString
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.editor.SaveCancelButtons

@Composable
fun ArtistEditorScreen(
    navController: NavController,
    artistEditorViewModel: ArtistEditorViewModel
) {
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Edit Artist",
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        Column(
            modifier = Modifier.fillMaxWidth(),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            val artist by artistEditorViewModel.artistState.collectAsState()
            TextField(
                value = artist.name,
                onValueChange = artistEditorViewModel::updateName,
                label = { Text("Name") },
                modifier = Modifier.fillMaxWidth()
            )
            TextField(
                value = toCommaSeparatedString(artist.genres.user),
                onValueChange = artistEditorViewModel::updateUserGenres,
                label = { Text("User Genres") },
                modifier = Modifier.fillMaxWidth()
            )
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