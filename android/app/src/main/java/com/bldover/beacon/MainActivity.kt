package com.bldover.beacon

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.imePadding
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
import com.bldover.beacon.data.spotify.SpotifyAuthResult
import com.bldover.beacon.data.spotify.SpotifyAuthResultBus

import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.NavigationBottomBar
import com.bldover.beacon.ui.screens.albums.AlbumDetailsScreen
import com.bldover.beacon.ui.screens.albums.AlbumDetailsViewModel
import com.bldover.beacon.ui.screens.analytics.AnalyticsArtistsScreen
import com.bldover.beacon.ui.screens.analytics.AnalyticsEventsScreen
import com.bldover.beacon.ui.screens.analytics.AnalyticsEventsViewModel
import com.bldover.beacon.ui.screens.analytics.AnalyticsGenresScreen
import com.bldover.beacon.ui.screens.analytics.AnalyticsListViewModel
import com.bldover.beacon.ui.screens.analytics.AnalyticsMonthsScreen
import com.bldover.beacon.ui.screens.analytics.AnalyticsOverviewScreen
import com.bldover.beacon.ui.screens.analytics.AnalyticsOverviewViewModel
import com.bldover.beacon.ui.screens.analytics.AnalyticsVenuesScreen
import com.bldover.beacon.ui.screens.analytics.AnalyticsYearsScreen
import com.bldover.beacon.ui.screens.editor.album.AlbumEditorScreen
import com.bldover.beacon.ui.screens.editor.album.AlbumEditorViewModel
import com.bldover.beacon.ui.screens.editor.album.AlbumsViewModel
import com.bldover.beacon.ui.screens.editor.artist.ArtistEditorScreen
import com.bldover.beacon.ui.screens.editor.artist.ArtistEditorViewModel
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorScreen
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel
import com.bldover.beacon.ui.screens.editor.event.EventEditorScreen
import com.bldover.beacon.ui.screens.editor.event.EventEditorViewModel
import com.bldover.beacon.ui.screens.editor.genre.GenreSelectorScreen
import com.bldover.beacon.ui.screens.editor.genre.GenreSelectorViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenueEditorScreen
import com.bldover.beacon.ui.screens.editor.venue.VenueEditorViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenueSelectorScreen
import com.bldover.beacon.ui.screens.editor.venue.VenueSelectorViewModel
import com.bldover.beacon.ui.screens.events.HistoryScreen
import com.bldover.beacon.ui.screens.events.PlannerScreen
import com.bldover.beacon.ui.screens.events.SavedEventsViewModel
import com.bldover.beacon.ui.screens.events.UpcomingEventsViewModel
import com.bldover.beacon.ui.screens.events.UpcomingScreen
import com.bldover.beacon.ui.screens.utility.ManageAlbumsScreen
import com.bldover.beacon.ui.screens.utility.ManageArtistsScreen
import com.bldover.beacon.ui.screens.utility.ManageGenresScreen
import com.bldover.beacon.ui.screens.utility.ManageVenuesScreen
import com.bldover.beacon.ui.screens.utility.SettingsState
import com.bldover.beacon.ui.screens.utility.SpotifyAuthNotification
import com.bldover.beacon.ui.screens.utility.SpotifyAuthScreen
import com.bldover.beacon.ui.screens.utility.SpotifyAuthViewModel
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
        handleSpotifyAuthIntent(intent)
        setContent {
            BeaconTheme {
                BeaconApp()
            }
        }
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        setIntent(intent)
        handleSpotifyAuthIntent(intent)
    }

    private fun handleSpotifyAuthIntent(intent: Intent?) {
        val data = intent?.data ?: return
        if (data.scheme != "beacon" || data.host != "spotify-auth-complete") return
        val status = data.getQueryParameter("status")
        val reason = data.getQueryParameter("reason")
        Timber.d("Spotify auth deep link: status=$status reason=$reason")
        val result = if (status == "ok") {
            SpotifyAuthResult.Success
        } else {
            SpotifyAuthResult.Failure(reason)
        }
        SpotifyAuthResultBus.post(result)
    }
}

