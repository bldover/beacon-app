package com.bldover.beacon.ui.screens.utility

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.screens.editor.record.RecordEditorViewModel
import com.bldover.beacon.ui.screens.editor.record.RecordsViewModel
import com.bldover.beacon.ui.screens.editor.record.SearchableRecordsList
import timber.log.Timber

@Composable
fun ManageRecordsScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    recordEditorViewModel: RecordEditorViewModel,
    recordsViewModel: RecordsViewModel = hiltViewModel()
) {
    Timber.d("composing ManageRecordsScreen")
    LaunchedEffect(Unit) {
        recordsViewModel.resetFilter()
    }
    val recordState by recordsViewModel.uiState.collectAsState()

    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Manage Records",
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        SearchableRecordsList(
            recordState = recordState,
            onSearchRecords = recordsViewModel::applyFilter,
            onRecordSelected = { record ->
                recordEditorViewModel.launchEditor(
                    navController = navController,
                    record = record,
                    onSave = { updated ->
                        recordsViewModel.updateRecord(
                            record = updated,
                            onSuccess = { navController.popBackStack() },
                            onError = { msg -> snackbarState.showSnackbar(msg) }
                        )
                    },
                    onDelete = { toDelete ->
                        recordsViewModel.deleteRecord(
                            record = toDelete,
                            onSuccess = { navController.popBackStack() },
                            onError = { msg -> snackbarState.showSnackbar(msg) }
                        )
                    }
                )
            },
            onNewRecord = {
                recordEditorViewModel.launchEditor(
                    navController = navController,
                    onSave = { newRecord ->
                        recordsViewModel.addRecord(
                            record = newRecord,
                            onSuccess = { navController.popBackStack() },
                            onError = { msg -> snackbarState.showSnackbar(msg) }
                        )
                    }
                )
            }
        )
    }
}
