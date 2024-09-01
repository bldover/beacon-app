package com.bldover.beacon.ui.screens.editor.venue

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AddCircle
import androidx.compose.material3.Icon
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Venue
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.editor.VenueCard
import timber.log.Timber

@Composable
fun VenueSelectorScreen(
    navController: NavController,
    venueSelectorViewModel: VenueSelectorViewModel,
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
        VenueSelectorContent(
            venueState = venueState,
            navController = navController,
            onSearchVenues = venuesViewModel::applyFilter,
            onVenueSelected = venueSelectorViewModel::selectVenue
        )
    }
}

@Composable
private fun VenueSelectorContent(
    venueState: VenueState,
    navController: NavController,
    onSearchVenues: (String) -> Unit,
    onVenueSelected: (Venue) -> Unit
) {
    Scaffold(
        topBar = {
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = venueState is VenueState.Success,
                onQueryChange = onSearchVenues
            )
        }
    ) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding)) {
            Spacer(modifier = Modifier.height(16.dp))
            when (venueState) {
                is VenueState.Success -> VenueList(venueState.filtered, navController, onVenueSelected)
                is VenueState.Error -> LoadErrorMessage()
                is VenueState.Loading -> LoadingSpinner()
            }
        }
    }
}

@Composable
private fun VenueList(
    venues: List<Venue>,
    navController: NavController,
    onVenueSelected: (Venue) -> Unit
) {
    ScrollableItemList(
        items = venues,
        topAnchor = { NewVenueCard() },
        getItemKey = { it.id!! }
    ) { venue ->
        VenueCard(
            venue = venue,
            onClick = {
                onVenueSelected(venue)
                navController.popBackStack()
            }
        )
    }
}

@Composable
private fun NewVenueCard() {
    BasicCard(
        modifier = Modifier.clickable { /* TODO: new venue */ }
    ) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
        ) {
            Text(text = "New Venue")
            Icon(
                imageVector = Icons.Default.AddCircle,
                contentDescription = "New Venue"
            )
        }
    }
}