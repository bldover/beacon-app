package com.bldover.beacon.ui.components.analytics

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.runtime.Composable
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.OrderField
import com.bldover.beacon.data.model.analytics.AnalyticsOrdering
import com.bldover.beacon.ui.components.common.OrderToggle

@Composable
fun AnalyticsSortBar(
    state: AnalyticsOrdering,
    onChange: (AnalyticsOrdering) -> Unit
) {
    LazyRow(
        horizontalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        item { OrderToggle(OrderField.COUNT, state) { field, order -> onChange(AnalyticsOrdering(field, order)) } }
        item { OrderToggle(OrderField.NAME, state) { field, order -> onChange(AnalyticsOrdering(field, order)) } }
    }
}
