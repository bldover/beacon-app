package com.bldover.beacon.ui.screens.editor.event

import androidx.compose.foundation.lazy.items
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.ui.components.common.AddNewCard
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import com.bldover.beacon.ui.components.editor.DateEditCard
import com.bldover.beacon.ui.components.editor.SaveableEditFieldsColumn
import com.bldover.beacon.ui.components.editor.ReducedMinSizeSwitch
import com.bldover.beacon.ui.components.editor.SwipeableArtistEditCard
import com.bldover.beacon.ui.components.editor.VenueEditCard
import com.bldover.beacon.ui.screens.editor.artist.ArtistEditorViewModel
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel
import com.bldover.beacon.ui.screens.editor.artist.ArtistsViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenueSelectorViewModel
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
                SaveableEditFieldsColumn(
                    onSave = { eventEditorViewModel.onSave() },
                    onCancel = { navController.popBackStack() },
                    showDelete = eventEditorViewModel.showDelete,
                    onDelete = { eventEditorViewModel.onDelete() }
                ) {
                    val headliner = event.artists.find { it.headliner }
                    item {
                        if (headliner != null) {
                            SwipeableArtistEditCard(
                                artist = headliner,
                                label = "Headliner",
                                onSwipe = { eventEditorViewModel.updateHeadliner(null) },
                                onClick = {
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
                            AddNewCard(
                                label = "Headliner",
                                onClick = {
                                    artistSelectorViewModel.launchSelector(navController) {
                                        eventEditorViewModel.updateHeadliner(it)
                                    }
                                }
                            )
                        }
                    }
                    val openers = event.artists.filter { !it.headliner }
                    items(items = openers, key = { it.name }) { opener ->
                        SwipeableArtistEditCard(
                            artist = opener,
                            label = "Opener",
                            onSwipe = eventEditorViewModel::removeOpener,
                            onClick = {
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
                        AddNewCard(
                            label = "Opener",
                            onClick = {
                                artistSelectorViewModel.launchSelector(navController) {
                                    eventEditorViewModel.addOpener(it)
                                }
                            }
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
                    if (!event.date.isBefore(LocalDate.now())) {
                        item {
                            ReducedMinSizeSwitch(
                                label = "Purchased",
                                checked = event.purchased,
                                onChange = { eventEditorViewModel.updatePurchased(it) }
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