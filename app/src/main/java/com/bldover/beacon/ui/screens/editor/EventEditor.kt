package com.bldover.beacon.ui.screens.editor

import androidx.compose.animation.animateContentSize
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.material3.Button
import androidx.compose.material3.DatePicker
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.rememberDatePickerState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.MutableState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.Artist
import com.bldover.beacon.data.model.Venue
import com.bldover.beacon.ui.components.ExpandableCard
import timber.log.Timber
import java.time.LocalDate

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DateEditCard(
    dateState: MutableState<LocalDate>,
    onClick: () -> Unit,
    onClose: () -> Unit
) {
    val datePickerState = rememberDatePickerState(
        initialSelectedDateMillis = dateState.value.toEpochDay() * 1000 * 60 * 60 * 24
    )
    EditCard(
        onClick = onClick,
        onClose = onClose,
        onSave = {
            val selectedDate = LocalDate.ofEpochDay(
                datePickerState.selectedDateMillis!! / 1000 / 60 / 60 / 24
            )
            dateState.value = selectedDate
        },
        expandContent = {
            DatePicker(state = datePickerState)
        }
    ) {
        SummaryLine(label = "Date") {
            Text(text = dateState.value.toString())
        }
    }
}

@Composable
fun VenueEditCard(
    venueState: MutableState<Venue>,
    onClick: () -> Unit,
    onClose: () -> Unit
) {
    EditCard(
        onClick = onClick,
        onClose = onClose,
        onSave = {
            // TODO: save venue
            venueState.value = venueState.value
        },
        expandContent = {
            // TODO: Add venue selection
            Box(modifier = Modifier
                .fillMaxWidth()
                .height(60.dp))
        }
    ) {
        SummaryLine(label = "Venue") {
            Text(
                text = venueState.value.name,
                textAlign = TextAlign.End
            )
        }
    }
}

@Composable
fun ArtistsEditCard(
    artistsState: MutableState<List<Artist>>,
    onClick: () -> Unit,
    onClose: () -> Unit
) {
    EditCard(
        onClick = onClick,
        onClose = onClose,
        onSave = {
            // TODO: save artists
            artistsState.value = artistsState.value
        },
        expandContent = {
            // TODO: Add artists selection
            Box(modifier = Modifier
                .fillMaxWidth()
                .height(60.dp))
        }
    ) {
        Column(
            modifier = Modifier.fillMaxWidth(),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            val headliner = artistsState.value.find { it.headliner }
            val openers = artistsState.value.filter { !it.headliner }
            Timber.d("Headliner: $headliner")
            Timber.d("Openers: $openers")
            if (headliner != null) {
                SummaryLine(label = "Headliner") {
                    Text(
                        text = headliner.name,
                        textAlign = TextAlign.End
                    )
                }
                Spacer(modifier = Modifier.height(8.dp))
                HorizontalDivider(color = MaterialTheme.colorScheme.onPrimaryContainer)
                Spacer(modifier = Modifier.height(8.dp))
            }
            if (openers.isNotEmpty()) {
                SummaryLine(label = "Openers") {
                    for (opener in openers) {
                        Text(
                            text = opener.name,
                            textAlign = TextAlign.End
                        )
                    }
                }
            }
        }
    }
}

@Composable
fun EditCard(
    onClick: () -> Unit,
    onClose: () -> Unit,
    onSave: () -> Unit,
    expandContent: @Composable () -> Unit,
    content: @Composable () -> Unit
) {
    var expanded by remember { mutableStateOf(false) }
    ExpandableCard(
        onClick = {
            onClick()
            expanded = true
        },
        modifier = Modifier.animateContentSize()
    ) {
        if (!expanded) {
            content()
        } else {
            expandContent()
            SaveCancelButtons(
                onCancel = {
                    onClose()
                    expanded = false
                },
                onSave = {
                    onSave()
                    onClose()
                    expanded = false
                }
            )
        }
    }
}

@Composable
fun SummaryLine(
    label: String,
    content: @Composable () -> Unit
) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically
    ) {
        Box(
            modifier = Modifier.weight(0.25f),
            contentAlignment = Alignment.CenterStart
        ) {
            Text(text = label)
        }
        Box(
            modifier = Modifier.weight(0.75f),
            contentAlignment = Alignment.CenterEnd
        ) {
            Column(
                modifier = Modifier.fillMaxWidth(),
                horizontalAlignment = Alignment.End
            ) {
                content()
            }
        }
    }
}

@Composable
fun SaveCancelButtons(
    onCancel: () -> Unit,
    onSave: () -> Unit
) {
    Row(
        horizontalArrangement = Arrangement.End,
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 8.dp)
    ) {
        Button(onClick = onCancel) {
            Text(text = "Cancel")
        }
        Spacer(modifier = Modifier.width(8.dp))
        Button(onClick = onSave) {
            Text(text = "Save")
        }
    }
}