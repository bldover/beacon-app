package com.bldover.beacon.ui.components

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Card
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.Event
import com.bldover.beacon.data.model.EventDetail
import java.time.format.DateTimeFormatter

@Composable
fun EventCard(
    event: Event,
    showPurchased: Boolean = false
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
        //.clickable { () -> { /* TODO: Something with event */} }
    ) {
        Column(
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
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
}

@Composable
fun EventDetailCard(
    event: EventDetail
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
        //.clickable { () -> { /* TODO: Something with event */} }
    ) {
        Column(
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
        ) {
            val artistsName = if (event.artists.isNotEmpty()) {
                event.artists.joinToString { it.name }
            }
            else {
                event.name
            }
            Text(
                text = artistsName,
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
            if (event.purchased) {
                Text(
                    text = "Going!",
                    style = MaterialTheme.typography.bodySmall
                )
            } else if (event.price != null) {
                Text(
                    text = "Price: $${event.formattedPrice}",
                    style = MaterialTheme.typography.bodySmall
                )
            }
        }
    }
}

