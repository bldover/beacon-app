package com.bldover.beacon.ui.components.common

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.itemsIndexed
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.RadioButton
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.derivedStateOf
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.runtime.snapshotFlow
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.focus.focusRequester
import androidx.compose.ui.text.TextRange
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.TextFieldValue
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.bldover.beacon.ui.components.editor.SaveCancelButtons
import kotlin.math.abs
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch

@Composable
fun YearPickerDialog(
    isVisible: Boolean,
    selectedYear: Int,
    onDismiss: () -> Unit,
    onYearSelected: (Int) -> Unit
) {
    if (!isVisible) return

    val years = (1000..9999).toList()
    val listState = rememberLazyListState()
    val coroutineScope = rememberCoroutineScope()

    val centeredIndex by remember {
        derivedStateOf {
            val layoutInfo = listState.layoutInfo
            val viewportCenter = layoutInfo.viewportSize.height / 2
            layoutInfo.visibleItemsInfo
                .minByOrNull { abs(it.offset + it.size / 2 - viewportCenter) }
                ?.index ?: years.indexOf(selectedYear)
        }
    }
    val centeredYear = years.getOrNull(centeredIndex) ?: selectedYear

    fun scrollToCenter(index: Int) {
        coroutineScope.launch {
            val layoutInfo = listState.layoutInfo
            val itemHeight = layoutInfo.visibleItemsInfo.firstOrNull()?.size ?: return@launch
            val viewportHeight = layoutInfo.viewportSize.height
            listState.animateScrollToItem(index, itemHeight / 2 - viewportHeight / 2)
        }
    }

    LaunchedEffect(isVisible) {
        val index = years.indexOf(selectedYear)
        if (index >= 0) {
            listState.scrollToItem(maxOf(0, index - 3))
            snapshotFlow { listState.layoutInfo.visibleItemsInfo }
                .first { it.isNotEmpty() }
            val itemHeight = listState.layoutInfo.visibleItemsInfo.first().size
            val viewportHeight = listState.layoutInfo.viewportSize.height
            listState.scrollToItem(index, itemHeight / 2 - viewportHeight / 2)
        }
    }

    LaunchedEffect(listState) {
        snapshotFlow { listState.isScrollInProgress }
            .collect { isScrolling ->
                if (isScrolling) return@collect
                val layoutInfo = listState.layoutInfo
                val itemHeight = layoutInfo.visibleItemsInfo.firstOrNull()?.size ?: return@collect
                val viewportHeight = layoutInfo.viewportSize.height
                val viewportCenter = viewportHeight / 2
                val centeredItem = layoutInfo.visibleItemsInfo
                    .minByOrNull { abs(it.offset + it.size / 2 - viewportCenter) }
                    ?: return@collect
                if (abs(centeredItem.offset - (viewportCenter - itemHeight / 2)) > 1) {
                    listState.animateScrollToItem(centeredItem.index, itemHeight / 2 - viewportHeight / 2)
                }
            }
    }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Select Year") },
        text = {
            LazyColumn(
                state = listState,
                modifier = Modifier.height(300.dp)
            ) {
                itemsIndexed(years) { index, year ->
                    val isCentered = index == centeredIndex
                    Text(
                        text = year.toString(),
                        modifier = Modifier
                            .fillMaxWidth()
                            .clickable {
                                if (isCentered) {
                                    onYearSelected(year)
                                    onDismiss()
                                } else {
                                    scrollToCenter(index)
                                }
                            }
                            .padding(vertical = 12.dp),
                        fontWeight = if (isCentered) FontWeight.Bold else FontWeight.Normal,
                        fontSize = if (isCentered) 22.sp else 16.sp,
                        color = if (isCentered) MaterialTheme.colorScheme.primary else Color.Unspecified,
                        textAlign = TextAlign.Center
                    )
                }
            }
        },
        confirmButton = {
            Button(onClick = {
                onYearSelected(centeredYear)
                onDismiss()
            }) { Text("OK") }
        },
        dismissButton = {
            Button(onClick = onDismiss) { Text("Cancel") }
        }
    )
}

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

@Composable
fun <T> RadioSelectorDialog(
    isVisible: Boolean,
    title: String,
    options: List<T>,
    selectedOption: T,
    getLabel: (T) -> String,
    onDismiss: () -> Unit,
    onOptionSelected: (T) -> Unit
) {
    if (!isVisible) return

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text(title) },
        text = {
            Column {
                options.forEach { option ->
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clickable { onOptionSelected(option); onDismiss() }
                            .padding(vertical = 8.dp),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        RadioButton(
                            selected = option == selectedOption,
                            onClick = { onOptionSelected(option); onDismiss() }
                        )
                        Text(
                            text = getLabel(option),
                            modifier = Modifier.padding(start = 8.dp)
                        )
                    }
                }
            }
        },
        confirmButton = {},
        dismissButton = {
            Button(onClick = onDismiss) { Text("Cancel") }
        }
    )
}
