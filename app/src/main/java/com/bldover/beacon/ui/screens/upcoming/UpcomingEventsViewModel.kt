package com.bldover.beacon.ui.screens.upcoming

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.EventDetail
import com.bldover.beacon.data.model.RecommendationThreshold
import com.bldover.beacon.data.repository.EventRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class UpcomingEventsState {
    data object Loading : UpcomingEventsState()
    data class Success(
        val events: List<EventDetail>,
        val filtered: List<EventDetail>
    ) : UpcomingEventsState()
    data class Error(val message: String) : UpcomingEventsState()
}

@HiltViewModel
class UpcomingEventsViewModel @Inject constructor(
    private val eventRepository: EventRepository
) : ViewModel() {

    private val _upcomingEventsState = MutableStateFlow<UpcomingEventsState>(UpcomingEventsState.Loading)
    val upcomingEventsState: StateFlow<UpcomingEventsState> = _upcomingEventsState.asStateFlow()

    private val _filterState = MutableStateFlow(RecommendationThreshold.NONE)
    val filterState: StateFlow<RecommendationThreshold> = _filterState.asStateFlow()

    init {
        loadData()
    }

    fun loadData() {
        Timber.i("Loading planned events for threshold $_filterState.value")
        viewModelScope.launch {
            _upcomingEventsState.value = UpcomingEventsState.Loading
            try {
                val events = when (_filterState.value) {
                    RecommendationThreshold.NONE -> eventRepository.getUpcomingEvents()
                    else -> eventRepository.getRecommendedEvents(_filterState.value)
                }
                Timber.i("Loaded ${events.size} planned events")
                _upcomingEventsState.value = UpcomingEventsState.Success(events, events)
            } catch (e: Exception) {
                Timber.e(e,"Failed to load planned events")
                _upcomingEventsState.value = UpcomingEventsState.Error("Failed to load events")
            }
        }
    }

    fun resetFilter() {
        Timber.i("Resetting upcoming event filter")
        if (_upcomingEventsState.value !is UpcomingEventsState.Success) {
            Timber.d("Resetting upcoming event filter - not in success state")
            return
        }
        _upcomingEventsState.value = UpcomingEventsState.Success(
            events = (_upcomingEventsState.value as UpcomingEventsState.Success).events,
            filtered = (_upcomingEventsState.value as UpcomingEventsState.Success).events
        )
        Timber.i("Reset upcoming event filter - success")
    }

    fun applyFilter(searchTerm: String) {
        Timber.i("Applying upcoming event filter for $searchTerm")
        when (_upcomingEventsState.value) {
            is UpcomingEventsState.Success -> {
                Timber.d("Applying upcoming event filter - in success state")
                val allEvents = (_upcomingEventsState.value as UpcomingEventsState.Success).events
                _upcomingEventsState.value = UpcomingEventsState.Success(
                    allEvents,
                    allEvents.filter { it.hasMatch(searchTerm) }
                )
            }
            else -> {
                Timber.d("Applying upcoming event filter - not in success state")
                return
            }
        }
        Timber.i("Applied upcoming event filter - success")
    }

    fun changeRecommendationThreshold(threshold: RecommendationThreshold) {
        Timber.i("Changing upcoming event recommendation threshold to $threshold")
        _filterState.value = threshold
        loadData()
    }
}