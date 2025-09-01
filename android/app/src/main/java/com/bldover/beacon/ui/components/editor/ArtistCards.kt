package com.bldover.beacon.ui.components.editor

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.compose.ui.text.style.TextDecoration
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.artist.ArtistType
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.DismissableCard
import com.bldover.beacon.ui.components.common.TextEntryDialog

@Composable
fun ArtistDetailsCard(
    artist: Artist,
    showAnyGenre: Boolean = false,
    onClick: () -> Unit = {}
) {
    BasicCard(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
    ) {
        Text(
            text = artist.name,
            style = MaterialTheme.typography.bodyLarge
        )
        Text(
            text = if (showAnyGenre) artist.genres.getTopGenre().orEmpty() else artist.genres.user.firstOrNull().orEmpty(),
            style = MaterialTheme.typography.bodySmall,
            color = if (artist.genres.hasUserGenre()) MaterialTheme.colorScheme.onSurface else MaterialTheme.colorScheme.onErrorContainer
        )
    }
}

@Composable
fun ArtistNameDialogEditCard(
    artist: Artist,
    onValueChange: (String) -> Unit
) {
    var showDialog by remember { mutableStateOf(false) }
    BasicCard(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = { showDialog = true })
    ) {
        SummaryLine(label = "Name") {
            Text(text = artist.name)
        }
    }

    TextEntryDialog(
        isVisible = showDialog,
        title = "Edit Name",
        label = "Artist Name",
        initialValue = artist.name,
        onDismiss = { showDialog = false },
        onSave = {
            onValueChange(it)
            showDialog = false
        }
    )
}

@Composable
fun SwipeableArtistEditCard(
    artist: Artist,
    artistType: ArtistType,
    onSwipe: (Artist) -> Unit,
    onClick: (() -> Unit)? = null,
) {
    Box(
        modifier = if (onClick == null) {
            Modifier
        } else {
            Modifier.clickable(onClick = onClick)
        }
    ) {
        DismissableCard(
            onDismiss = { onSwipe(artist) },
            border = if (artist.id.primary != null) BorderStroke(width = 1.dp, color = MaterialTheme.colorScheme.primary) else null,
        ) {
            SummaryLine(label = artistType.label) {
                Text(
                    text = artist.name
                )
                Text(
                    text = artist.genres.getTopGenre() ?: "",
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurface,
                    textDecoration = if (artist.genres.hasUserGenre()) TextDecoration.None else TextDecoration.Underline
                )
            }
        }
    }
}