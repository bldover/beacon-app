package com.bldover.beacon.ui.screens.saved

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.Direction
import com.bldover.beacon.data.model.OrderField
import com.bldover.beacon.data.model.event.Event
import com.bldover.beacon.data.model.event.EventOrdering
import com.bldover.beacon.data.repository.EventRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class SavedEventsState {
    data object Loading : SavedEventsState()
    data class Success(
        val allEvents: List<Event>,
        val filtered: List<Event>,
    ) : SavedEventsState()
    data class Error(val message: String) : SavedEventsState()
}

@HiltViewModel
class SavedEventsViewModel @Inject constructor(
    private val eventRepository: EventRepository
) : ViewModel() {

    private val _pastEventsState = MutableStateFlow<SavedEventsState>(SavedEventsState.Loading)
    val pastEventsState: StateFlow<SavedEventsState> = _pastEventsState.asStateFlow()
    private val _futureEventsState = MutableStateFlow<SavedEventsState>(SavedEventsState.Loading)
    val futureEventsState: StateFlow<SavedEventsState> = _futureEventsState.asStateFlow()

    private val _pastEventOrdering = MutableStateFlow(EventOrdering(OrderField.DATE, Direction.DESCENDING))
    val pastEventOrdering: StateFlow<EventOrdering> = _pastEventOrdering.asStateFlow()
    private val _futureEventOrdering = MutableStateFlow(EventOrdering(OrderField.DATE, Direction.ASCENDING))
    val futureEventOrdering: StateFlow<EventOrdering> = _futureEventOrdering.asStateFlow()

    init {
        loadData()
    }

    fun loadData() {
        Timber.i("Loading saved events")
        viewModelScope.launch {
            _pastEventsState.value = SavedEventsState.Loading
            _futureEventsState.value = SavedEventsState.Loading
            try {
                val pastEvents = eventRepository.getPastSavedEvents()
                val futureEvents = eventRepository.getFutureSavedEvents()
                Timber.i("Loaded ${pastEvents.size} past saved events and ${futureEvents.size} future saved events}")
                _pastEventsState.value = SavedEventsState.Success(pastEvents, pastEvents)
                _futureEventsState.value = SavedEventsState.Success(futureEvents, futureEvents)
                sortPastEvents(_pastEventOrdering.value)
                sortFutureEvents(_futureEventOrdering.value)
            } catch (e: Exception) {
                Timber.e(e,"Failed to load saved events")
                _pastEventsState.value = SavedEventsState.Error("Failed to load events")
                _futureEventsState.value = SavedEventsState.Error("Failed to load events")
            }
        }
    }

    fun resetPastEventFilter() {
        Timber.i("Resetting past saved event filter")
        if (_pastEventsState.value !is SavedEventsState.Success) {
            return
        }
        _pastEventsState.value = SavedEventsState.Success(
            allEvents = (_pastEventsState.value as SavedEventsState.Success).allEvents,
            filtered = (_pastEventsState.value as SavedEventsState.Success).allEvents
        )
    }

    fun resetFutureEventFilter() {
        Timber.i("Resetting future saved event filter")
        if (_futureEventsState.value !is SavedEventsState.Success) {
            return
        }
        _futureEventsState.value = SavedEventsState.Success(
            allEvents = (_futureEventsState.value as SavedEventsState.Success).allEvents,
            filtered = (_futureEventsState.value as SavedEventsState.Success).allEvents
        )
    }

    fun filterPastEvents(searchTerm: String) {
        Timber.i("Filtering past saved events by $searchTerm")
        when (_pastEventsState.value) {
            is SavedEventsState.Success -> {
                val allEvents = (_pastEventsState.value as SavedEventsState.Success).allEvents
                _pastEventsState.value = SavedEventsState.Success(
                    allEvents,
                    allEvents.filter { it.hasMatch(searchTerm) }
                )
            }
            else -> {
                Timber.w("Filtering past saved events by $searchTerm - not in success state")
                return
            }
        }
    }

    fun filterFutureEvents(searchTerm: String) {
        Timber.i("Filtering future saved events by $searchTerm")
        when (_futureEventsState.value) {
            is SavedEventsState.Success -> {
                val allEvents = (_futureEventsState.value as SavedEventsState.Success).allEvents
                _futureEventsState.value = SavedEventsState.Success(
                    allEvents,
                    allEvents.filter { it.hasMatch(searchTerm) }
                )
            }
            else -> {
                Timber.w("Filtering future saved events by $searchTerm - not in success state")
                return
            }
        }
    }

    fun sortPastEvents(ordering: EventOrdering) {
        Timber.i("Sorting past saved events by $ordering")
        when (_pastEventsState.value) {
            is SavedEventsState.Success -> {
                val state = (_pastEventsState.value as SavedEventsState.Success)
                _pastEventsState.value = SavedEventsState.Success(
                    state.allEvents.sortedWith(Comparator(ordering::compare)),
                    state.filtered.sortedWith(Comparator(ordering::compare))
                )
                _pastEventOrdering.value = ordering
            }
            else -> {
                Timber.w("Sorting past events by $ordering - not in success state")
                return
            }
        }
    }

    fun sortFutureEvents(ordering: EventOrdering) {
        Timber.i("Sorting future saved events by $ordering")
        when (_futureEventsState.value) {
            is SavedEventsState.Success -> {
                val state = (_futureEventsState.value as SavedEventsState.Success)
                _futureEventsState.value = SavedEventsState.Success(
                    state.allEvents.sortedWith(Comparator(ordering::compare)),
                    state.filtered.sortedWith(Comparator(ordering::compare))
                )
                _futureEventOrdering.value = ordering
            }
            else -> {
                Timber.w("Sorting future events by $ordering - not in success state")
                return
            }
        }
    }

    fun addEvent(
        event: Event,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            Timber.i("Saving event $event")
            if (!event.isPopulated()) {
                onError("Event is missing required fields")
            } else {
                try {
                    eventRepository.saveEvent(event)
                    Timber.i("Saved event $event")
                    onSuccess()
                    loadData()
                } catch (e: Exception) {
                    Timber.e(e, "Failed to save event $event")
                    onError("Error saving event ${event.artists.first().name}, try again later")
                }
            }
        }
    }

    fun updateEvent(
        event: Event,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            Timber.i("Updating event $event")
            if (!event.isPopulated()) onError("Event is missing required fields")
            else {
                try {
                    if (event.id.primary == null) eventRepository.saveEvent(event) else eventRepository.updateEvent(event)
                    Timber.i("Updated event $event")
                    onSuccess()
                    loadData()
                } catch (e: Exception) {
                    Timber.e(e, "Failed to update event $event")
                    onError("Error saving event ${event.artists.first().name}, try again later")
                }
            }
        }
    }

    fun deleteEvent(
        event: Event,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            Timber.i("Deleting event $event")
            try {
                eventRepository.deleteEvent(event)
                Timber.i("Deleted event $event")
                onSuccess()
                loadData()
            } catch (e: Exception) {
                Timber.e(e, "Failed to delete event $event")
                onError("Error deleting event ${event.artists.first().name}, try again later")
            }
        }
    }
}