package com.bldover.beacon.ui.components.common

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.KeyboardArrowDown
import androidx.compose.material.icons.filled.KeyboardArrowUp
import androidx.compose.material3.Card
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.EventOrdering
import com.bldover.beacon.data.model.Order
import com.bldover.beacon.data.model.OrderType
import com.bldover.beacon.data.model.RecommendationThreshold

@Composable
fun EventSearchUtilityBar(
    state: EventOrdering,
    onChange: (EventOrdering) -> Unit
) {
    LazyRow(
        horizontalArrangement = Arrangement.spacedBy(8.dp),
    ) {
        item { OrderToggle(OrderType.DATE, state, onChange) }
        item { OrderToggle(OrderType.VENUE, state, onChange) }
    }
}

@Composable
fun OrderToggle(
    orderType: OrderType,
    state: EventOrdering,
    onChange: (EventOrdering) -> Unit
) {
    val optionSelected = state.option == orderType
    FilterCard(
        onClick = {
            val order = if (optionSelected) {
                if (state.order == Order.ASCENDING) Order.DESCENDING
                else Order.ASCENDING
            }
            else Order.ASCENDING
            onChange(EventOrdering(option = orderType, order = order))
        },
        border = if (optionSelected) BorderStroke(1.dp, MaterialTheme.colorScheme.primary) else null
    ) {
        Row(
            modifier = Modifier
                .height(24.dp)
                .padding(if (!optionSelected) PaddingValues(horizontal = 4.dp) else PaddingValues()),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = orderType.label,
                style = MaterialTheme.typography.bodyMedium,
                textAlign = if (optionSelected) TextAlign.Start else TextAlign.Center
            )
            if (optionSelected) {
                Icon(
                    imageVector = if (state.order == Order.ASCENDING)
                        Icons.Default.KeyboardArrowDown else Icons.Default.KeyboardArrowUp,
                    contentDescription = if (state.order == Order.ASCENDING)
                        "Ascending" else "Descending"
                )
            }
        }
    }
}

@Composable
fun RecommendationSelectionBar(
    state: RecommendationThreshold,
    onChange: (RecommendationThreshold) -> Unit
) {
    LazyRow(
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        item {
            FilterCard(
                onClick = { onChange(RecommendationThreshold.NONE) },
                border = if (state == RecommendationThreshold.NONE) BorderStroke(1.dp, MaterialTheme.colorScheme.primary) else null
            ) {
                Box(
                    modifier = Modifier
                        .height(24.dp)
                        .padding(horizontal = 4.dp),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        text = "All",
                        style = MaterialTheme.typography.bodyMedium,
                        textAlign = TextAlign.Center
                    )
                }
            }
        }
        item {
            var expanded by remember { mutableStateOf(false) }
            FilterCard(
                onClick = { expanded = !expanded },
                border = if (state != RecommendationThreshold.NONE) BorderStroke(1.dp, MaterialTheme.colorScheme.primary) else null
            ) {
                Box(
                    modifier = Modifier
                        .height(24.dp)
                        .padding(horizontal = 4.dp),
                    contentAlignment = Alignment.Center
                ) {
                    RecommendationSelector(
                        expanded = expanded,
                        state = state,
                        onChange = onChange,
                        onDismiss = { expanded = false }
                    )
                }
            }
        }
    }
}

@Composable
fun RecommendationSelector(
    expanded: Boolean,
    state: RecommendationThreshold,
    onChange: (RecommendationThreshold) -> Unit,
    onDismiss: () -> Unit
) {
    Text(
        text = if (state == RecommendationThreshold.NONE) "Recommended" else "Recommended: ${state.label}",
        style = MaterialTheme.typography.bodyMedium
    )
    DropdownMenu(
        expanded = expanded,
        onDismissRequest = onDismiss
    ) {
        val options = listOf(
            RecommendationThreshold.LOW,
            RecommendationThreshold.MEDIUM,
            RecommendationThreshold.HIGH
        )
        options.forEach {
            DropdownMenuItem(
                text = { Text(it.label) },
                onClick = {
                    onDismiss()
                    onChange(it)
                }
            )
        }
    }
}

@Composable
fun FilterCard(
    onClick: () -> Unit,
    border: BorderStroke? = null,
    content: @Composable () -> Unit
) {
    Card(
        onClick = onClick,
        border = border
    ) {
        Box(modifier = Modifier.padding(vertical = 4.dp, horizontal = 8.dp)) {
            content()
        }
    }
}

@Preview
@Composable
fun EventSearchUtilityBarPreview() {
    EventSearchUtilityBar(
        state = EventOrdering(OrderType.DATE, Order.DESCENDING),
        onChange = {}
    )
}

@Preview
@Composable
fun RecommendationSelectionBarPreview() {
    RecommendationSelectionBar(
        state = RecommendationThreshold.LOW,
        onChange = {}
    )
}