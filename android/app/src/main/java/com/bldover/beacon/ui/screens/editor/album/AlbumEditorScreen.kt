package com.bldover.beacon.ui.screens.editor.album

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.Switch
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.AddNewCard
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TextEntryDialog
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.common.YearPickerDialog
import com.bldover.beacon.ui.components.editor.DeleteButton
import com.bldover.beacon.ui.components.editor.SaveCancelButtons
import com.bldover.beacon.ui.components.editor.SummaryLine
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel

@Composable
fun AlbumEditorScreen(
    navController: NavController,
    albumEditorViewModel: AlbumEditorViewModel,
    artistSelectorViewModel: ArtistSelectorViewModel
) {
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Edit Album",
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        val album by albumEditorViewModel.albumState.collectAsState()
        var showNameDialog by remember { mutableStateOf(false) }
        var showYearPicker by remember { mutableStateOf(false) }

        LazyColumn(
            modifier = Modifier.fillMaxWidth(),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            item {
                BasicCard(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { showNameDialog = true }
                ) {
                    SummaryLine(label = "Name") {
                        Text(text = album.name)
                    }
                }
                TextEntryDialog(
                    isVisible = showNameDialog,
                    title = "Edit Name",
                    label = "Album Name",
                    initialValue = album.name,
                    onDismiss = { showNameDialog = false },
                    onSave = {
                        albumEditorViewModel.updateName(it)
                        showNameDialog = false
                    }
                )
            }

            item {
                if (album.artist.isPopulated()) {
                    BasicCard(modifier = Modifier.clickable {
                        artistSelectorViewModel.launchSelector(navController) {
                            albumEditorViewModel.updateArtist(it)
                        }
                    }) {
                        SummaryLine(label = "Artist") {
                            Text(text = album.artist.name)
                        }
                    }
                } else {
                    AddNewCard(
                        label = "Select Artist",
                        onClick = {
                            artistSelectorViewModel.launchSelector(navController) {
                                albumEditorViewModel.updateArtist(it)
                            }
                        }
                    )
                }
            }

            item {
                BasicCard(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { showYearPicker = true }
                ) {
                    SummaryLine(label = "Year") {
                        Text(text = if (album.year == 0) "" else album.year.toString())
                    }
                }
                YearPickerDialog(
                    isVisible = showYearPicker,
                    selectedYear = album.year,
                    onDismiss = { showYearPicker = false },
                    onYearSelected = {
                        albumEditorViewModel.updateYear(it)
                        showYearPicker = false
                    }
                )
            }

            item {
                BasicCard {
                    SummaryLine(label = "Signed") {
                        Switch(
                            checked = album.signed,
                            onCheckedChange = albumEditorViewModel::updateSigned
                        )
                    }
                }
            }

            item {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = 8.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    if (albumEditorViewModel.showDelete) {
                        DeleteButton(onDelete = { albumEditorViewModel.onDelete() })
                    }
                    Spacer(modifier = Modifier.weight(1f))
                    SaveCancelButtons(
                        onCancel = { navController.popBackStack() },
                        onSave = { albumEditorViewModel.onSave() }
                    )
                }
            }
        }
    }
}
