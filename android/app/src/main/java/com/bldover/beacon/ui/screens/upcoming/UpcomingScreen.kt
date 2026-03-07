package com.bldover.beacon.ui.screens.upcoming

import androidx.compose.foundation.ExperimentalFoundationApi
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.data.model.event.Event
import com.bldover.beacon.data.model.event.EventDetail
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.RecommendationSelectionBar
import com.bldover.beacon.ui.components.common.RefreshButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.common.UpcomingEventCard
import com.bldover.beacon.ui.screens.editor.event.EventEditorViewModel
import com.bldover.beacon.ui.screens.saved.SavedEventsViewModel
import timber.log.Timber
import java.time.LocalDate
import java.time.format.DateTimeFormatter

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

@OptIn(ExperimentalFoundationApi::class)
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
            when (upcomingState.value) {
                is UpcomingEventsState.Success -> {
                    val events = (upcomingState.value as UpcomingEventsState.Success).filtered
                    val groupedEvents = events
                        .groupBy { it.date }
                        .mapValues { (_, eventsForDate) ->
                            eventsForDate.sortedWith(
                                compareByDescending<EventDetail> { it.purchased && it.id.primary != null }
                                    .thenByDescending { it.id.primary != null }
                                    .thenByDescending { it.rank ?: Float.MIN_VALUE }
                            )
                        }
                        .toList()
                        .sortedBy { it.first }

                    LazyColumn(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(innerPadding),
                        verticalArrangement = Arrangement.spacedBy(16.dp)
                    ) {
                        groupedEvents.forEach { (date, eventsForDate) ->
                            stickyHeader {
                                DateHeader(date = date)
                            }

                            items(
                                items = eventsForDate,
                                key = { event -> event.uniqueId() }
                            ) { event ->
                                UpcomingEventCard(
                                    event = event,
                                    accented = event.isSaved(),
                                    onClick = {
                                        eventEditorViewModel.launchEditor(
                                            navController = navController,
                                            event = event.asEvent(),
                                            onSave = { event: Event ->
                                                savedEventsViewModel.addEvent(
                                                    event = event,
                                                    onSuccess = {
                                                        navController.popBackStack()
                                                        snackbarState.showSnackbar("Event saved")
                                                    },
                                                    onError = { msg ->
                                                        snackbarState.showSnackbar(
                                                            msg
                                                        )
                                                    }
                                                )
                                            }
                                        )
                                    }
                                )
                            }
                        }
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

@Composable
fun DateHeader(date: LocalDate) {
    Text(
        text = date.format(DateTimeFormatter.ofPattern("MMMM d")),
        style = MaterialTheme.typography.titleMedium,
        modifier = Modifier
            .fillMaxWidth()
            .background(MaterialTheme.colorScheme.background)
            .padding(horizontal = 16.dp, vertical = 8.dp)
    )
}