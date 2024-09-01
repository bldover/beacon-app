package com.bldover.beacon.ui.components.common

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import timber.log.Timber

@Composable
fun ScreenFrame(
    topBar: @Composable () -> Unit = {},
    snackbarHost: @Composable () -> Unit = {},
    content: @Composable () -> Unit
) {
    Timber.d("composing ScreenFrame")
    Scaffold(
        topBar = topBar,
        snackbarHost = snackbarHost
    ) { innerPadding ->
        Box(modifier = Modifier.padding(innerPadding)) {
            Box(modifier = Modifier.padding(start = 16.dp, end = 16.dp)) {
                content()
            }
        }
    }
}