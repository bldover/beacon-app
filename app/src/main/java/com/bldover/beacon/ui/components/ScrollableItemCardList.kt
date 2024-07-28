package com.bldover.beacon.ui.components

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun <T> ScrollableItemCardList(
    items: List<T>,
    modifier: Modifier = Modifier,
    verticalArrangement: Arrangement.Vertical = Arrangement.spacedBy(16.dp),
    createItem: @Composable (T) -> Unit = { }
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = verticalArrangement
    ) {
        this.items(items = items) { item: T ->
            createItem(item)
        }
    }
}