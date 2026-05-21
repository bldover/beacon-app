package com.bldover.beacon.ui.components.editor

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.DismissableCard

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
fun SwipeableArtistEditCard(
    artist: Artist,
    label: String,
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
            SummaryLine(label = label) {
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