@Composable
fun BeaconApp(
    userSettingsViewModel: UserSettingsViewModel = hiltViewModel(),
    savedEventsViewModel: SavedEventsViewModel = hiltViewModel(),
    upcomingEventsViewModel: UpcomingEventsViewModel = hiltViewModel(),
    artistSelectorViewModel: ArtistSelectorViewModel = hiltViewModel(),
    venueSelectorViewModel: VenueSelectorViewModel = hiltViewModel(),
    genreSelectorViewModel: GenreSelectorViewModel = hiltViewModel(),
    eventEditorViewModel: EventEditorViewModel = hiltViewModel(),
    artistEditorViewModel: ArtistEditorViewModel = hiltViewModel(),
    venueEditorViewModel: VenueEditorViewModel = hiltViewModel(),
    albumEditorViewModel: AlbumEditorViewModel = hiltViewModel(),
    albumsViewModel: AlbumsViewModel = hiltViewModel(),
    albumDetailsViewModel: AlbumDetailsViewModel = hiltViewModel(),
    analyticsOverviewViewModel: AnalyticsOverviewViewModel = hiltViewModel(),
    analyticsListViewModel: AnalyticsListViewModel = hiltViewModel(),
    analyticsEventsViewModel: AnalyticsEventsViewModel = hiltViewModel(),
    spotifyAuthViewModel: SpotifyAuthViewModel = hiltViewModel()
) {
    Timber.d("composing BeaconApp")
    val navController = rememberNavController()
    val settings by userSettingsViewModel.userSettings.collectAsState()
    val coroutineScope = rememberCoroutineScope()
    val snackbarHostState = remember { SnackbarHostState() }

    SpotifyAuthNotification(spotifyAuthViewModel = spotifyAuthViewModel)

    Scaffold(
        bottomBar = { NavigationBottomBar(navController = navController) },
        snackbarHost = { SnackbarHost(hostState = snackbarHostState, modifier = Modifier.imePadding()) }
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
                        UtilityScreen(
                            navController = navController
                        )
                    }
                    composable(Screen.SPOTIFY_AUTH.name) {
                        SpotifyAuthScreen(
                            navController = navController,
                            snackbarState = snackbarState,
                            spotifyAuthViewModel = spotifyAuthViewModel
                        )
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
                            artistEditorViewModel = artistEditorViewModel,
                            venueSelectorViewModel = venueSelectorViewModel,
                            eventEditorViewModel = eventEditorViewModel
                        )
                    }
                    composable(Screen.SELECT_VENUE.name) {
                        VenueSelectorScreen(
                            navController = navController,
                            venueSelectorViewModel = venueSelectorViewModel,
                            venueEditorViewModel = venueEditorViewModel
                        )
                    }
                    composable(Screen.SELECT_ARTIST.name) {
                        ArtistSelectorScreen(
                            navController = navController,
                            artistSelectorViewModel = artistSelectorViewModel,
                            artistEditorViewModel = artistEditorViewModel
                        )
                    }
                    composable(Screen.EDIT_ARTIST.name) {
                        ArtistEditorScreen(
                            navController = navController,
                            artistEditorViewModel = artistEditorViewModel,
                            genreSelectorViewModel = genreSelectorViewModel
                        )
                    }
                    composable(Screen.SELECT_GENRE.name) {
                        GenreSelectorScreen(
                            navController = navController,
                            genreSelectorViewModel = genreSelectorViewModel
                        )
                    }
                    composable(Screen.EDIT_VENUE.name) {
                        VenueEditorScreen(
                            navController = navController,
                            venueEditorViewModel = venueEditorViewModel
                        )
                    }
                    composable(Screen.MANAGE_VENUES.name) {
                        ManageVenuesScreen(
                            navController = navController,
                            snackbarState = snackbarState,
                            venueEditorViewModel = venueEditorViewModel
                        )
                    }
                    composable(Screen.MANAGE_ARTISTS.name) {
                        ManageArtistsScreen(
                            navController = navController,
                            snackbarState = snackbarState,
                            artistEditorViewModel = artistEditorViewModel
                        )
                    }
                    composable(Screen.MANAGE_GENRES.name) {
                        ManageGenresScreen(
                            navController = navController
                        )
                    }
                    composable(Screen.MANAGE_ALBUMS.name) {
                        ManageAlbumsScreen(
                            navController = navController,
                            snackbarState = snackbarState,
                            albumEditorViewModel = albumEditorViewModel,
                            albumDetailsViewModel = albumDetailsViewModel,
                            albumsViewModel = albumsViewModel
                        )
                    }
                    composable(Screen.ALBUM_DETAILS.name) {
                        AlbumDetailsScreen(
                            navController = navController,
                            albumDetailsViewModel = albumDetailsViewModel,
                            onEdit = { album ->
                                albumEditorViewModel.launchEditor(
                                    navController = navController,
                                    album = album,
                                    onSave = { updated ->
                                        albumsViewModel.updateAlbum(
                                            album = updated,
                                            onSuccess = {
                                                albumDetailsViewModel.updateAlbum(updated)
                                                navController.popBackStack()
                                            },
                                            onError = { msg -> snackbarState.showSnackbar(msg) }
                                        )
                                    },
                                    onDelete = { toDelete ->
                                        albumsViewModel.deleteAlbum(
                                            album = toDelete,
                                            onSuccess = {
                                                navController.popBackStack(Screen.MANAGE_ALBUMS.name, false)
                                            },
                                            onError = { msg -> snackbarState.showSnackbar(msg) }
                                        )
                                    }
                                )
                            }
                        )
                    }
                    composable(Screen.EDIT_ALBUM.name) {
                        AlbumEditorScreen(
                            navController = navController,
                            albumEditorViewModel = albumEditorViewModel,
                            artistSelectorViewModel = artistSelectorViewModel,
                            genreSelectorViewModel = genreSelectorViewModel
                        )
                    }
                    composable(Screen.ANALYTICS.name) {
                        AnalyticsOverviewScreen(
                            navController = navController,
                            overviewViewModel = analyticsOverviewViewModel,
                            eventsViewModel = analyticsEventsViewModel
                        )
                    }
                    composable(Screen.ANALYTICS_YEARS.name) {
                        AnalyticsYearsScreen(
                            navController = navController,
                            listViewModel = analyticsListViewModel,
                            eventsViewModel = analyticsEventsViewModel
                        )
                    }
                    composable(Screen.ANALYTICS_MONTHS.name) {
                        AnalyticsMonthsScreen(
                            navController = navController,
                            listViewModel = analyticsListViewModel,
                            eventsViewModel = analyticsEventsViewModel
                        )
                    }
                    composable(Screen.ANALYTICS_ARTISTS.name) {
                        AnalyticsArtistsScreen(
                            navController = navController,
                            listViewModel = analyticsListViewModel,
                            eventsViewModel = analyticsEventsViewModel
                        )
                    }
                    composable(Screen.ANALYTICS_VENUES.name) {
                        AnalyticsVenuesScreen(
                            navController = navController,
                            listViewModel = analyticsListViewModel,
                            eventsViewModel = analyticsEventsViewModel
                        )
                    }
                    composable(Screen.ANALYTICS_GENRES.name) {
                        AnalyticsGenresScreen(
                            navController = navController,
                            listViewModel = analyticsListViewModel,
                            eventsViewModel = analyticsEventsViewModel
                        )
                    }
                    composable(Screen.ANALYTICS_EVENTS.name) {
                        AnalyticsEventsScreen(
                            navController = navController,
                            snackbarState = snackbarState,
                            eventsViewModel = analyticsEventsViewModel,
                            eventEditorViewModel = eventEditorViewModel
                        )
                    }
                }
            }
        }
    }
}