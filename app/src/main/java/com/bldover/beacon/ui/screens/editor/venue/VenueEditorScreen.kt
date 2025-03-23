package com.bldover.beacon.ui.screens.editor.venue

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
        Column(
            modifier = Modifier.fillMaxWidth(),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            val venue by venueEditorViewModel.venueState.collectAsState()
            TextField(
                value = venue.name,
                onValueChange = venueEditorViewModel::updateName,
                label = { Text("Name") },
                modifier = Modifier.fillMaxWidth()
            )
            TextField(
                value = venue.city,
                onValueChange = venueEditorViewModel::updateCity,
                label = { Text("City") },
                modifier = Modifier.fillMaxWidth()
            )
            TextField(
                value = venue.state,
                onValueChange = venueEditorViewModel::updateState,
                label = { Text("State") },
                modifier = Modifier.fillMaxWidth()
            )
            Row(
                horizontalArrangement = Arrangement.End,
                modifier = Modifier.fillMaxWidth()
            ) {
                SaveCancelButtons(
                    onSave = {
                        venueEditorViewModel.onSave()
                    },
                    onCancel = { navController.popBackStack() }
                )
            }
        }
    }
}