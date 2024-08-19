package com.bldover.beacon.ui.screens.editor

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.SheetState
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.Artist
import com.bldover.beacon.data.model.Event
import com.bldover.beacon.data.model.Venue
import com.bldover.beacon.ui.components.BasicCard
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EventEditBottomSheet(
    event: Event,
    modalSheetState: SheetState,
    onClosed: () -> Unit,
) {
    val scope = rememberCoroutineScope()
    ModalBottomSheet(
        onDismissRequest = { onClosed() },
        sheetState = modalSheetState
    ) {
        Column(modifier = Modifier.padding(start = 16.dp, end = 16.dp, bottom = 16.dp)) {
            EditSheet(
                event = event,
                onCancel = {
                    hideSheet(scope = scope, modalSheetState = modalSheetState, onClosed = onClosed)
                },
                onSave = {
                    // TODO: save event it
                    hideSheet(scope = scope, modalSheetState = modalSheetState, onClosed = onClosed)
                }
            )
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
fun hideSheet(
    scope: CoroutineScope,
    modalSheetState: SheetState,
    onClosed: () -> Unit
) {
    scope.launch {
        modalSheetState.hide()
    }.invokeOnCompletion {
        if (!modalSheetState.isVisible)
            onClosed()
    }
}

@OptIn(ExperimentalMaterial3Api::class)
fun showSheet(
    scope: CoroutineScope,
    modalSheetState: SheetState
) {
    scope.launch {
        modalSheetState.show()
    }
}

private enum class EditField {
    DATE, VENUE, ARTISTS, NONE
}

@Composable
fun EditSheet(
    event: Event,
    onCancel: () -> Unit = {},
    onSave: (Event) -> Unit = {}
) {
    var selected by remember { mutableStateOf(EditField.NONE) }
    val artistsState = remember { mutableStateOf(event.artists) }
    val dateState = remember { mutableStateOf(event.date) }
    val venueState = remember { mutableStateOf(event.venue) }
    Column(modifier = Modifier.fillMaxWidth()) {
        if (showField(EditField.ARTISTS, selected)) {
            ArtistsEditCard(
                artistsState = artistsState,
                onClick = { selected = EditField.ARTISTS },
                onClose = { selected = EditField.NONE }
            )
            Spacer(modifier = Modifier.height(16.dp))
        }
        if (showField(EditField.DATE, selected)) {
            DateEditCard(
                dateState = dateState,
                onClick = { selected = EditField.DATE },
                onClose = { selected = EditField.NONE }
            )
            Spacer(modifier = Modifier.height(16.dp))
        }
        if (showField(EditField.VENUE, selected)) {
            VenueEditCard(
                venueState = venueState,
                onClick = { selected = EditField.VENUE },
                onClose = { selected = EditField.NONE }
            )
        }
    }
    if (selected == EditField.NONE) {
        SaveCancelButtons(
            onCancel = onCancel,
            onSave = {
                // todo: get updated event
                onSave(event)
            }
        )
    }
}

private fun showField(field: EditField, selected: EditField): Boolean {
    return field == selected || selected == EditField.NONE
}

@Composable
fun EditVenue(
    venue: Venue,
    onClick: () -> Unit,
    onClose: () -> Unit,
    onSave: () -> Unit
) {
    BasicCard {
        Text(text = "Venue")
        Text(text = venue.name)
    }
}

@Composable
fun EditArtists(
    artists: List<Artist>,
    onClick: () -> Unit,
    onClose: () -> Unit,
    onSave: () -> Unit
) {
    BasicCard {
        Column {
            Text(text = "Artists")
            for (artist in artists) {
                Text(text = artist.name)
            }
        }
    }
}