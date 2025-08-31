package com.bldover.beacon.ui.screens.upcoming

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.data.model.event.Event
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.RecommendationSelectionBar
import com.bldover.beacon.ui.components.common.RefreshButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.common.UpcomingEventCard
import com.bldover.beacon.ui.screens.editor.event.EventEditorViewModel
import com.bldover.beacon.ui.screens.saved.SavedEventsViewModel
import timber.log.Timber

@Composable
fun UpcomingScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    eventEditorViewModel: EventEditorViewModel,
    savedEventsViewModel: SavedEventsViewModel,
    upcomingEventsViewModel: UpcomingEventsViewModel
) {
    Timber.d("composing UpcomingScreen")
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = Screen.UPCOMING_EVENTS.title,
                trailingIcon = { RefreshButton { upcomingEventsViewModel.loadData() } }
            )
        }
    ) {
        UpcomingEventList(
            navController = navController,
            snackbarState = snackbarState,
            eventEditorViewModel = eventEditorViewModel,
            savedEventsViewModel = savedEventsViewModel,
            upcomingEventsViewModel = upcomingEventsViewModel
        )
    }
}

@Composable
fun UpcomingEventList(
    navController: NavController,
    snackbarState: SnackbarState,
    eventEditorViewModel: EventEditorViewModel,
    savedEventsViewModel: SavedEventsViewModel,
    upcomingEventsViewModel: UpcomingEventsViewModel
) {
    LaunchedEffect(true) {
        upcomingEventsViewModel.resetFilter()
    }
    val upcomingState = upcomingEventsViewModel.upcomingEventsState.collectAsState()
    Scaffold(
        topBar = {
            Column {
                BasicSearchBar(
                    modifier = Modifier.fillMaxWidth(),
                    enabled = upcomingState.value is UpcomingEventsState.Success
                ) {
                    upcomingEventsViewModel.applyFilter(it)
                }
                RecommendationSelectionBar(
                    state = upcomingEventsViewModel.filterState.collectAsState().value,
                    onChange = { upcomingEventsViewModel.changeRecommendationThreshold(it) }
                )
            }
        }
    ) { innerPadding ->
        Column {
            Spacer(modifier = Modifier.height(16.dp))
            when (upcomingState.value) {
                is UpcomingEventsState.Success -> {
                    ScrollableItemList(
                        items = (upcomingState.value as UpcomingEventsState.Success).filtered,
                        modifier = Modifier.padding(innerPadding)
                    ) {
                        UpcomingEventCard(
                            event = it,
                            accented = it.id.primary != null,
                            onClick = {
                                eventEditorViewModel.launchEditor(
                                    navController = navController,
                                    event = it.asEvent(),
                                    onSave = { event: Event ->
                                        savedEventsViewModel.addEvent(
                                            event = event,
                                            onSuccess = {
                                                navController.popBackStack()
                                                snackbarState.showSnackbar("Event saved")
                                            },
                                            onError = { msg -> snackbarState.showSnackbar(msg) }
                                        )
                                    }
                                )
                            }
                        )
                    }
                }
                is UpcomingEventsState.Error -> {
                    LoadErrorMessage()
                }
                is UpcomingEventsState.Loading -> {
                    LoadingSpinner()
                }
            }
        }
    }
}