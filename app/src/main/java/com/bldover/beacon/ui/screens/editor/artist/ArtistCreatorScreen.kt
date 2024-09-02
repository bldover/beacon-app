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
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.editor.SaveCancelButtons

@Composable
fun ArtistCreatorScreen(
    navController: NavController,
    artistCreatorViewModel: ArtistCreatorViewModel
) {
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "New Artist",
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        Column(
            modifier = Modifier.fillMaxWidth(),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            val artist by artistCreatorViewModel.artistState.collectAsState()
            TextField(
                value = artist.name,
                onValueChange = artistCreatorViewModel::updateName,
                label = { Text("Name") },
                modifier = Modifier.fillMaxWidth()
            )
            TextField(
                value = artist.genre,
                onValueChange = artistCreatorViewModel::updateGenre,
                label = { Text("Genre") },
                modifier = Modifier.fillMaxWidth()
            )
            Row(
                horizontalArrangement = Arrangement.End,
                modifier = Modifier.fillMaxWidth()
            ) {
                SaveCancelButtons(
                    onSave = {
                        artistCreatorViewModel.onSave()
                        navController.popBackStack()
                    },
                    onCancel = { navController.popBackStack() }
                )
            }
        }
    }
}