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
fun DateEditCardOld(
    dateState: MutableState<LocalDate>,
    onClick: () -> Unit,
    onClose: () -> Unit
) {
    var expanded by remember { mutableStateOf(false) }
    val datePickerState = rememberDatePickerState(
        initialSelectedDateMillis = dateState.value.toEpochDay() * 1000 * 60 * 60 * 24
    )
    ExpandableCard(
        onClick = {
            onClick()
            expanded = true
        },
        modifier = Modifier.animateContentSize()
    ) {
        if (!expanded) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Box(
                    modifier = Modifier.weight(0.25f),
                    contentAlignment = Alignment.CenterStart
                ) {
                    Text(text = "Date")
                }
                Box(
                    modifier = Modifier.weight(0.75f),
                    contentAlignment = Alignment.CenterEnd
                ) {
                    Text(text = dateState.value.toString())
                }
            }
        } else {
            DatePicker(state = datePickerState)
            SaveCancelButtons(
                onCancel = {
                    onClose()
                    expanded = false
                },
                onSave = {
                    val selectedDate = LocalDate.ofEpochDay(
                        datePickerState.selectedDateMillis!! / 1000 / 60 / 60 / 24
                    )
                    dateState.value = selectedDate
                    onClose()
                    expanded = false
                }
            )
        }
    }
}

@Composable
fun VenueEditCardOld(
    venueState: MutableState<Venue>,
    onClick: () -> Unit,
    onClose: () -> Unit
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
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Box(
                    modifier = Modifier.weight(0.25f),
                    contentAlignment = Alignment.CenterStart
                ) {
                    Text(text = "Venue")
                }
                Box(
                    modifier = Modifier.weight(0.75f),
                    contentAlignment = Alignment.CenterEnd
                ) {
                    Text(
                        text = venueState.value.name,
                        textAlign = TextAlign.End
                    )
                }
            }
        } else {
            // TODO: Add venue selection
            Box(modifier = Modifier
                .fillMaxWidth()
                .height(60.dp))
            SaveCancelButtons(
                onCancel = {
                    onClose()
                    expanded = false
                },
                onSave = {
                    // TODO: save venue
                    venueState.value = venueState.value
                    onClose()
                    expanded = false
                }
            )
        }
    }
}

@Composable
fun ArtistsEditCardOld(
    artistsState: MutableState<List<Artist>>,
    onClick: () -> Unit,
    onClose: () -> Unit
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
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Box(
                            modifier = Modifier.weight(0.25f),
                            contentAlignment = Alignment.CenterStart
                        ) {
                            Text(text = "Headliner")
                        }
                        Box(
                            modifier = Modifier.weight(0.75f),
                            contentAlignment = Alignment.CenterEnd
                        ) {
                            Text(
                                text = headliner.name,
                                textAlign = TextAlign.End
                            )
                        }
                    }
                    Spacer(modifier = Modifier.height(8.dp))
                    HorizontalDivider(color = MaterialTheme.colorScheme.onPrimaryContainer)
                    Spacer(modifier = Modifier.height(8.dp))
                }
                if (openers.isNotEmpty()) {
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Box(
                            modifier = Modifier.weight(0.25f),
                            contentAlignment = Alignment.CenterStart
                        ) {
                            Text(text = "Openers")
                        }
                        Box(
                            modifier = Modifier.weight(0.75f),
                            contentAlignment = Alignment.CenterEnd
                        ) {
                            Column(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalAlignment = Alignment.End
                            ) {
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
        } else {
            // TODO: Add artists selection
            Box(modifier = Modifier
                .fillMaxWidth()
                .height(60.dp))
            SaveCancelButtons(
                onCancel = {
                    onClose()
                    expanded = false
                },
                onSave = {
                    // TODO: save artists
                    artistsState.value = artistsState.value
                    onClose()
                    expanded = false
                }
            )
        }
    }
}

@Composable
fun SaveCancelButtonsOld(
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