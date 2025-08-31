package com.bldover.beacon.ui.components.editor

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AddCircle
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.ui.components.common.AddNewCard
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.DismissableCard
import com.bldover.beacon.ui.components.common.TextEntryDialog

@Composable
fun GenreCard(
    genre: String,
    onClick: () -> Unit,
    hasAccentBorder: Boolean = false
) {
    BasicCard(
        modifier = Modifier.clickable { onClick() },
        border = if (hasAccentBorder) BorderStroke(width = 1.dp, color = MaterialTheme.colorScheme.primary) else null
    ) {
        Text(
            text = genre,
            style = MaterialTheme.typography.bodyMedium
        )
    }
}

@Composable
fun SwipeableGenreCard(
    genre: String,
    onSwipe: (String) -> Unit,
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
            onDismiss = { onSwipe(genre) },
        ) {
            SummaryLine(label = "Genre") {
                Text(
                    text = genre,
                    style = MaterialTheme.typography.bodyMedium
                )
            }
        }
    }
}

@Composable
fun NewGenreDialogEditCard(
    onNewGenre: (String) -> Unit
) {
    var isVisible by remember { mutableStateOf(false) }

    AddNewCard(
        label = "Genre",
        onClick = { isVisible = true }
    )

    TextEntryDialog(
        isVisible = isVisible,
        title = "New Genre",
        label = "Genre Name",
        onDismiss = { isVisible = false },
        onSave = {
            onNewGenre(it)
            isVisible = false
        }
    )
}