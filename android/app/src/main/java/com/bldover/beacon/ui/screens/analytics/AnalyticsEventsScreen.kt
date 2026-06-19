package com.bldover.beacon.ui.screens.analytics

import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.navigation.NavController
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.data.model.event.Event
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.screens.editor.event.EventEditorViewModel
import com.bldover.beacon.ui.screens.events.SavedEventsList
import timber.log.Timber

@Composable
fun AnalyticsEventsScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    eventsViewModel: AnalyticsEventsViewModel,
    eventEditorViewModel: EventEditorViewModel
) {
    Timber.d("composing AnalyticsEventsScreen")
    val selection by eventsViewModel.selection.collectAsState()
    val eventsState by eventsViewModel.eventsState.collectAsState()
    val ordering by eventsViewModel.ordering.collectAsState()
    val title = selection?.name ?: "Events"
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = title,
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        SavedEventsList(
            eventsState = eventsState,
            filterState = ordering,
            onSearchChange = eventsViewModel::filterEvents,
            onFilterChange = eventsViewModel::sortEvents,
            onEventClick = { eventId ->
                eventEditorViewModel.launchEditor(
                    navController = navController,
                    eventId = eventId,
                    onSave = { event: Event ->
                        eventsViewModel.updateEvent(
                            event = event,
                            onSuccess = {
                                navController.popBackStack()
                                snackbarState.showSnackbar("Event saved")
                            },
                            onError = { msg -> snackbarState.showSnackbar(msg) }
                        )
                    },
                    onDelete = { event: Event ->
                        eventsViewModel.deleteEvent(
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
