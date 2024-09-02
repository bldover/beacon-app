package com.bldover.beacon.ui.screens.saved

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.Event
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
                Timber.i("Loaded ${pastEvents.size} past saved events")
                Timber.i("Loaded ${pastEvents.size} future saved events")
                _pastEventsState.value = SavedEventsState.Success(pastEvents, pastEvents)
                _futureEventsState.value = SavedEventsState.Success(futureEvents, futureEvents)
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
            Timber.d("Resetting past saved event filter - not in success state")
            return
        }
        _pastEventsState.value = SavedEventsState.Success(
            allEvents = (_pastEventsState.value as SavedEventsState.Success).allEvents,
            filtered = (_pastEventsState.value as SavedEventsState.Success).allEvents
        )
        Timber.i("Resetting past saved event filter - success")
    }

    fun resetFutureEventFilter() {
        Timber.i("Resetting future saved event filter")
        if (_futureEventsState.value !is SavedEventsState.Success) {
            Timber.d("Resetting future saved event filter - not in success state")
            return
        }
        _futureEventsState.value = SavedEventsState.Success(
            allEvents = (_futureEventsState.value as SavedEventsState.Success).allEvents,
            filtered = (_futureEventsState.value as SavedEventsState.Success).allEvents
        )
        Timber.i("Resetting future saved event filter - success")
    }

    fun filterPastEvents(searchTerm: String) {
        Timber.i("Filtering past saved events by $searchTerm")
        when (_pastEventsState.value) {
            is SavedEventsState.Success -> {
                Timber.i("Filtering past saved events by $searchTerm - in success state")
                val allEvents = (_pastEventsState.value as SavedEventsState.Success).allEvents
                _pastEventsState.value = SavedEventsState.Success(
                    allEvents,
                    allEvents.filter { it.hasMatch(searchTerm) }
                )
            }
            else -> {
                Timber.d("Filtering past saved events by $searchTerm - not in success state")
                return
            }
        }
        Timber.i("Filtering past saved events by $searchTerm - success")
    }

    fun filterFutureEvents(searchTerm: String) {
        Timber.i("Filtering future saved events by $searchTerm")
        when (_futureEventsState.value) {
            is SavedEventsState.Success -> {
                Timber.i("Filtering future saved events by $searchTerm - in success state")
                val allEvents = (_futureEventsState.value as SavedEventsState.Success).allEvents
                _futureEventsState.value = SavedEventsState.Success(
                    allEvents,
                    allEvents.filter { it.hasMatch(searchTerm) }
                )
            }
            else -> {
                Timber.d("Filtering future saved events by $searchTerm - not in success state")
                return
            }
        }
        Timber.i("Filtering future saved events by $searchTerm - success")
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
                    if (event.id == null) eventRepository.saveEvent(event)
                    else eventRepository.updateEvent(event)
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

    fun isSaved(event: Event): Boolean {
        return when (futureEventsState.value) {
            is SavedEventsState.Success -> {
                val savedEvents = (futureEventsState.value as SavedEventsState.Success).allEvents
                savedEvents.any { it.ticketmasterId != null && event.ticketmasterId == it.ticketmasterId }
            }
            else -> false
        }
    }
}