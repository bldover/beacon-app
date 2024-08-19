package com.bldover.beacon.ui.screens.upcoming

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.bldover.beacon.ui.components.BasicSearchBar
import com.bldover.beacon.ui.components.EventDetailCard
import com.bldover.beacon.ui.components.LoadErrorMessage
import com.bldover.beacon.ui.components.LoadingSpinner
import com.bldover.beacon.ui.components.ScrollableItemList

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun UpcomingScreen(
    upcomingViewModel: UpcomingViewModel = hiltViewModel()
) {
    val plannerState = upcomingViewModel.uiState.collectAsState()
    Scaffold(
        topBar = {
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = plannerState.value is UiState.Success
            ) {
                upcomingViewModel.handleEvent(UiEvent.ApplySearchFilter(it))
            }
        }
    ) { innerPadding ->
        Column {
            Spacer(modifier = Modifier.height(16.dp))
            when (plannerState.value) {
                is UiState.Success -> {
                    ScrollableItemList(
                        items = (plannerState.value as UiState.Success).filtered,
                        modifier = Modifier.padding(innerPadding)
                    ) {
                        EventDetailCard(it)
                    }
                }

                is UiState.Error -> {
                    LoadErrorMessage()
                }

                is UiState.Loading -> {
                    LoadingSpinner()
                }
            }
        }
    }
}