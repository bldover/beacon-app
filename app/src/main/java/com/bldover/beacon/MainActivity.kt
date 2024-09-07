package com.bldover.beacon

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.NavigationBottomBar
import com.bldover.beacon.ui.screens.editor.artist.ArtistCreatorScreen
import com.bldover.beacon.ui.screens.editor.artist.ArtistCreatorViewModel
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorScreen
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel
import com.bldover.beacon.ui.screens.editor.event.EventEditorScreen
import com.bldover.beacon.ui.screens.editor.event.EventEditorViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenueCreatorScreen
import com.bldover.beacon.ui.screens.editor.venue.VenueCreatorViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenueSelectorScreen
import com.bldover.beacon.ui.screens.editor.venue.VenueSelectorViewModel
import com.bldover.beacon.ui.screens.saved.HistoryScreen
import com.bldover.beacon.ui.screens.saved.PlannerScreen
import com.bldover.beacon.ui.screens.saved.SavedEventsViewModel
import com.bldover.beacon.ui.screens.upcoming.UpcomingEventsViewModel
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
    savedEventsViewModel: SavedEventsViewModel = hiltViewModel(),
    upcomingEventsViewModel: UpcomingEventsViewModel = hiltViewModel(),
    artistSelectorViewModel: ArtistSelectorViewModel = hiltViewModel(),
    venueSelectorViewModel: VenueSelectorViewModel = hiltViewModel(),
    eventEditorViewModel: EventEditorViewModel = hiltViewModel(),
    artistCreatorViewModel: ArtistCreatorViewModel = hiltViewModel(),
    venueCreatorViewModel: VenueCreatorViewModel = hiltViewModel()
) {
    Timber.d("composing BeaconApp")
    val navController = rememberNavController()
    val settings by userSettingsViewModel.userSettings.collectAsState()
    val coroutineScope = rememberCoroutineScope()
    val snackbarHostState = remember { SnackbarHostState() }

    Scaffold(
        bottomBar = { NavigationBottomBar(navController = navController) },
        snackbarHost = { SnackbarHost(hostState = snackbarHostState) }
    ) { innerPadding ->
        Timber.d("composing BeaconApp - content")
        val snackbarState = remember { SnackbarState(coroutineScope, snackbarHostState) }
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
                        HistoryScreen(
                            navController = navController,
                            snackbarState = snackbarState,
                            savedEventsViewModel = savedEventsViewModel,
                            eventEditorViewModel = eventEditorViewModel
                        )
                    }
                    composable(Screen.CONCERT_PLANNER.name) {
                        PlannerScreen(
                            navController = navController,
                            snackbarState = snackbarState,
                            savedEventsViewModel = savedEventsViewModel,
                            eventEditorViewModel = eventEditorViewModel
                        )
                    }
                    composable(Screen.UPCOMING_EVENTS.name) {
                        UpcomingScreen(
                            navController = navController,
                            snackbarState = snackbarState,
                            eventEditorViewModel = eventEditorViewModel,
                            savedEventsViewModel = savedEventsViewModel,
                            upcomingEventsViewModel = upcomingEventsViewModel
                        )
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
                    composable(Screen.EDIT_EVENT.name) {
                        EventEditorScreen(
                            navController = navController,
                            snackbarState = snackbarState,
                            artistSelectorViewModel = artistSelectorViewModel,
                            artistCreatorViewModel = artistCreatorViewModel,
                            venueSelectorViewModel = venueSelectorViewModel,
                            eventEditorViewModel = eventEditorViewModel
                        )
                    }
                    composable(Screen.SELECT_VENUE.name) {
                        VenueSelectorScreen(
                            navController = navController,
                            venueSelectorViewModel = venueSelectorViewModel,
                            venueCreatorViewModel = venueCreatorViewModel
                        )
                    }
                    composable(Screen.SELECT_ARTIST.name) {
                        ArtistSelectorScreen(
                            navController = navController,
                            artistSelectorViewModel = artistSelectorViewModel,
                            artistCreatorViewModel = artistCreatorViewModel
                        )
                    }
                    composable(Screen.CREATE_ARTIST.name) {
                        ArtistCreatorScreen(
                            navController = navController,
                            artistCreatorViewModel = artistCreatorViewModel
                        )
                    }
                    composable(Screen.CREATE_VENUE.name) {
                        VenueCreatorScreen(
                            navController = navController,
                            venueCreatorViewModel = venueCreatorViewModel
                        )
                    }
                }
            }
        }
    }
}