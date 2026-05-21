package com.bldover.beacon.ui.screens.editor.venue

import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.editor.SaveableEditFieldsColumn

@Composable
fun VenueEditorScreen(
    navController: NavController,
    venueEditorViewModel: VenueEditorViewModel
) {
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Edit Venue",
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        val venue by venueEditorViewModel.venueState.collectAsState()
        SaveableEditFieldsColumn (
            onSave = { venueEditorViewModel.onSave() },
            onCancel = { navController.popBackStack() }
        ) {
            item {
                TextField(
                    value = venue.name,
                    onValueChange = venueEditorViewModel::updateName,
                    label = { Text("Name") },
                    modifier = Modifier.fillMaxWidth()
                )
            }
            item {
                TextField(
                    value = venue.city,
                    onValueChange = venueEditorViewModel::updateCity,
                    label = { Text("City") },
                    modifier = Modifier.fillMaxWidth()
                )
            }
            item {
                TextField(
                    value = venue.state,
                    onValueChange = venueEditorViewModel::updateState,
                    label = { Text("State") },
                    modifier = Modifier.fillMaxWidth()
                )
            }
        }
    }
}