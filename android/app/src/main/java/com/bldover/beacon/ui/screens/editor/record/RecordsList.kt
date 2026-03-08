package com.bldover.beacon.ui.screens.editor.record

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AddCircle
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.record.Record
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScrollableItemList

@Composable
fun SearchableRecordsList(
    recordState: RecordState,
    onSearchRecords: (String) -> Unit,
    onRecordSelected: (Record) -> Unit,
    onNewRecord: () -> Unit
) {
    Scaffold(
        topBar = {
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = recordState is RecordState.Success,
                onQueryChange = onSearchRecords
            )
        }
    ) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding)) {
            Spacer(modifier = Modifier.height(16.dp))
            when (recordState) {
                is RecordState.Success -> RecordList(
                    recordState.filtered,
                    onRecordSelected,
                    onNewRecord
                )
                is RecordState.Error -> LoadErrorMessage()
                is RecordState.Loading -> LoadingSpinner()
            }
        }
    }
}

@Composable
private fun RecordList(
    records: List<Record>,
    onRecordSelected: (Record) -> Unit,
    onNewRecord: () -> Unit
) {
    ScrollableItemList(
        items = records,
        topAnchor = { NewRecordCard(onNewRecord) },
        getItemKey = { it.id!! }
    ) { record ->
        RecordDetailsCard(
            record = record,
            onClick = { onRecordSelected(record) }
        )
    }
}

@Composable
private fun RecordDetailsCard(
    record: Record,
    onClick: () -> Unit
) {
    BasicCard(modifier = Modifier.clickable { onClick() }) {
        Text(
            text = "${record.artist.name} - ${record.name}",
            style = MaterialTheme.typography.bodyLarge
        )
        Text(
            text = record.year.toString(),
            style = MaterialTheme.typography.bodySmall
        )
    }
}

@Composable
private fun NewRecordCard(onNewRecord: () -> Unit) {
    BasicCard(modifier = Modifier.clickable { onNewRecord() }) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            Text(text = "New Record")
            Icon(
                imageVector = Icons.Default.AddCircle,
                contentDescription = "New Record"
            )
        }
    }
}
