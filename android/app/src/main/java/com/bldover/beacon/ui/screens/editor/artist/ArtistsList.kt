package com.bldover.beacon.ui.screens.editor.artist

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
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.artist.ArtistOrdering
import com.bldover.beacon.ui.components.common.ArtistSearchUtilityBar
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.editor.ArtistDetailsCard

@Composable
fun SearchableArtistsList(
    artistState: ArtistState,
    showAnyGenre: Boolean = false,
    orderingState: ArtistOrdering,
    onSearchArtists: (String) -> Unit,
    onOrderingChange: (ArtistOrdering) -> Unit,
    onArtistSelected: (Artist) -> Unit,
    onNewArtist: () -> Unit
) {
    Scaffold(
        topBar = {
            Column {
                BasicSearchBar(
                    modifier = Modifier.fillMaxWidth(),
                    enabled = artistState is ArtistState.Success,
                    onQueryChange = onSearchArtists
                )
                ArtistSearchUtilityBar(
                    state = orderingState,
                    onChange = onOrderingChange
                )
            }
        }
    ) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding)) {
            Spacer(modifier = Modifier.height(16.dp))
            when (artistState) {
                is ArtistState.Success -> ArtistList(
                    artistState.filtered,
                    showAnyGenre,
                    onArtistSelected,
                    onNewArtist
                )
                is ArtistState.Error -> LoadErrorMessage()
                is ArtistState.Loading -> LoadingSpinner()
            }
        }
    }
}

@Composable
private fun ArtistList(
    artists: List<Artist>,
    showAnyGenre: Boolean = false,
    onArtistSelected: (Artist) -> Unit,
    onNewArtist: () -> Unit
) {
    ScrollableItemList(
        items = artists,
        topAnchor = { NewArtistCard(onNewArtist) },
        getItemKey = { it.id.primary!! }
    ) { artist ->
        ArtistDetailsCard(
            artist = artist,
            showAnyGenre = showAnyGenre,
            onClick = {
                onArtistSelected(artist)
            }
        )
    }
}

@Composable
private fun NewArtistCard(
    onNewArtist: () -> Unit
) {
    BasicCard(
        modifier = Modifier.clickable { onNewArtist() }
    ) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
        ) {
            Text(text = "New Artist")
            Icon(
                imageVector = Icons.Default.AddCircle,
                contentDescription = "New Artist"
            )
        }
    }
}