package com.bldover.beacon.data.model

import androidx.compose.material3.SnackbarHostState
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch

class SnackbarState(
    val appScope: CoroutineScope,
    val snackbarHostState: SnackbarHostState
) {
    fun showSnackbar(message: String) {
        appScope.launch {
            snackbarHostState.showSnackbar(message)
        }
    }
}