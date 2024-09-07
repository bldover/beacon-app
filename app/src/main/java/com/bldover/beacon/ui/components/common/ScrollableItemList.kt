package com.bldover.beacon.ui.components.common

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import timber.log.Timber

@Composable
fun <T> ScrollableItemList(
    items: List<T>,
    modifier: Modifier = Modifier,
    verticalArrangement: Arrangement.Vertical = Arrangement.spacedBy(16.dp),
    topAnchor: @Composable (() -> Unit)? = null,
    getItemKey: ((T) -> (String))? = null,
    createItem: @Composable (T) -> Unit = {}
) {
    Timber.d("composing ScrollableItemList with ${items.size} items")
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = verticalArrangement,
    ) {
        topAnchor?.let { item { it.invoke() } }
        items(
            items = items,
            key = getItemKey?.let { getKey -> { item: T -> getKey(item) } }
        ) { item: T ->
            Timber.d("composing ScrollableItemList - item $item")
            createItem(item)
        }
    }
}