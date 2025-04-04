package com.bldover.beacon.ui.screens.editor.venue

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import timber.log.Timber

@Composable
fun VenueSelectorScreen(
    navController: NavController,
    venueSelectorViewModel: VenueSelectorViewModel,
    venueEditorViewModel: VenueEditorViewModel,
    venuesViewModel: VenuesViewModel = hiltViewModel()
) {
    Timber.d("composing VenueSelectorScreen")
    LaunchedEffect(Unit) {
        venuesViewModel.resetFilter()
    }
    val venueState by venuesViewModel.uiState.collectAsState()

    ScreenFrame(
        topBar = { TitleTopBar(
            title = "Select Venue",
            leadingIcon = { BackButton(navController = navController) }
        ) }
    ) {
        SearchableVenuesList(
            venueState = venueState,
            orderingState = venuesViewModel.ordering.collectAsState().value,
            onSearchVenues = venuesViewModel::applyFilter,
            onOrderingChange = venuesViewModel::applyOrdering,
            onVenueSelected = {
                venueSelectorViewModel.selectVenue(it)
                navController.popBackStack()
            },
            onNewVenue = {
                venueEditorViewModel.launchEditor(
                    navController = navController,
                    onSave = {
                        venueSelectorViewModel.selectVenue(it)
                        navController.popBackStack()
                        navController.popBackStack()
                    }
                )
            }
        )
    }
}