package com.bldover.beacon

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavType.Companion.StringType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.NavigationBottomBar
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorScreen
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel
import com.bldover.beacon.ui.screens.editor.event.EventEditorScreen
import com.bldover.beacon.ui.screens.editor.event.EventEditorViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenueSelectorScreen
import com.bldover.beacon.ui.screens.editor.venue.VenueSelectorViewModel
import com.bldover.beacon.ui.screens.saved.HistoryScreen
import com.bldover.beacon.ui.screens.saved.PlannerScreen
import com.bldover.beacon.ui.screens.upcoming.UpcomingScreen
import com.bldover.beacon.ui.screens.utility.SettingsState
import com.bldover.beacon.ui.screens.utility.UserSettingsScreen
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
    userSettingsViewModel: UserSettingsViewModel = hiltViewModel(),
    artistSelectorViewModel: ArtistSelectorViewModel = hiltViewModel(),
    venueSelectorViewModel: VenueSelectorViewModel = hiltViewModel(),
    eventEditorViewModel: EventEditorViewModel = hiltViewModel()
) {
    Timber.d("composing BeaconApp")
    val navController = rememberNavController()
    val settings by userSettingsViewModel.userSettings.collectAsState()

    Scaffold(
        bottomBar = { NavigationBottomBar(navController = navController) }
    ) { innerPadding ->
        Timber.d("composing BeaconApp - content")
        when (settings) {
            is SettingsState.Loading -> {
                Timber.d("composing BeaconApp - settings loading")
                LoadingSpinner()
            }
            is SettingsState.Success -> {
                Timber.d("composing BeaconApp - settings loaded: $settings")
                val startDestination = rememberSaveable(settings) {
                    Screen.fromOrDefault((settings as SettingsState.Success).data.startScreen).name
                }

                NavHost(
                    navController = navController,
                    startDestination = startDestination,
                    modifier = Modifier.padding(innerPadding)
                ) {
                    composable(Screen.CONCERT_HISTORY.name) {
                        HistoryScreen(navController = navController)
                    }
                    composable(Screen.CONCERT_PLANNER.name) {
                        PlannerScreen(navController = navController)
                    }
                    composable(Screen.UPCOMING_EVENTS.name) {
                        UpcomingScreen(navController = navController)
                    }
                    composable(Screen.UTILITIES.name) {
                        UtilityScreen(navController = navController)
                    }
                    composable(Screen.USER_SETTINGS.name) {
                        UserSettingsScreen(
                            navController = navController,
                            userSettingsViewModel = userSettingsViewModel
                        )
                    }
                    composable(
                        route = Screen.EDIT_EVENT.name + "/{uuid}/{eventId}",
                        arguments = listOf(
                            navArgument("uuid") { type = StringType },
                            navArgument("eventId") {
                                type = StringType
                            }
                        )
                    ) { backStackEntry ->
                        val eventId = backStackEntry.arguments?.getString("eventId")
                        val uuid = backStackEntry.arguments?.getString("uuid")!!
                        EventEditorScreen(
                            navController = navController,
                            eventId = if (eventId == " ") null else eventId,
                            uuid = uuid,
                            artistSelectorViewModel = artistSelectorViewModel,
                            venueSelectorViewModel = venueSelectorViewModel,
                            eventEditorViewModel = eventEditorViewModel
                        )
                    }
                    composable(Screen.SELECT_VENUE.name) {
                        VenueSelectorScreen(
                            navController = navController,
                            venueSelectorViewModel = venueSelectorViewModel
                        )
                    }
                    composable(Screen.SELECT_ARTIST.name) {
                        ArtistSelectorScreen(
                            navController = navController,
                            artistSelectorViewModel = artistSelectorViewModel
                        )
                    }
                }
            }
        }
    }
}