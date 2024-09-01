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
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.RefreshButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.common.TitleTopBar
import timber.log.Timber

@Composable
fun UpcomingScreen(
    navController: NavController,
    upcomingEventsViewModel: UpcomingEventsViewModel = hiltViewModel()
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
        UpcomingEventList(upcomingEventsViewModel)
    }
}

@Composable
fun UpcomingEventList(
    upcomingEventsViewModel: UpcomingEventsViewModel
) {
    LaunchedEffect(true) {
        upcomingEventsViewModel.resetFilter()
    }
    val upcomingState = upcomingEventsViewModel.upcomingEventsState.collectAsState()
    Scaffold(
        topBar = {
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = upcomingState.value is UpcomingEventsState.Success
            ) {
                upcomingEventsViewModel.applyFilter(it)
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
                            onClick = { /* TODO */ }
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