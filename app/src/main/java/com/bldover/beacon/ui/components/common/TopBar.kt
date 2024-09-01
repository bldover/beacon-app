package com.bldover.beacon.ui.components.common

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.IconButtonDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import timber.log.Timber

@Composable
fun TitleTopBar(
    title: String,
    leadingIcon: @Composable () -> Unit = {},
    trailingIcon: @Composable () -> Unit = {}
) {
    Timber.d("composing title top bar with title $title")
    Box(
        modifier = Modifier
            .fillMaxWidth()
            .padding(16.dp)
    ) {
        Box(modifier = Modifier.align(Alignment.CenterStart)) {
            leadingIcon()
        }
        Box(modifier = Modifier.align(Alignment.Center)) {
            Text(
                text = title,
                style = MaterialTheme.typography.titleLarge
            )
        }
        Box(modifier = Modifier.align(Alignment.CenterEnd)) {
            trailingIcon()
        }
    }
}

@Composable
fun BackButton(navController: NavController) {
    IconButton(
        onClick = { navController.popBackStack() },
        colors = IconButtonDefaults.iconButtonColors()
    ) {
        Icon(
            imageVector = Icons.Default.ArrowBack, //todo: update
            contentDescription = "Back button"
        )
    }
}

@Composable
fun RefreshButton(
    onClick: () -> Unit
) {
    IconButton(onClick = onClick) {
        Icon(
            imageVector = Icons.Default.Refresh,
            contentDescription = "Refresh button"
        )
    }
}

@Composable
fun AddButton(
    onClick: () -> Unit
) {
    IconButton(onClick = onClick) {
        Icon(
            imageVector = Icons.Default.Add,
            contentDescription = "Add button"
        )
    }
}