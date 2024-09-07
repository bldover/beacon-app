package com.bldover.beacon.ui.screens.editor.event

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Artist
import com.bldover.beacon.data.model.Event
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.Venue
import com.bldover.beacon.data.repository.ArtistRepository
import com.bldover.beacon.data.repository.EventRepository
import com.bldover.beacon.data.repository.VenueRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.async
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import java.time.LocalDate
import javax.inject.Inject

sealed class EventEditorState {
    data object Loading : EventEditorState()
    data class Success(val event: Event) : EventEditorState()
    data class Error(val message: String) : EventEditorState()
}

@HiltViewModel
class EventEditorViewModel @Inject constructor(
    private val eventRepository: EventRepository,
    private val artistRepository: ArtistRepository,
    private val venueRepository: VenueRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow<EventEditorState>(EventEditorState.Loading)
    val uiState: StateFlow<EventEditorState> = _uiState.asStateFlow()

    private var onSave: (Event) -> Unit = {}
    private var onDelete: (Event) -> Unit = {}
    var showDelete: Boolean = false
        private set

    fun launchEditor(
        navController: NavController,
        eventId: String,
        onSave: (Event) -> Unit,
        onDelete: (Event) -> Unit
    ) {
        loadEventId(eventId)
        this.onSave = onSave
        this.onDelete = onDelete
        this.showDelete = true
        navController.navigate(Screen.EDIT_EVENT.name)
    }

    fun launchEditor(
        navController: NavController,
        onSave: (Event) -> Unit
    ) {
        loadDefaultEvent()
        this.onSave = onSave
        this.onDelete = {}
        this.showDelete = false
        navController.navigate(Screen.EDIT_EVENT.name)
    }

    fun launchEditor(
        navController: NavController,
        event: Event,
        onSave: (Event) -> Unit
    ) {
        viewModelScope.launch {
            val deferredArtists = async { artistRepository.getArtists() }
            val deferredVenue = async { venueRepository.getVenues() }
            val savedArtists = deferredArtists.await()
            val savedVenues = deferredVenue.await()
            event.artists = replaceArtistsWithSaved(event.artists, savedArtists)
            event.venue = replaceVenueWithSaved(event.venue, savedVenues)
            _uiState.value = EventEditorState.Success(event)
        }
        this.onSave = onSave
        this.onDelete = {}
        this.showDelete = false
        navController.navigate(Screen.EDIT_EVENT.name)
    }

    private fun loadDefaultEvent() {
        val event = Event(
            artists = emptyList(),
            date = LocalDate.now(),
            venue = Venue(name = "", city = "", state = ""),
            purchased = false
        )
        _uiState.value = EventEditorState.Success(event)
    }

    private fun loadEventId(eventId: String) {
        Timber.i("Loading edit event ID $eventId")
        viewModelScope.launch {
            _uiState.value = EventEditorState.Loading
            try {
                val eventReq = async { eventRepository.getEvent(eventId) }
                val deferredArtists = async { artistRepository.getArtists() }
                val deferredVenues = async { venueRepository.getVenues() }
                val event = eventReq.await()
                val savedArtists = deferredArtists.await()
                val savedVenues = deferredVenues.await()
                event.artists = replaceArtistsWithSaved(event.artists, savedArtists)
                event.venue = replaceVenueWithSaved(event.venue, savedVenues)
                _uiState.value = EventEditorState.Success(event.copy())
            } catch (e: Exception) {
                Timber.e(e,"Failed to load event $eventId")
                _uiState.value = EventEditorState.Error("Failed to load event")
            }
        }
    }

    private fun replaceArtistsWithSaved(artists: List<Artist>, savedArtists: List<Artist>): List<Artist> {
        return artists.map { artist ->
            val saved = savedArtists.find { it.name.equals(artist.name, ignoreCase = true) }
            saved?.let {
                it.headliner = artist.headliner
                it
            } ?: artist }
    }

    private fun replaceVenueWithSaved(venue: Venue, savedVenues: List<Venue>): Venue {
        return savedVenues.find { it.name.equals(venue.name, ignoreCase = true)
                && it.city.equals(venue.city, ignoreCase = true)} ?: venue
    }

    fun updateHeadliner(headliner: Artist?) {
        Timber.i("Updating headliner $headliner")
        if (_uiState.value !is EventEditorState.Success) {
            Timber.d("Updating headliner - not in success state")
            return
        }
        val state = (_uiState.value as EventEditorState.Success)
        Timber.d("Updating headliner - previous artists ${state.event.artists}")
        val artists = state.event.artists.toMutableList().apply {
            removeAll { it.headliner }
            headliner?.let {
                it.headliner = true
                add(it)
            }
        }
        Timber.d("Updating headliner - new artists $artists")
        _uiState.value = EventEditorState.Success(state.event.copy(artists = artists))
        Timber.i("Updated headliner - success")
    }

    fun addOpener(opener: Artist) {
        Timber.i("Adding opener $opener")
        if (_uiState.value !is EventEditorState.Success) {
            Timber.d("Updating openers - not in success state")
            return
        }
        val state = (_uiState.value as EventEditorState.Success)
        Timber.d("Adding opener - previous artists ${state.event.artists}")
        val artists = state.event.artists.toMutableList().apply {
            add(opener)
        }
        Timber.d("Adding opener - new artists $artists")
        _uiState.value = EventEditorState.Success(state.event.copy(artists = artists))
        Timber.i("Adding opener - success")
    }

    fun updateOpener(opener: Artist, newOpener: Artist) {
        Timber.i("Updating opener $opener to $newOpener")
        if (_uiState.value !is EventEditorState.Success) {
            Timber.d("Updating openers - not in success state")
            return
        }
        val state = (_uiState.value as EventEditorState.Success)
        Timber.d("Updating opener - previous artists ${state.event.artists}")
        val artists = state.event.artists.toMutableList().apply {
            val i = indexOf(opener)
            if (i == -1) {
                Timber.w("Updating opener - opener $opener not found")
                return
            }
            removeAt(i)
            add(i, newOpener)
        }
        Timber.d("Updating opener - new artists $artists")
        _uiState.value = EventEditorState.Success(state.event.copy(artists = artists))
        Timber.i("Updating opener - success")
    }

    fun removeOpener(opener: Artist) {
        Timber.i("Removing opener $opener")
        if (_uiState.value !is EventEditorState.Success) {
            Timber.d("Updating openers - not in success state")
            return
        }
        val state = (_uiState.value as EventEditorState.Success)
        Timber.d("Removing opener - previous artists ${state.event.artists}")
        val artists = state.event.artists.toMutableList().apply {
            remove(opener)
        }
        Timber.d("Removing opener - new artists $artists")
        _uiState.value = EventEditorState.Success(state.event.copy(artists = artists))
        Timber.i("Removing opener - success")
    }

    fun updateVenue(venue: Venue) {
        Timber.i("Updating venue $venue")
        if (_uiState.value !is EventEditorState.Success) {
            Timber.d("Updating venue - not in success state")
            return
        }
        val state = (_uiState.value as EventEditorState.Success)
        Timber.d("Updating venue - previous venue ${state.event.venue}")
        _uiState.value = EventEditorState.Success(state.event.copy(venue = venue))
        Timber.i("Updated venue - success")
    }

    fun updateDate(date: LocalDate) {
        Timber.i("Updating date $date")
        if (_uiState.value !is EventEditorState.Success) {
            Timber.d("Updating date - not in success state")
            return
        }
        val state = (_uiState.value as EventEditorState.Success)
        Timber.d("Updating date - previous date ${state.event.date}")
        _uiState.value = EventEditorState.Success(state.event.copy(date = date))
        Timber.i("Updated date - success")
    }

    fun updatePurchased(purchased: Boolean) {
        Timber.i("Updating purchased $purchased")
        if (_uiState.value !is EventEditorState.Success) {
            Timber.d("Updating purchased - not in success state")
            return
        }
        val state = (_uiState.value as EventEditorState.Success)
        Timber.d("Updating purchased - previous purchased ${state.event.purchased}")
        _uiState.value = EventEditorState.Success(state.event.copy(purchased = purchased))
        Timber.i("Updated purchased - success")
    }

    fun onSave() {
        onSave((_uiState.value as EventEditorState.Success).event)
    }

    fun onDelete() {
        onDelete((_uiState.value as EventEditorState.Success).event)
    }
}