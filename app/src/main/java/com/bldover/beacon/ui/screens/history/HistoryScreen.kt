package com.bldover.beacon.ui.screens.history

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Card
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.bldover.beacon.ui.components.BasicSearchBar
import com.bldover.beacon.ui.components.ScrollableItemCardList
import com.bldover.beacon.ui.theme.BeaconTheme
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HistoryScreen() {
    var events by remember { mutableStateOf(listOf(
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent(),
        getDummyEvent()
    )) }
    var filteredEvents by remember { mutableStateOf(events) }
    Box(
        modifier = Modifier.padding(horizontal = 16.dp)
    ) {
        Scaffold(
            topBar = {
                BasicSearchBar(
                    modifier = Modifier.fillMaxWidth()
                ) {
                    filteredEvents = events
                        .filter { e -> e.artists.any { a -> a.name.contains(it, true) } }
                        .toList()
                }
            }
        ) { innerPadding ->
            Column {
                Spacer(modifier = Modifier.height(16.dp))
                ScrollableItemCardList(
                    items = filteredEvents,
                    modifier = Modifier.padding(innerPadding)
                ) { event ->
                    EventDetails(event)
                }
            }
        }
    }
}

@Composable
fun EventDetails(event: Event) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            //.clickable { () -> { /* TODO: Something with event */} }
    ) {
        Column(
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
        ) {
            Text(
                text = event.artists.joinToString { a -> a.name },
                style = MaterialTheme.typography.bodyLarge
            )
            Text(
                text = "${event.date?.format(DateTimeFormatter.ISO_DATE)}",
                style = MaterialTheme.typography.bodySmall
            )
            Text(
                text = "${event.venue?.name}",
                style = MaterialTheme.typography.bodySmall
            )
        }
    }
}

@Preview
@Composable
fun HistoryScreenPreview() {
    BeaconTheme(darkTheme = true, dynamicColor = true) {
        HistoryScreen()
    }
}