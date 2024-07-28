package com.bldover.beacon

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import com.bldover.beacon.ui.components.NavigationBottomBar
import com.bldover.beacon.ui.components.TitleTopBar
import com.bldover.beacon.ui.screens.history.HistoryScreen
import com.bldover.beacon.ui.screens.UpcomingScreen
import com.bldover.beacon.ui.screens.UtilityScreen
import com.bldover.beacon.ui.theme.BeaconTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            BeaconTheme {
                BeaconApplication()
            }
        }
    }
}

enum class ActiveScreen(val title: String, val shortDesc: String) {
    CONCERT_HISTORY("Concert History", "History"),
    HISTORY_DETAIL("", ""),
    UPCOMING_EVENTS("Upcoming Events", "Upcoming"),
    UTILITIES("Utilities", "Utilities")
}

@Composable
fun BeaconApplication() {
    var activeScreen by remember { mutableStateOf(ActiveScreen.CONCERT_HISTORY) }
    Scaffold(
        topBar = { TitleTopBar(activeScreen = activeScreen) },
        bottomBar = { NavigationBottomBar(
            activeScreen = activeScreen,
            changeScreen = fun (a) {activeScreen = a}
        ) }
    ) { innerPadding ->
        Box(modifier = Modifier.padding(innerPadding)) {
            when (activeScreen) {
                ActiveScreen.CONCERT_HISTORY -> HistoryScreen()
                ActiveScreen.HISTORY_DETAIL -> {}
                ActiveScreen.UPCOMING_EVENTS -> UpcomingScreen()
                ActiveScreen.UTILITIES -> UtilityScreen()
            }
        }
    }
}

@Preview
@Composable
fun BeaconApplicationPreview() {
    BeaconTheme(darkTheme = true, dynamicColor = true) {
        BeaconApplication()
    }
}