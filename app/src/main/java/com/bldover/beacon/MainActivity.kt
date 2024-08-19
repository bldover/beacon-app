package com.bldover.beacon

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.ui.components.NavigationBottomBar
import com.bldover.beacon.ui.components.TitleTopBar
import com.bldover.beacon.ui.screens.history.HistoryScreen
import com.bldover.beacon.ui.screens.planner.PlannerScreen
import com.bldover.beacon.ui.screens.upcoming.UpcomingScreen
import com.bldover.beacon.ui.screens.utility.SettingsState
import com.bldover.beacon.ui.screens.utility.UserSettingsList
import com.bldover.beacon.ui.screens.utility.UserSettingsViewModel
import com.bldover.beacon.ui.screens.utility.UtilityScreen
import com.bldover.beacon.ui.theme.BeaconTheme
import dagger.hilt.android.AndroidEntryPoint
import timber.log.Timber

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            BeaconTheme {
                BeaconApp()
            }
        }
    }
}

@Composable
fun BeaconApp(
    userSettingsViewModel: UserSettingsViewModel = hiltViewModel()
) {
    val navController = rememberNavController()
    when (val settings = userSettingsViewModel.userSettings.collectAsState().value) {
        is SettingsState.Loading -> {}
        is SettingsState.Success -> {
            Timber.d("settings loaded: $settings")
            AppFrame(navController = navController) {
                NavHost(
                    navController = navController,
                    startDestination = remember {
                        Screen.fromOrDefault(settings.data.startScreen).name
                    }
                ) {
                    composable(Screen.CONCERT_HISTORY.name) { HistoryScreen() }
                    composable(Screen.CONCERT_PLANNER.name) { PlannerScreen() }
                    composable(Screen.UPCOMING_EVENTS.name) { UpcomingScreen() }
                    composable(Screen.UTILITIES.name) { UtilityScreen(navController) }
                    composable(Screen.USER_SETTINGS.name) { UserSettingsList() }
                }
            }
        }
    }
}

@Composable
fun AppFrame(
    navController: NavController,
    content: @Composable () -> Unit
) {
    val activeScreen = Screen.fromOrDefault(navController.currentBackStackEntryAsState().value?.destination?.route)
    Scaffold(
        topBar = { TitleTopBar(
            title = activeScreen.title,
            showBackButton = activeScreen.subScreen,
            navController = navController
        ) },
        bottomBar = { NavigationBottomBar(navController = navController) }
    ) { innerPadding ->
        Box(modifier = Modifier.padding(innerPadding)) {
            Box(modifier = Modifier.padding(16.dp)) {
                content()
            }
        }
    }
}
