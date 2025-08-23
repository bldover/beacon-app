package com.bldover.beacon.ui.screens.editor.event

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.editor.AddArtistCard
import com.bldover.beacon.ui.components.editor.ArtistType
import com.bldover.beacon.ui.components.editor.DateEditCard
import com.bldover.beacon.ui.components.editor.DeleteButton
import com.bldover.beacon.ui.components.editor.PurchasedSwitch
import com.bldover.beacon.ui.components.editor.SaveCancelButtons
import com.bldover.beacon.ui.components.editor.SwipeableArtistEditCard
import com.bldover.beacon.ui.components.editor.VenueEditCard
import com.bldover.beacon.ui.screens.editor.artist.ArtistEditorViewModel
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel
import com.bldover.beacon.ui.screens.editor.artist.ArtistsViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenueSelectorViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenuesViewModel
import timber.log.Timber
import java.time.LocalDate

@Composable
fun EventEditorScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    artistSelectorViewModel: ArtistSelectorViewModel,
    artistEditorViewModel: ArtistEditorViewModel,
    venueSelectorViewModel: VenueSelectorViewModel,
    eventEditorViewModel: EventEditorViewModel,
    artistsViewModel: ArtistsViewModel = hiltViewModel()
) {
    Timber.d("composing EventEditorScreen")
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = Screen.EDIT_EVENT.title,
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        val eventState by eventEditorViewModel.uiState.collectAsState()
        var isError = false
        when (eventState) {
            is EventEditorState.Success -> {
                val event = (eventState as EventEditorState.Success).event
                LazyColumn(verticalArrangement = Arrangement.spacedBy(16.dp)) {
                    val headliner = event.artists.find { it.headliner }
                    item {
                        if (headliner != null) {
                            SwipeableArtistEditCard(
                                artist = headliner,
                                artistType = ArtistType.HEADLINER,
                                onSwipe = { eventEditorViewModel.updateHeadliner(null) },
                                onSelect = {
                                    artistEditorViewModel.launchEditor(
                                        navController = navController,
                                        artist = headliner,
                                        onSave = { updated ->
                                            artistsViewModel.upsertArtist(
                                                artist = updated,
                                                onSuccess = {
                                                    eventEditorViewModel.updateHeadliner(it)
                                                    navController.popBackStack()
                                                },
                                                onError = { err ->
                                                    Timber.e(err)
                                                    snackbarState.showSnackbar("Failed to save artist")
                                                }
                                            )
                                        }
                                    )
                                }
                            )
                        } else {
                            AddArtistCard(
                                artistType = ArtistType.HEADLINER,
                                onSelect = eventEditorViewModel::updateHeadliner,
                                navController = navController,
                                artistSelectorViewModel = artistSelectorViewModel
                            )
                        }
                    }
                    val openers = event.artists.filter { !it.headliner }
                    items(items = openers) { opener ->
                        SwipeableArtistEditCard(
                            artist = opener,
                            artistType = ArtistType.OPENER,
                            onSwipe = eventEditorViewModel::removeOpener,
                            onSelect = {
                                artistEditorViewModel.launchEditor(
                                    navController = navController,
                                    artist = opener,
                                    onSave = { updated ->
                                        artistsViewModel.upsertArtist(
                                            artist = updated,
                                            onSuccess = {
                                                eventEditorViewModel.updateOpener(opener, it)
                                                navController.popBackStack()
                                            },
                                            onError = { err ->
                                                Timber.e(err)
                                                snackbarState.showSnackbar("Failed to save artist")
                                            }
                                        )
                                    }
                                )
                            }
                        )
                    }
                    item {
                        AddArtistCard(
                            artistType = ArtistType.OPENER,
                            onSelect = eventEditorViewModel::addOpener,
                            navController = navController,
                            artistSelectorViewModel = artistSelectorViewModel
                        )
                    }
                    item {
                        DateEditCard(
                            date = event.date,
                            onChange = { eventEditorViewModel.updateDate(it) }
                        )
                    }
                    item {
                        VenueEditCard(
                            venue = event.venue,
                            onClick = {
                                venueSelectorViewModel.launchSelector(navController) {
                                    eventEditorViewModel.updateVenue(it)
                                }
                            }
                        )
                    }
                    item {
                        val futureEvent = !event.date.isBefore(LocalDate.now())
                        PurchasedSwitch(
                            checked = !futureEvent || event.purchased,
                            enabled = futureEvent,
                            onChange = { eventEditorViewModel.updatePurchased(it) }
                        )
                    }
                    item {
                        Row(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(vertical = 8.dp),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            if (eventEditorViewModel.showDelete) {
                                DeleteButton(onDelete = { eventEditorViewModel.onDelete() })
                            }
                            Spacer(modifier = Modifier.weight(1f))
                            SaveCancelButtons(
                                onCancel = { navController.popBackStack() },
                                onSave = { eventEditorViewModel.onSave() }
                            )
                        }
                    }
                }
            }
            is EventEditorState.Error -> {
                isError = true
            }
            is EventEditorState.Loading -> {
                LoadingSpinner()
            }
        }
        LaunchedEffect(isError) {
            if (isError) { snackbarState.showSnackbar("Failed to load event") }
        }
    }
}