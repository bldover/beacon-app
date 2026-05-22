package com.bldover.beacon.ui.screens.editor.album

import android.Manifest
import android.content.Context
import android.content.Intent
import android.content.pm.PackageManager
import android.net.Uri
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.PickVisualMediaRequest
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.core.content.ContextCompat
import androidx.core.content.FileProvider
import androidx.navigation.NavController
import com.bldover.beacon.data.model.album.AlbumFormat
import com.bldover.beacon.ui.components.common.AddNewCard
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.RadioSelectorDialog
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TextEntryCard
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.common.YearPickerDialog
import com.bldover.beacon.ui.components.editor.ReducedMinSizeSwitch
import com.bldover.beacon.ui.components.editor.SaveableEditFieldsColumn
import com.bldover.beacon.ui.components.editor.SummaryCard
import com.bldover.beacon.ui.components.editor.SwipeableArtistEditCard
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel
import com.bldover.beacon.ui.screens.editor.genre.GenreSelectorViewModel
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
    artistSelectorViewModel: ArtistSelectorViewModel,
    genreSelectorViewModel: GenreSelectorViewModel
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

        var showYearPicker by remember { mutableStateOf(false) }
        var showFormatPicker by remember { mutableStateOf(false) }
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

        SaveableEditFieldsColumn (
            onSave = { albumEditorViewModel.onSave() },
            onCancel = { navController.popBackStack() },
            showDelete = albumEditorViewModel.showDelete,
            onDelete = { albumEditorViewModel.onDelete() }
        ) {
            item {
                TextEntryCard(
                    label = "Album",
                    value = album.name,
                    dialogTitle = "Edit Name",
                    dialogLabel = "Album Name",
                    onValueChange = albumEditorViewModel::updateName
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
                SummaryCard(
                    label = "Year",
                    onClick = { showYearPicker = true }
                ) {
                    Text(text = if (album.year == 0) "" else album.year.toString())
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
                SummaryCard(
                    label = "Format",
                    onClick = { showFormatPicker = true }
                ) {
                    Text(text = album.format.displayName)
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
                ReducedMinSizeSwitch(
                    label = "Limited Edition",
                    checked = album.limitedEdition,
                    onChange = albumEditorViewModel::updateLimitedEdition
                )
            }

            item {
                TextEntryCard(
                    label = "Variant",
                    value = album.variant,
                    dialogTitle = "Edit Variant",
                    dialogLabel = "Variant",
                    onValueChange = albumEditorViewModel::updateVariant
                )
            }

            item {
                SummaryCard(
                    label = "Genre",
                    onClick = {
                        genreSelectorViewModel.launchSelector(
                            navController = navController,
                            artist = null,
                            onSelect = { albumEditorViewModel.updateGenre(it) }
                        )
                    }
                ) {
                    Text(text = album.genre)
                }
            }

            item {
                SummaryCard(
                    label = "Cover Image",
                    onClick = { showImageSourceDialog = true }
                ) {
                    Text(text = if (album.coverImageUri != null) "Image selected" else "None")
                }
            }

            item {
                ReducedMinSizeSwitch(
                    label = "Signed",
                    checked = album.signed,
                    onChange = albumEditorViewModel::updateSigned
                )
            }

            item {
                ReducedMinSizeSwitch(
                    label = "Wishlisted",
                    checked = album.wishlisted,
                    onChange = albumEditorViewModel::updateWishlisted
                )
            }

            item {
                TextEntryCard(
                    label = "Notes",
                    value = album.notes,
                    dialogTitle = "Edit Notes",
                    dialogLabel = "Notes",
                    onValueChange = albumEditorViewModel::updateNotes
                )
            }
        }
    }
}
