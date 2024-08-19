package com.bldover.beacon.ui.screens.history

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.bldover.beacon.data.model.Event
import com.bldover.beacon.ui.components.BasicSearchBar
import com.bldover.beacon.ui.components.EventCard
import com.bldover.beacon.ui.screens.editor.EventEditBottomSheet
import com.bldover.beacon.ui.components.LoadErrorMessage
import com.bldover.beacon.ui.components.LoadingSpinner
import com.bldover.beacon.ui.components.ScrollableItemList
import com.bldover.beacon.ui.screens.editor.showSheet

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HistoryScreen(
    historyViewModel: HistoryViewModel = hiltViewModel()
) {
    val historyState = historyViewModel.uiState.collectAsState()
    var selectedEvent by remember { mutableStateOf<Event?>(null) }
    var showEditorSheet by remember { mutableStateOf(false) }
    val editorSheetState = rememberModalBottomSheetState(skipPartiallyExpanded = true)
    val scope = rememberCoroutineScope()
    Scaffold(
        topBar = {
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = historyState.value is UiState.Success
            ) {
                historyViewModel.handleEvent(UiEvent.ApplySearchFilter(it))
            }
        }
    ) { innerPadding ->
        Column {
            Spacer(modifier = Modifier.height(16.dp))
            when (historyState.value) {
                is UiState.Success -> {
                    ScrollableItemList(
                        items = (historyState.value as UiState.Success).filtered,
                        modifier = Modifier.padding(innerPadding),
                        getItemKey = { it.id }
                    ) {
                        Box(
                            modifier = Modifier.clickable {
                                selectedEvent = it
                                showEditorSheet = true
                                showSheet(scope, editorSheetState)
                            }
                        ) {
                            EventCard(it)
                        }
                    }
                    if (showEditorSheet) {
                        EventEditBottomSheet(
                            event = selectedEvent!!,
                            modalSheetState = editorSheetState,
                            onClosed = { showEditorSheet = false }
                        )
                    }
                }
                is UiState.Error -> {
                    LoadErrorMessage()
                }
                is UiState.Loading -> {
                    LoadingSpinner()
                }
            }
        }
    }
}