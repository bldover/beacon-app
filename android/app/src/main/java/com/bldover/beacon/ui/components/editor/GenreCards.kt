package com.bldover.beacon.ui.components.editor

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
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.BasicOutlinedCard
import com.bldover.beacon.ui.components.common.DismissableCard
import com.bldover.beacon.ui.components.editor.SummaryLine
import com.bldover.beacon.ui.screens.editor.genre.GenreSelectorViewModel

@Composable
fun SwipeableGenreCard(
    genre: String,
    onSwipe: (String) -> Unit,
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
fun AddGenreCard(
    onSelect: (String) -> Unit,
    navController: NavController,
    genreSelectorViewModel: GenreSelectorViewModel,
) {
    Box(
        modifier = Modifier.clickable {
            genreSelectorViewModel.launchSelector(navController) {
                onSelect(it)
            }
        }
    ) {
        BasicOutlinedCard {
            SummaryLine(label = "Add Genre") {
                Icon(
                    imageVector = Icons.Default.AddCircle,
                    contentDescription = "Add Genre"
                )
            }
        }
    }
}