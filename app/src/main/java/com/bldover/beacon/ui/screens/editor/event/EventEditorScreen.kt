package com.bldover.beacon.ui.screens.editor.event

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
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
import com.bldover.beacon.ui.components.editor.SwipableArtistEditCard
import com.bldover.beacon.ui.components.editor.VenueEditCard
import com.bldover.beacon.ui.screens.editor.artist.ArtistSelectorViewModel
import com.bldover.beacon.ui.screens.editor.venue.VenueSelectorViewModel
import com.bldover.beacon.ui.screens.saved.SavedEventsViewModel
import kotlinx.coroutines.launch
import timber.log.Timber
import java.time.LocalDate

@Composable
fun EventEditorScreen(
    navController: NavController,
    eventId: String? = null,
    uuid: String,
    artistSelectorViewModel: ArtistSelectorViewModel,
    venueSelectorViewModel: VenueSelectorViewModel,
    eventEditorViewModel: EventEditorViewModel,
    savedEventsViewModel: SavedEventsViewModel = hiltViewModel()
) {
    Timber.d("composing EventEditorScreen")
    LaunchedEffect(eventId, uuid) {
        eventEditorViewModel.loadEvent(eventId, uuid)
    }
    val coroutineScope = rememberCoroutineScope()
    val snackbarHostState = remember { SnackbarHostState() }
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = Screen.EDIT_EVENT.title,
                leadingIcon = { BackButton(navController = navController) }
            )
        },
        snackbarHost = { SnackbarHost(hostState = snackbarHostState) }
    ) {
        val eventState by eventEditorViewModel.uiState.collectAsState()
        Timber.d("composing EventEditorScreen - content - eventState $eventState")
        var isError = false
        when (eventState) {
            is EventEditorState.Success -> {
                Timber.d("composing EventEditorScreen - content - success")
                val event = (eventState as EventEditorState.Success).tempEvent
                LazyColumn(
                    verticalArrangement = Arrangement.spacedBy(16.dp),

                    ) {
                    val headliner = event.artists.find { it.headliner }
                    item(key = headliner?.id ?: "") {
                        if (headliner != null) {
                            SwipableArtistEditCard(
                                artist = headliner,
                                artistType = ArtistType.HEADLINER,
                                onSelect = eventEditorViewModel::updateHeadliner,
                                onSwipe = { eventEditorViewModel.updateHeadliner(null) },
                                navController = navController,
                                artistSelectorViewModel = artistSelectorViewModel
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
                    items(
                        items = openers,
                        key = { it.id ?: it.name }
                    ) { opener ->
                        SwipableArtistEditCard(
                            artist = opener,
                            artistType = ArtistType.OPENER,
                            onSwipe = eventEditorViewModel::removeOpener,
                            navController = navController,
                            artistSelectorViewModel = artistSelectorViewModel
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
                            navController = navController,
                            onChange = { eventEditorViewModel.updateVenue(it) },
                            venueSelectorViewModel = venueSelectorViewModel
                        )
                    }
                    item {
                        val futureEvent = event.date.isAfter(LocalDate.now())
                        PurchasedSwitch(
                            checked = !futureEvent || event.purchased,
                            enabled = futureEvent,
                            onChange = { eventEditorViewModel.updatePurchased(it) }
                        )
                    }
                    item {
                        Row(
                            horizontalArrangement = Arrangement.End,
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(vertical = 8.dp),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            if (eventId != null) {
                                DeleteButton(
                                    onDelete = {
                                        savedEventsViewModel.deleteEvent(
                                            event = event,
                                            onSuccess = {
                                                coroutineScope.launch {
                                                    navController.popBackStack()
                                                    snackbarHostState.showSnackbar("Event deleted")
                                                }
                                            },
                                            onError = { msg ->
                                                coroutineScope.launch {
                                                    snackbarHostState.showSnackbar(msg)
                                                }
                                            }
                                        )
                                    }
                                )
                                Spacer(modifier = Modifier.weight(1f))
                            }
                            SaveCancelButtons(
                                onCancel = { navController.popBackStack() },
                                onSave = {
                                    savedEventsViewModel.updateEvent(
                                        event = event,
                                        onSuccess = {
                                            coroutineScope.launch {
                                                navController.popBackStack()
                                                snackbarHostState.showSnackbar("Event saved")
                                            }
                                        },
                                        onError = { msg ->
                                            coroutineScope.launch {
                                                snackbarHostState.showSnackbar(msg)
                                            }
                                        }
                                    )
                                }
                            )
                        }
                    }
                }
            }
            is EventEditorState.Error -> {
                Timber.d("composing EventEditorScreen - content - error")
                isError = true
            }
            is EventEditorState.Loading -> {
                Timber.d("composing EventEditorScreen - content - loading")
                LoadingSpinner()
            }
        }
        LaunchedEffect(isError) {
            Timber.d("composing EventEditorScreen - content - isError launched effect $isError")
            if (isError) {
                coroutineScope.launch {
                    snackbarHostState.showSnackbar("Failed to load event")
                }
            }
        }
    }
}