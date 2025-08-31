package com.bldover.beacon.ui.components.editor

import androidx.compose.material3.Switch
import androidx.compose.runtime.Composable
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicOutlinedCard
import timber.log.Timber

@Composable
fun PurchasedSwitch(
    checked: Boolean,
    onChange: (Boolean) -> Unit
) {
    Timber.d("compose PurchasedSwitch : checked=$checked")
    BasicCard {
        SummaryLine(label = "Purchased") {
            Switch(
                checked = checked,
                enabled = true,
                onCheckedChange = onChange
            )
        }
    }
}
