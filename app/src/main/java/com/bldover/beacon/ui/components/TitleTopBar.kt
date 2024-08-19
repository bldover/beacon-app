package com.bldover.beacon.ui.components

import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.CenterAlignedTopAppBar
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.IconButtonDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import androidx.navigation.compose.currentBackStackEntryAsState
import com.bldover.beacon.data.model.Screen

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TitleTopBar(
    title: String,
    showBackButton: Boolean = false,
    navController: NavController
) {
    CenterAlignedTopAppBar(
        navigationIcon = { if (showBackButton) BackButton(navController) },
        title = { Text(text = title) }
    )
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