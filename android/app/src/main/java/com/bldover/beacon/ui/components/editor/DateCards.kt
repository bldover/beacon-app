package com.bldover.beacon.ui.components.editor

import androidx.compose.material3.DatePicker
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Text
import androidx.compose.material3.rememberDatePickerState
import androidx.compose.runtime.Composable
import com.bldover.beacon.ui.components.common.ExpandableCard
import timber.log.Timber
import java.time.LocalDate

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DateEditCard(
    date: LocalDate,
    onChange: (LocalDate) -> Unit,
) {
    Timber.d("composing DateEditCard : $date")
    ExpandableCard(
        headlineContent = {
            SummaryLine(label = "Date") {
                Text(text = date.toString())
            }
        }
    ) {
        val datePickerState = rememberDatePickerState(
            initialSelectedDateMillis = date.toEpochDay() * 1000 * 60 * 60 * 24
        )
        DatePicker(state = datePickerState)
        val selectedDate = LocalDate.ofEpochDay(
            datePickerState.selectedDateMillis!! / 1000 / 60 / 60 / 24
        )
        if (selectedDate != date) {
            onChange(selectedDate)
        }
    }
}