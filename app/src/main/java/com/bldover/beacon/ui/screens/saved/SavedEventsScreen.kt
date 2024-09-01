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
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.ui.components.common.AddButton
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.RefreshButton
import com.bldover.beacon.ui.components.common.SavedEventCard
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.common.TitleTopBar
import timber.log.Timber
import java.util.UUID

@Composable
fun HistoryScreen(
    navController: NavController,
    savedEventsViewModel: SavedEventsViewModel = hiltViewModel()
) {
    Timber.d("composing HistoryScreen")
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = Screen.CONCERT_HISTORY.title,
                trailingIcon = {
                    TopBarIcons(
                        navController = navController,
                        savedEventsViewModel = savedEventsViewModel
                    )
                }
            )
        }
    ) {
        SavedEventsList(
            history = true,
            navController = navController,
            savedEventsViewModel = savedEventsViewModel
        )
    }
}

@Composable
fun PlannerScreen(
    navController: NavController,
    savedEventsViewModel: SavedEventsViewModel = hiltViewModel()
) {
    Timber.d("composing PlannerScreen")
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = Screen.CONCERT_PLANNER.title,
                trailingIcon = {
                    TopBarIcons(
                        navController = navController,
                        savedEventsViewModel = savedEventsViewModel
                    )
                }
            )
        }
    ) {
        SavedEventsList(
            history = false,
            navController = navController,
            savedEventsViewModel = savedEventsViewModel
        )
    }
}

@Composable
fun TopBarIcons(
    navController: NavController,
    savedEventsViewModel: SavedEventsViewModel
) {
    Row {
        AddButton {
            val uuid = UUID.randomUUID().toString()
            navController.navigate("${Screen.EDIT_EVENT.name}/$uuid/ ")
        }
        Spacer(modifier = Modifier.padding(2.dp))
        RefreshButton { savedEventsViewModel.loadData() }
    }
}

@Composable
fun SavedEventsList(
    history: Boolean,
    navController: NavController,
    savedEventsViewModel: SavedEventsViewModel
) {
    LaunchedEffect(history) {
        if (history) savedEventsViewModel.resetPastEventFilter() else savedEventsViewModel.resetFutureEventFilter()
    }

    val eventsState = remember(history) {
        if (history) savedEventsViewModel.pastEventsState else savedEventsViewModel.futureEventsState
    }.collectAsState()

    Timber.d("composing SavedEventsList : history=$history")

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
        Timber.d("composing SavedEventsList - content")
        Box(modifier = Modifier.padding(innerPadding)) {
            SavedEventsListContent(
                eventsState = eventsState.value,
                onEventClick = remember(navController) {
                    { eventId ->
                        val uuid = UUID.randomUUID().toString()
                        navController.navigate("${Screen.EDIT_EVENT.name}/$uuid/$eventId")
                    }
                }
            )
        }
    }
}

@Composable
private fun SavedEventsListContent(
    eventsState: SavedEventsState,
    onEventClick: (String) -> Unit
) {
    Column {
        Spacer(modifier = Modifier.height(16.dp))
        when (eventsState) {
            is SavedEventsState.Success -> {
                Timber.d("composing SavedEventsList - content - success")
                ScrollableItemList(
                    items = eventsState.filtered,
                    getItemKey = { it.id!! }
                ) { event ->
                    SavedEventCard(
                        event = event,
                        onClick = { onEventClick(event.id!!) }
                    )
                }
            }
            is SavedEventsState.Error -> {
                Timber.d("composing SavedEventsList - content - error")
                LoadErrorMessage()
            }
            is SavedEventsState.Loading -> {
                Timber.d("composing SavedEventsList - content - loading")
                LoadingSpinner()
            }
        }
    }
}