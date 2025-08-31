package com.bldover.beacon.ui.components.common

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.focus.focusRequester
import androidx.compose.ui.text.TextRange
import androidx.compose.ui.text.input.TextFieldValue
import com.bldover.beacon.ui.components.editor.SaveCancelButtons

@Composable
fun TextEntryDialog(
    isVisible: Boolean,
    title: String,
    label: String?,
    initialValue: String = "",
    onDismiss: () -> Unit,
    onSave: (String) -> Unit
) {
    if (!isVisible) return

    var value by remember { mutableStateOf(TextFieldValue(
        text = initialValue,
        selection = TextRange(initialValue.length)
    )) }
    val focusRequester = remember { FocusRequester() }

    LaunchedEffect(isVisible) {
        focusRequester.requestFocus()
    }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text(title) },
        text = {
            Column {
                TextField(
                    value = value,
                    onValueChange = { value = it },
                    label = { if (label != null) Text(label) },
                    modifier = Modifier
                        .fillMaxWidth()
                        .focusRequester(focusRequester),
                    singleLine = true
                )
            }
        },
        confirmButton = {
            SaveCancelButtons(
                onCancel = onDismiss,
                onSave = {
                    onSave(value.text)
                }
            )
        },
        dismissButton = null
    )
}
