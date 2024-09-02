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
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Artist
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.editor.ArtistDetailsCard
import timber.log.Timber

@Composable
fun ArtistSelectorScreen(
    navController: NavController,
    artistSelectorViewModel: ArtistSelectorViewModel,
    artistCreatorViewModel: ArtistCreatorViewModel,
    artistsViewModel: ArtistsViewModel = hiltViewModel()
) {
    Timber.d("composing ArtistSelectorScreen")
    LaunchedEffect(Unit) {
        artistsViewModel.resetFilter()
    }
    val artistState by artistsViewModel.uiState.collectAsState()

    ScreenFrame(
        topBar = { TitleTopBar(
            title = "Select Artist",
            leadingIcon = { BackButton(navController = navController) }
        ) }
    ) {
        ArtistSelectorContent(
            artistState = artistState,
            navController = navController,
            onSearchArtists = artistsViewModel::applyFilter,
            onArtistSelected = artistSelectorViewModel::selectArtist,
            onNewArtist = {
                artistCreatorViewModel.launchCreator(
                    navController = navController,
                    onSave = {
                        artistSelectorViewModel.selectArtist(it)
                        navController.popBackStack()
                    }
                )
            }
        )
    }
}

@Composable
private fun ArtistSelectorContent(
    artistState: ArtistState,
    navController: NavController,
    onSearchArtists: (String) -> Unit,
    onArtistSelected: (Artist) -> Unit,
    onNewArtist: () -> Unit
) {
    Scaffold(
        topBar = {
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = artistState is ArtistState.Success,
                onQueryChange = onSearchArtists
            )
        }
    ) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding)) {
            Spacer(modifier = Modifier.height(16.dp))
            when (artistState) {
                is ArtistState.Success -> ArtistList(artistState.filtered, navController, onArtistSelected, onNewArtist)
                is ArtistState.Error -> LoadErrorMessage()
                is ArtistState.Loading -> LoadingSpinner()
            }
        }
    }
}

@Composable
private fun ArtistList(
    artists: List<Artist>,
    navController: NavController,
    onArtistSelected: (Artist) -> Unit,
    onNewArtist: () -> Unit
) {
    ScrollableItemList(
        items = artists,
        topAnchor = { NewArtistCard(onNewArtist) },
        getItemKey = { it.id!! }
    ) { artist ->
        ArtistDetailsCard(
            artist = artist,
            onClick = {
                onArtistSelected(artist)
                navController.popBackStack()
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