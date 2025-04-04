package com.bldover.beacon.ui.components.common

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Card
import androidx.compose.material3.CardColors
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.event.Event
import com.bldover.beacon.data.model.event.EventDetail
import java.time.format.DateTimeFormatter

@Composable
fun accentCardColors(): CardColors {
    return CardColors(
        containerColor = MaterialTheme.colorScheme.primaryContainer,
        contentColor = MaterialTheme.colorScheme.onPrimaryContainer,
        disabledContainerColor = MaterialTheme.colorScheme.primaryContainer,
        disabledContentColor = MaterialTheme.colorScheme.onPrimaryContainer
    )
}

@Composable
fun SavedEventCard(
    event: Event,
    accented: Boolean,
    onClick: () -> Unit = {}
) {
    BasicCard(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick),
        border = if (accented) BorderStroke(width = 1.dp, color = MaterialTheme.colorScheme.primary) else null
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
    }
}

@Composable
fun UpcomingEventCard(
    event: EventDetail,
    accented: Boolean,
    onClick: () -> Unit = {}
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick),
        border = if (accented) BorderStroke(width = 1.dp, color = MaterialTheme.colorScheme.primary) else null
    ) {
        Column(
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
        ) {
            if (event.artists.isNotEmpty()) {
                val artists = event.artists.joinToString { it.name }
                Text(
                    text = artists,
                    style = MaterialTheme.typography.bodyLarge
                )
                if (event.name.isNotBlank() && event.name != event.artists.first().name) {
                    Text(
                        text = event.name,
                        style = MaterialTheme.typography.bodySmall
                    )
                }
            } else {
                Text(
                    text = event.name,
                    style = MaterialTheme.typography.bodyLarge
                )
            }
            Text(
                text = event.date.format(DateTimeFormatter.ISO_DATE),
                style = MaterialTheme.typography.bodySmall
            )
            Text(
                text = event.venue.name,
                style = MaterialTheme.typography.bodySmall
            )
            if (event.price != null) {
                Text(
                    text = "Price: $${event.formattedPrice}",
                    style = MaterialTheme.typography.bodySmall
                )
            }
            if (event.rank != null) {
                Text(
                    text = "Rank: ${event.rank}",
                    style = MaterialTheme.typography.bodySmall
                )
            }
            event.artistRanks
                ?.asSequence()
                ?.mapNotNull { it.relatedArtists }
                ?.flatten()
                ?.distinct()
                ?.filterNot { relatedArtist ->
                    event.artists
                        .map { it.name }
                        .any { relatedArtist.equals(it, ignoreCase = true) }
                }
                ?.joinToString()
                ?.takeIf { it.isNotBlank() }
                ?.let {
                    Text(
                        text = "Related Artists: $it",
                        style = MaterialTheme.typography.bodySmall
                    )
                }
        }
    }
}