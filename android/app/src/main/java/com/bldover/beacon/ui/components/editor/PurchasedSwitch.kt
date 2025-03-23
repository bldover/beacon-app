package com.bldover.beacon.ui.components.editor

import androidx.compose.material3.Switch
import androidx.compose.runtime.Composable
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicOutlinedCard
import timber.log.Timber

@Composable
fun PurchasedSwitch(
    checked: Boolean,
    enabled: Boolean = true,
    onChange: (Boolean) -> Unit
) {
    Timber.d("compose PurchasedSwitch : checked=$checked, enabled=$enabled")
    if (enabled) {
        BasicCard {
            SummaryLine(label = "Purchased") {
                Switch(
                    checked = checked,
                    enabled = true,
                    onCheckedChange = onChange
                )
            }
        }
    } else {
        BasicOutlinedCard {
            SummaryLine(label = "Purchased") {
                Switch(
                    checked = false,
                    enabled = false,
                    onCheckedChange = onChange
                )
            }
        }
    }
}