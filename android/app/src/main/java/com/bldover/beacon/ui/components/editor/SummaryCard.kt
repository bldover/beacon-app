package com.bldover.beacon.ui.components.editor

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import com.bldover.beacon.ui.components.common.BasicCard

@Composable
fun SummaryCard(
    label: String,
    modifier: Modifier = Modifier,
    onClick: () -> Unit = {},
    showBorder: Boolean = false,
    borderStroke: BorderStroke? = null,
    content: @Composable () -> Unit
) {
    BasicCard(
        modifier = if (onClick != {}) modifier.clickable { onClick() } else modifier,
        border = if (showBorder) borderStroke else null
    ) {
        SummaryLine(label = label) {
            content()
        }
    }
}

@Composable
fun SummaryLine(
    label: String,
    content: @Composable () -> Unit
) {
    val labelWidthWeight = 0.35f
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically
    ) {
        Box(
            modifier = Modifier.weight(labelWidthWeight),
            contentAlignment = Alignment.CenterStart
        ) {
            Text(text = label)
        }
        Box(
            modifier = Modifier.weight(1 - labelWidthWeight),
            contentAlignment = Alignment.CenterEnd
        ) {
            Column(
                modifier = Modifier.fillMaxWidth(),
                horizontalAlignment = Alignment.End
            ) {
                content()
            }
        }
    }
}