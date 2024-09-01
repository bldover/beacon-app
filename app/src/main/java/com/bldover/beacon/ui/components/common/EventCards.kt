package com.bldover.beacon.ui.components.common

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import com.bldover.beacon.data.model.Event
import java.time.format.DateTimeFormatter

@Composable
fun SavedEventCard(
    event: Event,
    onClick: () -> Unit = {},
    showPurchased: Boolean = false
) {
    BasicCard(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
    ) {
        Text(
            text = event.artists.joinToString { it.name },
            style = MaterialTheme.typography.bodyLarge
        )
        Text(
            text = event.date.format(DateTimeFormatter.ISO_DATE),
            style = MaterialTheme.typography.bodySmall
        )
        Text(
            text = event.venue.name,
            style = MaterialTheme.typography.bodySmall
        )
        if (showPurchased) {
            Text(
                text = "Purchased: ${event.purchased}",
                style = MaterialTheme.typography.bodySmall
            )
        }
    }
}