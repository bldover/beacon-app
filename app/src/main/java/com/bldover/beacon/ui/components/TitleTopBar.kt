package com.bldover.beacon.ui.components

import androidx.compose.material3.CenterAlignedTopAppBar
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.IconButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.tooling.preview.Preview
import com.bldover.beacon.ActiveScreen
import com.bldover.beacon.ui.theme.BeaconTheme

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TitleTopBar(
    activeScreen: ActiveScreen
) {
    CenterAlignedTopAppBar(
        title = { Text(text = activeScreen.title) },
        navigationIcon = {
            when (activeScreen) {
                ActiveScreen.HISTORY_DETAIL -> IconButton(onClick = { /*TODO*/ }, ) {

                }
                else -> Unit
            }
        }
    )
}

@Preview
@Composable
fun TitleTopBarPreview(
    activeScreen: ActiveScreen = ActiveScreen.CONCERT_HISTORY
) {
    BeaconTheme(darkTheme = true, dynamicColor = true) {
        TitleTopBar(activeScreen)
    }
}