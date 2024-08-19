package com.bldover.beacon.ui.components

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun <T> ScrollableItemList(
    items: List<T>,
    modifier: Modifier = Modifier,
    verticalArrangement: Arrangement.Vertical = Arrangement.spacedBy(16.dp),
    getItemKey: ((T) -> (String))? = null,
    createItem: @Composable (T) -> Unit = {}
) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = verticalArrangement,
    ) {
        this.items(
            items = items,
            key = getItemKey?.let { getKey -> { item: T -> getKey(item) } }
        ) { item: T ->
            createItem(item)
        }
    }
}