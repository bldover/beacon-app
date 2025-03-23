package com.bldover.beacon.ui.components.editor

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.Venue
import com.bldover.beacon.ui.components.common.BasicCard
import timber.log.Timber

@Composable
fun VenueCard(
    venue: Venue,
    onClick: () -> Unit = {}
) {
    BasicCard(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
    ) {
        Text(
            text = venue.name,
            style = MaterialTheme.typography.bodyLarge
        )
        Text(
            text = "${venue.city}, ${venue.state}",
            style = MaterialTheme.typography.bodySmall
        )
    }
}

@Composable
fun VenueEditCard(
    venue: Venue,
    onClick: () -> Unit = {}
) {
    Timber.d("composing VenueEditCard : $venue")
    BasicCard(
        modifier = Modifier.clickable { onClick() },
        border = if (venue.id != null) BorderStroke(width = 1.dp, color = MaterialTheme.colorScheme.primary) else null
    ) {
        SummaryLine(label = "Venue") {
            Text(
                text = venue.name
            )
            Text(
                text = "${venue.city}, ${venue.state}",
                style = MaterialTheme.typography.bodySmall
            )
        }
    }
}