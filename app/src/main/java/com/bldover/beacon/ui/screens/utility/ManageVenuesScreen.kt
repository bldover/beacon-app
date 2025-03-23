package com.bldover.beacon.ui.screens.utility

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.screens.editor.venue.SearchableVenuesList
import com.bldover.beacon.ui.screens.editor.venue.VenueEditorViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenuesViewModel
import timber.log.Timber

@Composable
fun ManageVenuesScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    venueEditorViewModel: VenueEditorViewModel,
    venuesViewModel: VenuesViewModel = hiltViewModel()
) {
    Timber.d("composing ManageVenuesScreen")
    LaunchedEffect(Unit) {
        venuesViewModel.resetFilter()
    }
    val venueState by venuesViewModel.uiState.collectAsState()

    ScreenFrame(
        topBar = { TitleTopBar(
            title = "Manage Venues",
            leadingIcon = { BackButton(navController = navController) }
        ) }
    ) {
        SearchableVenuesList(
            venueState = venueState,
            onSearchVenues = venuesViewModel::applyFilter,
            onVenueSelected = {
                venueEditorViewModel.launchEditor(
                    navController = navController,
                    venue = it,
                    onSave = { updated ->
                        venuesViewModel.updateVenue(
                            venue = updated,
                            onSuccess = {
                                navController.popBackStack()
                            },
                            onError = { err ->
                                Timber.e(err)
                                snackbarState.showSnackbar("Failed to save venue")
                            }
                        )
                    }
                )
            },
            onNewVenue = {
                venueEditorViewModel.launchEditor(
                    navController = navController,
                    onSave = {
                        venuesViewModel.addVenue(
                            venue = it,
                            onSuccess = {
                                navController.popBackStack()
                            },
                            onError = { err ->
                                Timber.e(err)
                                snackbarState.showSnackbar("Failed to save venue")
                            }
                        )
                    }
                )
            }
        )
    }
}