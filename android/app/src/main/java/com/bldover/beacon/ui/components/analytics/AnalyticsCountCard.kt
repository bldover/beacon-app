package com.bldover.beacon.ui.components.analytics

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import com.bldover.beacon.data.model.analytics.AnalyticsCount
import com.bldover.beacon.ui.components.common.BasicCard

@Composable
fun AnalyticsCountCard(
    count: AnalyticsCount,
    onClick: (() -> Unit)? = null
) {
    val baseModifier = Modifier.fillMaxWidth()
    val modifier = if (onClick != null) baseModifier.clickable(onClick = onClick) else baseModifier
    BasicCard(modifier = modifier) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = count.name,
                style = MaterialTheme.typography.bodyLarge
            )
            Text(
                text = count.count.toString(),
                style = MaterialTheme.typography.bodyLarge
            )
        }
    }
}
