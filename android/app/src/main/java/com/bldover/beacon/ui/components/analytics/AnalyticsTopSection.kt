package com.bldover.beacon.ui.components.analytics

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.analytics.AnalyticsCount

@Composable
fun AnalyticsTopSection(
    title: String,
    items: List<AnalyticsCount>,
    onSeeAll: () -> Unit,
    onItemClick: ((AnalyticsCount) -> Unit)? = null
) {
    Column(modifier = Modifier.fillMaxWidth()) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = title,
                style = MaterialTheme.typography.titleMedium
            )
            TextButton(onClick = onSeeAll) {
                Text(text = "See all")
            }
        }
        if (items.isEmpty()) {
            Text(
                text = "No data",
                style = MaterialTheme.typography.bodyMedium,
                modifier = Modifier.padding(vertical = 4.dp)
            )
        } else {
            Column(
                verticalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                items.forEach { item ->
                    AnalyticsCountCard(
                        count = item,
                        onClick = onItemClick?.let { { it(item) } }
                    )
                }
            }
        }
    }
}
