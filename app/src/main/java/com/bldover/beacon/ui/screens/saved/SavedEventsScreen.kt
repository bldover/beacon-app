package com.bldover.beacon.ui.screens.saved

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Event
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.ui.components.common.AddButton
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.RefreshButton
import com.bldover.beacon.ui.components.common.SavedEventCard
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.screens.editor.event.EventEditorViewModel
import timber.log.Timber

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
            history = true,
            navController = navController,
            snackbarState = snackbarState,
            savedEventsViewModel = savedEventsViewModel,
            eventEditorViewModel = eventEditorViewModel
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
            history = false,
            navController = navController,
            snackbarState = snackbarState,
            savedEventsViewModel = savedEventsViewModel,
            eventEditorViewModel = eventEditorViewModel
        )
    }
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
    history: Boolean,
    navController: NavController,
    snackbarState: SnackbarState,
    savedEventsViewModel: SavedEventsViewModel,
    eventEditorViewModel: EventEditorViewModel
) {
    LaunchedEffect(history) {
        if (history) savedEventsViewModel.resetPastEventFilter() else savedEventsViewModel.resetFutureEventFilter()
    }
    val eventsState = remember(history) {
        if (history) savedEventsViewModel.pastEventsState else savedEventsViewModel.futureEventsState
    }.collectAsState()

    Scaffold(
        topBar = {
            val isEnabled = eventsState.value is SavedEventsState.Success
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = isEnabled,
                onQueryChange = remember(history) {
                    { query ->
                        if (history) savedEventsViewModel.filterPastEvents(query)
                        else savedEventsViewModel.filterFutureEvents(query)
                    }
                }
            )
        }
    ) { innerPadding ->
        Box(modifier = Modifier.padding(innerPadding)) {
            SavedEventsListContent(
                eventsState = eventsState.value,
                highlightPurchased = !history,
                onEventClick = {
                    eventEditorViewModel.launchEditor(
                        navController = navController,
                        eventId = it,
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
            )
        }
    }
}

@Composable
private fun SavedEventsListContent(
    eventsState: SavedEventsState,
    highlightPurchased: Boolean,
    onEventClick: (String) -> Unit
) {
    Column {
        Spacer(modifier = Modifier.height(16.dp))
        when (eventsState) {
            is SavedEventsState.Success -> {
                ScrollableItemList(
                    items = eventsState.filtered,
                    getItemKey = { it.id!! }
                ) { event ->
                    SavedEventCard(
                        event = event,
                        highlighted = highlightPurchased && event.purchased,
                        onClick = { onEventClick(event.id!!) }
                    )
                }
            }
            is SavedEventsState.Error -> {
                LoadErrorMessage()
            }
            is SavedEventsState.Loading -> {
                LoadingSpinner()
            }
        }
    }
}