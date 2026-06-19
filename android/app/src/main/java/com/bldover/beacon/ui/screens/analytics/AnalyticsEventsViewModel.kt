package com.bldover.beacon.ui.screens.analytics

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Direction
import com.bldover.beacon.data.model.OrderField
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.analytics.AnalyticsCategory
import com.bldover.beacon.data.model.event.Event
import com.bldover.beacon.data.model.event.EventOrdering
import com.bldover.beacon.data.repository.AnalyticsRepository
import com.bldover.beacon.data.repository.EventRepository
import com.bldover.beacon.ui.screens.events.SavedEventsState
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

data class AnalyticsEventsSelection(
    val category: AnalyticsCategory,
    val key: String,
    val name: String
)

@HiltViewModel
class AnalyticsEventsViewModel @Inject constructor(
    private val analyticsRepository: AnalyticsRepository,
    private val eventRepository: EventRepository
) : ViewModel() {

    private val _selection = MutableStateFlow<AnalyticsEventsSelection?>(null)
    val selection: StateFlow<AnalyticsEventsSelection?> = _selection.asStateFlow()

    private val _eventsState = MutableStateFlow<SavedEventsState>(SavedEventsState.Loading)
    val eventsState: StateFlow<SavedEventsState> = _eventsState.asStateFlow()

    private val _ordering = MutableStateFlow(EventOrdering(OrderField.DATE, Direction.DESCENDING))
    val ordering: StateFlow<EventOrdering> = _ordering.asStateFlow()

    fun launch(
        navController: NavController,
        category: AnalyticsCategory,
        key: String,
        name: String
    ) {
        _selection.value = AnalyticsEventsSelection(category, key, name)
        loadEvents()
        navController.navigate(Screen.ANALYTICS_EVENTS.name)
    }

    fun loadEvents() {
        val current = _selection.value ?: return
        Timber.i("Loading analytics events for ${current.category}=${current.key}")
        viewModelScope.launch {
            _eventsState.value = SavedEventsState.Loading
            try {
                val events = analyticsRepository.getEvents(current.category, current.key)
                val sorted = events.sortedWith(Comparator(_ordering.value::compare))
                _eventsState.value = SavedEventsState.Success(sorted, sorted)
            } catch (e: Exception) {
                Timber.e(e, "Failed to load analytics events")
                _eventsState.value = SavedEventsState.Error("Failed to load events")
            }
        }
    }

    fun filterEvents(searchTerm: String) {
        val state = _eventsState.value as? SavedEventsState.Success ?: return
        _eventsState.value = SavedEventsState.Success(
            state.allEvents,
            state.allEvents.filter { it.hasMatch(searchTerm) }
        )
    }

    fun sortEvents(ordering: EventOrdering) {
        _ordering.value = ordering
        val state = _eventsState.value as? SavedEventsState.Success ?: return
        _eventsState.value = SavedEventsState.Success(
            state.allEvents.sortedWith(Comparator(ordering::compare)),
            state.filtered.sortedWith(Comparator(ordering::compare))
        )
    }

    fun updateEvent(
        event: Event,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            Timber.i("Updating event from analytics $event")
            if (!event.isPopulated()) {
                onError("Event is missing required fields")
                return@launch
            }
            try {
                if (event.id.primary == null) eventRepository.saveEvent(event)
                else eventRepository.updateEvent(event)
                onSuccess()
                loadEvents()
            } catch (e: Exception) {
                Timber.e(e, "Failed to update event $event")
                onError("Error saving event, try again later")
            }
        }
    }

    fun deleteEvent(
        event: Event,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            Timber.i("Deleting event from analytics $event")
            try {
                eventRepository.deleteEvent(event)
                onSuccess()
                loadEvents()
            } catch (e: Exception) {
                Timber.e(e, "Failed to delete event $event")
                onError("Error deleting event, try again later")
            }
        }
    }
}
