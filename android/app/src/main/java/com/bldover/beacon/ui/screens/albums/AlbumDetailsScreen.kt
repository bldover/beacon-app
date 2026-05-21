package com.bldover.beacon.ui.screens.albums

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ExperimentalLayoutApi
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.aspectRatio
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.outlined.Album
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedCard
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import coil.compose.AsyncImage
import com.bldover.beacon.data.model.album.Album
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar

private val SignedAccent = Color(0xFFD4AF37)
private val WishlistedAccent = Color(0xFF4CAF50)

@Composable
fun AlbumDetailsScreen(
    navController: NavController,
    albumDetailsViewModel: AlbumDetailsViewModel,
    onEdit: (Album) -> Unit
) {
    val album by albumDetailsViewModel.albumState.collectAsState()

    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Album Details",
                leadingIcon = { BackButton(navController = navController) },
                trailingIcon = {
                    IconButton(onClick = { onEdit(album) }) {
                        Icon(
                            imageVector = Icons.Filled.Edit,
                            contentDescription = "Edit album"
                        )
                    }
                }
            )
        }
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .verticalScroll(rememberScrollState()),
            verticalArrangement = Arrangement.spacedBy(24.dp)
        ) {
            CoverHeader(name = album.name, uri = album.coverImageUri)
            ArtistChips(artists = album.artists)
            AttributeRow(signed = album.signed, wishlisted = album.wishlisted)
            DetailsSection(album = album)
            Spacer(modifier = Modifier.height(16.dp))
        }
    }
}

@Composable
private fun CoverHeader(name: String, uri: String?) {
    Column(
        modifier = Modifier.fillMaxWidth(),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        Surface(
            modifier = Modifier
                .fillMaxWidth(0.7f)
                .aspectRatio(1f),
            shape = RoundedCornerShape(16.dp),
            color = MaterialTheme.colorScheme.surfaceContainerHigh,
            tonalElevation = 2.dp
        ) {
            if (uri != null) {
                AsyncImage(
                    model = uri,
                    contentDescription = "Album cover",
                    modifier = Modifier
                        .fillMaxWidth()
                        .clip(RoundedCornerShape(16.dp))
                )
            } else {
                Box(
                    modifier = Modifier.fillMaxWidth(),
                    contentAlignment = Alignment.Center
                ) {
                    Icon(
                        imageVector = Icons.Outlined.Album,
                        contentDescription = "No cover image",
                        modifier = Modifier.fillMaxWidth(0.4f).aspectRatio(1f),
                        tint = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }
        }
        if (name.isNotBlank()) {
            Text(
                text = name,
                style = MaterialTheme.typography.titleLarge,
                fontWeight = FontWeight.SemiBold
            )
        }
    }
}

@OptIn(ExperimentalLayoutApi::class)
@Composable
private fun ArtistChips(artists: List<Artist>) {
    if (artists.isEmpty()) return
    FlowRow(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        artists.forEach { artist ->
            OutlinedCard(
                shape = RoundedCornerShape(50),
                border = BorderStroke(1.dp, MaterialTheme.colorScheme.outline)
            ) {
                Text(
                    text = artist.name,
                    style = MaterialTheme.typography.bodyMedium,
                    modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
                )
            }
        }
    }
}

@OptIn(ExperimentalLayoutApi::class)
@Composable
private fun AttributeRow(signed: Boolean, wishlisted: Boolean) {
    if (!signed && !wishlisted) return
    FlowRow(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        if (signed) {
            AttributeChip(label = "Signed", color = SignedAccent)
        }
        if (wishlisted) {
            AttributeChip(label = "Wishlisted", color = WishlistedAccent)
        }
    }
}

@Composable
private fun AttributeChip(label: String, color: Color) {
    OutlinedCard(
        shape = RoundedCornerShape(50),
        border = BorderStroke(1.5.dp, color)
    ) {
        Text(
            text = label,
            style = MaterialTheme.typography.labelLarge,
            color = color,
            fontWeight = FontWeight.SemiBold,
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
        )
    }
}

@Composable
private fun DetailsSection(album: Album) {
    val rows = buildList {
        add("Format" to formatLine(album))
        add("Released" to album.year.toString())
        if (album.genre.isNotBlank()) add("Genre" to album.genre)
        if (album.notes.isNotBlank()) add("Notes" to album.notes)
    }
    val labelWidth = 96.dp
    Column(
        modifier = Modifier.fillMaxWidth(),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        rows.forEach { (label, value) ->
            DetailRow(label = label, value = value, labelWidth = labelWidth)
        }
    }
}

@Composable
private fun DetailRow(
    label: String,
    value: String,
    labelWidth: androidx.compose.ui.unit.Dp
) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        verticalAlignment = Alignment.Top
    ) {
        Text(
            text = label,
            style = MaterialTheme.typography.bodyMedium,
            fontWeight = FontWeight.SemiBold,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            modifier = Modifier.width(labelWidth)
        )
        Spacer(modifier = Modifier.width(12.dp))
        Text(
            text = value,
            style = MaterialTheme.typography.bodyMedium,
            modifier = Modifier.fillMaxWidth()
        )
    }
}

private fun formatLine(album: Album): String {
    val parts = buildList {
        add(album.format.displayName)
        if (album.limitedEdition) add("Limited Edition")
        if (album.variant.isNotBlank()) add(album.variant)
    }
    return parts.joinToString(" · ")
}
