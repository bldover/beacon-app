package com.bldover.beacon.ui.components.editor

import androidx.compose.material3.LocalMinimumInteractiveComponentSize
import androidx.compose.material3.Switch
import androidx.compose.runtime.Composable
import androidx.compose.runtime.CompositionLocalProvider
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import timber.log.Timber

@Preview
@Composable
fun ReducedMinSizeSwitchPreview() {
    ReducedMinSizeSwitch(
        label = "Purchased",
        checked = true,
        onChange = {}
    )
}

@Composable
fun ReducedMinSizeSwitch(
    label: String = "",
    checked: Boolean,
    minimumSize: Dp = 32.dp,
    onChange: (Boolean) -> Unit
) {
    Timber.d("compose ReducedMinSizeSwitch : checked=$checked")
    SummaryCard(label = label) {
        CompositionLocalProvider(LocalMinimumInteractiveComponentSize provides minimumSize) {
            Switch(
                checked = checked,
                enabled = true,
                onCheckedChange = onChange
            )
        }
    }
}
