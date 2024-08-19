package com.bldover.beacon.ui.screens.upcoming

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.EventDetail
import com.bldover.beacon.data.repository.EventRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class UiState {
    data object Loading : UiState()
    data class Success(
        val events: List<EventDetail>,
        val filtered: List<EventDetail>
    ) : UiState()
    data class Error(val message: String) : UiState()
}

sealed class UiEvent {
    data object RefreshData : UiEvent()
    data class ApplySearchFilter(val searchTerm: String) : UiEvent()
}

@HiltViewModel
class UpcomingViewModel @Inject constructor(
    private val eventRepository: EventRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow<UiState>(UiState.Loading)
    val uiState: StateFlow<UiState> = _uiState.asStateFlow()

    init {
        loadData()
    }

    fun handleEvent(event: UiEvent) {
        when (event) {
            is UiEvent.RefreshData -> loadData()
            is UiEvent.ApplySearchFilter -> updateFilter(event.searchTerm)
        }
    }

    private fun loadData() {
        Timber.i("Loading planned events")
        viewModelScope.launch {
            _uiState.value = UiState.Loading
            try {
                val events = eventRepository.getUpcomingEvents()
                Timber.i("Loaded ${events.size} planned events")
                _uiState.value = UiState.Success(events, events)
            } catch (e: Exception) {
                Timber.e(e,"Failed to load planned events")
                _uiState.value = UiState.Error(e.message ?: "unknown error")
            }
        }
    }

    private fun updateFilter(searchTerm: String) {
        when (_uiState.value) {
            is UiState.Success -> {
                val allEvents = (_uiState.value as UiState.Success).events
                _uiState.value = UiState.Success(
                    allEvents,
                    allEvents.filter { it.hasMatch(searchTerm) }
                )
            }
            else -> return
        }
    }
}