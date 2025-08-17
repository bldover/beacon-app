package com.bldover.beacon.ui.screens.saved

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.event.Event
import com.bldover.beacon.data.model.event.EventOrdering
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.data.model.artist.GenreInfo
import com.bldover.beacon.data.model.venue.Venue
import com.bldover.beacon.ui.components.common.AddButton
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.EventSearchUtilityBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.RefreshButton
import com.bldover.beacon.ui.components.common.SavedEventCard
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.screens.editor.event.EventEditorViewModel
import timber.log.Timber
import java.time.LocalDate

@Composable
fun HistoryScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    savedEventsViewModel: SavedEventsViewModel,
    eventEditorViewModel: EventEditorViewModel
) {
    Timber.d("composing HistoryScreen")
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = Screen.CONCERT_HISTORY.title,
                trailingIcon = {
                    TopBarIcons(
                        navController = navController,
                        snackbarState = snackbarState,
                        savedEventsViewModel = savedEventsViewModel,
                        eventEditorViewModel = eventEditorViewModel
                    )
                }
            )
        }
    ) {
        SavedEventsList(
            eventsState = savedEventsViewModel.pastEventsState.collectAsState().value,
            filterState = savedEventsViewModel.pastEventOrdering.collectAsState().value,
            accentPurchased = false,
            onSearchChange = savedEventsViewModel::filterPastEvents,
            onFilterChange = savedEventsViewModel::sortPastEvents,
            onEventClick = {
                launchEventEditor(
                    eventId = it,
                    eventEditorViewModel = eventEditorViewModel,
                    savedEventsViewModel = savedEventsViewModel,
                    navController = navController,
                    snackbarState = snackbarState
                )
            }
        )
    }
}

@Composable
fun PlannerScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    savedEventsViewModel: SavedEventsViewModel,
    eventEditorViewModel: EventEditorViewModel
) {
    Timber.d("composing PlannerScreen")
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = Screen.CONCERT_PLANNER.title,
                trailingIcon = {
                    TopBarIcons(
                        navController = navController,
                        snackbarState = snackbarState,
                        savedEventsViewModel = savedEventsViewModel,
                        eventEditorViewModel = eventEditorViewModel
                    )
                }
            )
        }
    ) {
        SavedEventsList(
            eventsState = savedEventsViewModel.futureEventsState.collectAsState().value,
            filterState = savedEventsViewModel.futureEventOrdering.collectAsState().value,
            accentPurchased = true,
            onSearchChange = { savedEventsViewModel.filterFutureEvents(it) },
            onFilterChange = { savedEventsViewModel.sortFutureEvents(it) },
            onEventClick = {
                launchEventEditor(
                    eventId = it,
                    eventEditorViewModel = eventEditorViewModel,
                    savedEventsViewModel = savedEventsViewModel,
                    navController = navController,
                    snackbarState = snackbarState
                )
            }
        )
    }
}

private fun launchEventEditor(
    eventId: String,
    eventEditorViewModel: EventEditorViewModel,
    savedEventsViewModel: SavedEventsViewModel,
    navController: NavController,
    snackbarState: SnackbarState
) {
    eventEditorViewModel.launchEditor(
        navController = navController,
        eventId = eventId,
        onSave = { event: Event ->
            savedEventsViewModel.updateEvent(
                event = event,
                onSuccess = {
                    navController.popBackStack()
                    snackbarState.showSnackbar("Event saved")
                },
                onError = { msg -> snackbarState.showSnackbar(msg) }
            )
        },
        onDelete = { event: Event ->
            savedEventsViewModel.deleteEvent(
                event = event,
                onSuccess = {
                    navController.popBackStack()
                    snackbarState.showSnackbar("Event deleted")
                },
                onError = { msg -> snackbarState.showSnackbar(msg) }
            )
        }
    )
}

@Composable
fun TopBarIcons(
    navController: NavController,
    snackbarState: SnackbarState,
    savedEventsViewModel: SavedEventsViewModel,
    eventEditorViewModel: EventEditorViewModel
) {
    Row {
        AddButton {
            eventEditorViewModel.launchEditor(
                navController = navController,
                onSave = {
                    savedEventsViewModel.addEvent(
                        it,
                        onSuccess = {
                            navController.popBackStack()
                            snackbarState.showSnackbar("Event saved")
                        },
                        onError = { msg -> snackbarState.showSnackbar(msg) }
                    )
                }
            )
        }
        Spacer(modifier = Modifier.padding(2.dp))
        RefreshButton { savedEventsViewModel.loadData() }
    }
}

@Composable
fun SavedEventsList(
    eventsState: SavedEventsState,
    filterState: EventOrdering,
    accentPurchased: Boolean,
    onSearchChange: (String) -> Unit,
    onFilterChange: (EventOrdering) -> Unit,
    onEventClick: (String) -> Unit
) {
    Scaffold(
        topBar = {
            Column {
                BasicSearchBar(
                    modifier = Modifier.fillMaxWidth(),
                    enabled = eventsState is SavedEventsState.Success,
                    onQueryChange = onSearchChange
                )
                EventSearchUtilityBar(
                    state = filterState,
                    onChange = onFilterChange
                )
            }
        }
    ) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding)) {
            Spacer(modifier = Modifier.padding(8.dp))
            when (eventsState) {
                is SavedEventsState.Success -> {
                    ScrollableItemList(
                        items = eventsState.filtered
                    ) { event ->
                        SavedEventCard(
                            event = event,
                            accented = accentPurchased && event.purchased,
                            onClick = { onEventClick(event.id!!) }
                        )
                    }
                }
                is SavedEventsState.Error -> { LoadErrorMessage() }
                is SavedEventsState.Loading -> { LoadingSpinner() }
            }
        }
    }
}

@Preview
@Composable
fun SavedEventsListPreview() {
    val events = listOf(
        Event(
            id = "1",
            artists = listOf(Artist("123", "Test Artist", GenreInfo(listOf("Test Genre"), emptyList(), emptyList()))),
            venue = Venue("123", "Test Venue", "Test City", "Test State"),
            date = LocalDate.now(),
            purchased = false
        )
    )
    SavedEventsList(
        eventsState = SavedEventsState.Success(
            allEvents = events,
            filtered = events
        ),
        filterState = EventOrdering(),
        accentPurchased = false,
        onSearchChange = {},
        onFilterChange = {},
        onEventClick = {}
    )
}