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
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.venue.Venue
import com.bldover.beacon.data.model.venue.VenueOrdering
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.common.VenueSearchUtilityBar
import com.bldover.beacon.ui.components.editor.VenueCard

@Composable
fun SearchableVenuesList(
    venueState: VenueState,
    orderingState: VenueOrdering,
    onSearchVenues: (String) -> Unit,
    onOrderingChange: (VenueOrdering) -> Unit,
    onVenueSelected: (Venue) -> Unit,
    onNewVenue: () -> Unit
) {
    Scaffold(
        topBar = {
            Column {
                BasicSearchBar(
                    modifier = Modifier.fillMaxWidth(),
                    enabled = venueState is VenueState.Success,
                    onQueryChange = onSearchVenues
                )
                VenueSearchUtilityBar(
                    state = orderingState,
                    onChange = onOrderingChange
                )
            }
        }
    ) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding)) {
            Spacer(modifier = Modifier.height(16.dp))
            when (venueState) {
                is VenueState.Success -> VenueList(
                    venueState.filtered,
                    onVenueSelected,
                    onNewVenue
                )
                is VenueState.Error -> LoadErrorMessage()
                is VenueState.Loading -> LoadingSpinner()
            }
        }
    }
}

@Composable
private fun VenueList(
    venues: List<Venue>,
    onVenueSelected: (Venue) -> Unit,
    onNewVenue: () -> Unit
) {
    ScrollableItemList(
        items = venues,
        topAnchor = { NewVenueCard(onNewVenue) },
        getItemKey = { it.id!! }
    ) { venue ->
        VenueCard(
            venue = venue,
            onClick = {
                onVenueSelected(venue)
            }
        )
    }
}

@Composable
private fun NewVenueCard(
    onClick: () -> Unit
) {
    BasicCard(
        modifier = Modifier.clickable { onClick() }
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