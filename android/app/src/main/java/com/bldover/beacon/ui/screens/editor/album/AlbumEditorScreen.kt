package com.bldover.beacon.ui.screens.editor.album

import android.Manifest
import android.content.Context
import android.content.Intent
import android.content.pm.PackageManager
import android.net.Uri
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.PickVisualMediaRequest
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
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
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.core.content.ContextCompat
import androidx.core.content.FileProvider
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.AddNewCard
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.data.model.album.AlbumFormat
import com.bldover.beacon.ui.components.common.RadioSelectorDialog
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TextEntryDialog
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.common.YearPickerDialog
import com.bldover.beacon.ui.components.editor.DeleteButton
import com.bldover.beacon.ui.components.editor.SaveCancelButtons
import com.bldover.beacon.ui.components.editor.SummaryLine
import com.bldover.beacon.ui.components.editor.SwipeableArtistEditCard
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel
import java.io.File

private fun createCameraUri(context: Context): Uri {
    val imagesDir = File(context.filesDir, "cover_images")
    imagesDir.mkdirs()
    val imageFile = File.createTempFile("cover_", ".jpg", imagesDir)
    return FileProvider.getUriForFile(context, "${context.packageName}.fileprovider", imageFile)
}

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
        val context = LocalContext.current

        var showNameDialog by remember { mutableStateOf(false) }
        var showYearPicker by remember { mutableStateOf(false) }
        var showFormatPicker by remember { mutableStateOf(false) }
        var showVariantDialog by remember { mutableStateOf(false) }
        var showNotesDialog by remember { mutableStateOf(false) }
        var showImageSourceDialog by remember { mutableStateOf(false) }

        val cameraImageUri = remember { mutableStateOf<Uri?>(null) }

        val cameraLauncher = rememberLauncherForActivityResult(
            contract = ActivityResultContracts.TakePicture()
        ) { success ->
            if (success) {
                cameraImageUri.value?.let { albumEditorViewModel.uploadCoverImage(it.toString()) }
            }
        }

        val galleryLauncher = rememberLauncherForActivityResult(
            contract = ActivityResultContracts.PickVisualMedia()
        ) { uri ->
            uri?.let {
                runCatching {
                    context.contentResolver.takePersistableUriPermission(
                        it, Intent.FLAG_GRANT_READ_URI_PERMISSION
                    )
                }
                albumEditorViewModel.uploadCoverImage(it.toString())
            }
        }

        val cameraPermissionLauncher = rememberLauncherForActivityResult(
            contract = ActivityResultContracts.RequestPermission()
        ) { granted ->
            if (granted) {
                val uri = createCameraUri(context)
                cameraImageUri.value = uri
                cameraLauncher.launch(uri)
            }
        }

        if (showImageSourceDialog) {
            AlertDialog(
                onDismissRequest = { showImageSourceDialog = false },
                title = { Text("Cover Image") },
                text = {
                    Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                        Button(
                            modifier = Modifier.fillMaxWidth(),
                            onClick = {
                                showImageSourceDialog = false
                                if (ContextCompat.checkSelfPermission(context, Manifest.permission.CAMERA)
                                    == PackageManager.PERMISSION_GRANTED
                                ) {
                                    val uri = createCameraUri(context)
                                    cameraImageUri.value = uri
                                    cameraLauncher.launch(uri)
                                } else {
                                    cameraPermissionLauncher.launch(Manifest.permission.CAMERA)
                                }
                            }
                        ) { Text("Take Photo") }
                        Button(
                            modifier = Modifier.fillMaxWidth(),
                            onClick = {
                                showImageSourceDialog = false
                                galleryLauncher.launch(
                                    PickVisualMediaRequest(ActivityResultContracts.PickVisualMedia.ImageOnly)
                                )
                            }
                        ) { Text("Choose from Gallery") }
                        Button(
                            modifier = Modifier.fillMaxWidth(),
                            onClick = {
                                showImageSourceDialog = false
                                albumEditorViewModel.clearCoverImage()
                            }
                        ) { Text("Clear Image") }
                    }
                },
                confirmButton = {},
                dismissButton = {
                    Button(onClick = { showImageSourceDialog = false }) { Text("Cancel") }
                }
            )
        }

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

            items(items = album.artists, key = { it.name }) { artist ->
                SwipeableArtistEditCard(
                    artist = artist,
                    label = "Artist",
                    onSwipe = albumEditorViewModel::removeArtist,
                    onClick = {
                        artistSelectorViewModel.launchSelector(navController) {
                            albumEditorViewModel.replaceArtist(artist, it)
                        }
                    }
                )
            }

            item {
                AddNewCard(
                    label = "Add Artist",
                    onClick = {
                        artistSelectorViewModel.launchSelector(navController) {
                            albumEditorViewModel.addArtist(it)
                        }
                    }
                )
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
                BasicCard(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { showFormatPicker = true }
                ) {
                    SummaryLine(label = "Format") {
                        Text(text = album.format.displayName)
                    }
                }
                RadioSelectorDialog(
                    isVisible = showFormatPicker,
                    title = "Select Format",
                    options = AlbumFormat.entries,
                    selectedOption = album.format,
                    getLabel = { it.displayName },
                    onDismiss = { showFormatPicker = false },
                    onOptionSelected = {
                        albumEditorViewModel.updateFormat(it)
                        showFormatPicker = false
                    }
                )
            }

            item {
                BasicCard(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { showVariantDialog = true }
                ) {
                    SummaryLine(label = "Variant") {
                        Text(text = album.variant)
                    }
                }
                TextEntryDialog(
                    isVisible = showVariantDialog,
                    title = "Edit Variant",
                    label = "Variant",
                    initialValue = album.variant,
                    onDismiss = { showVariantDialog = false },
                    onSave = {
                        albumEditorViewModel.updateVariant(it)
                        showVariantDialog = false
                    }
                )
            }

            item {
                BasicCard(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { showNotesDialog = true }
                ) {
                    SummaryLine(label = "Notes") {
                        Text(text = album.notes)
                    }
                }
                TextEntryDialog(
                    isVisible = showNotesDialog,
                    title = "Edit Notes",
                    label = "Notes",
                    initialValue = album.notes,
                    onDismiss = { showNotesDialog = false },
                    onSave = {
                        albumEditorViewModel.updateNotes(it)
                        showNotesDialog = false
                    }
                )
            }

            item {
                BasicCard(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { showImageSourceDialog = true }
                ) {
                    SummaryLine(label = "Cover Image") {
                        Text(text = if (album.coverImageUri != null) "Image selected" else "None")
                    }
                }
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
                BasicCard {
                    SummaryLine(label = "Wishlisted") {
                        Switch(
                            checked = album.wishlisted,
                            onCheckedChange = albumEditorViewModel::updateWishlisted
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
