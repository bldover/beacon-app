package com.bldover.beacon.ui.components.editor

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AddCircle
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicOutlinedCard
import com.bldover.beacon.ui.components.common.DismissableCard
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel

@Composable
fun ArtistDetailsCard(
    artist: Artist,
    showAllGenres: Boolean = false,
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
            text = if (showAllGenres) artist.genres.getGenres().joinToString(", ") else artist.genres.getTopGenre() ?: "",
            style = MaterialTheme.typography.bodySmall
        )
    }
}

enum class ArtistType(val label: String) {
    HEADLINER("Headliner"),
    OPENER("Opener");
}

@Composable
fun SwipeableArtistEditCard(
    artist: Artist,
    artistType: ArtistType,
    onSwipe: (Artist) -> Unit,
    onSelect: (() -> Unit)? = null,
) {
    Box(
        modifier = if (onSelect == null) {
            Modifier
        } else {
            Modifier.clickable(onClick = onSelect)
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
                    color = MaterialTheme.colorScheme.onSurface
                )
            }
        }
    }
}

@Composable
fun AddArtistCard(
    artistType: ArtistType,
    onSelect: (Artist) -> Unit,
    navController: NavController,
    artistSelectorViewModel: ArtistSelectorViewModel,
) {
    Box(
        modifier = Modifier.clickable {
            artistSelectorViewModel.launchSelector(navController) {
                onSelect(it)
            }
        }
    ) {
        BasicOutlinedCard {
            SummaryLine(label = artistType.label) {
                Icon(
                    imageVector = Icons.Default.AddCircle,
                    contentDescription = "Add ${artistType.label}"
                )
            }
        }
    }
}