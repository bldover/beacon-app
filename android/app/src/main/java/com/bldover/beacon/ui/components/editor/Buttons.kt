package com.bldover.beacon.ui.components.editor

import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonColors
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun SaveCancelButtons(
    onCancel: () -> Unit,
    onSave: () -> Unit
) {
    Row(modifier = Modifier.padding(vertical = 8.dp)) {
        Button(onClick = onCancel) {
            Text(text = "Cancel")
        }
        Spacer(modifier = Modifier.width(8.dp))
        Button(onClick = onSave) {
            Text(text = "Save")
        }
    }
}

@Composable
fun DeleteButton(
    onDelete: () -> Unit
) {
    Button(
        colors = ButtonColors(
            containerColor = MaterialTheme.colorScheme.errorContainer,
            contentColor = MaterialTheme.colorScheme.onErrorContainer,
            disabledContainerColor = MaterialTheme.colorScheme.errorContainer,
            disabledContentColor = MaterialTheme.colorScheme.onErrorContainer
        ),
        onClick = onDelete
    ) {
        Text(text = "Delete")
    }
}