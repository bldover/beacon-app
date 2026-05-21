package com.bldover.beacon.ui.components.editor

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyListScope
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun SaveableEditFieldsColumn(
    onSave: () -> Unit,
    onCancel: () -> Unit,
    modifier: Modifier = Modifier,
    showDelete: Boolean = false,
    onDelete: () -> Unit = {},
    verticalArrangement: Arrangement.Vertical = Arrangement.spacedBy(16.dp),
    content: LazyListScope.() -> Unit
) {

    LazyColumn(
        modifier = modifier.fillMaxWidth(),
        verticalArrangement = verticalArrangement
    ) {
        content()
        item {
            Row(
                horizontalArrangement = Arrangement.End,
                modifier = Modifier.fillMaxWidth()
            ) {
                if (showDelete) {
                    DeleteButton(onDelete = onDelete)
                    Spacer(modifier = Modifier.weight(1f))
                }
                SaveCancelButtons(
                    onSave = onSave,
                    onCancel = onCancel
                )
            }
        }
    }
}