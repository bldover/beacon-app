package com.bldover.beacon.ui.screens.editor.venue

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.Direction
import com.bldover.beacon.data.model.OrderField
import com.bldover.beacon.data.model.venue.Venue
import com.bldover.beacon.data.model.venue.VenueOrdering
import com.bldover.beacon.data.repository.VenueRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class VenueState {
    data object Loading : VenueState()
    data class Success(
        val venues: List<Venue>,
        val filtered: List<Venue>
    ) : VenueState()
    data class Error(val message: String) : VenueState()
}

@HiltViewModel
class VenuesViewModel @Inject constructor(
    private val venueRepository: VenueRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow<VenueState>(VenueState.Loading)
    val uiState: StateFlow<VenueState> = _uiState.asStateFlow()
    private val _ordering = MutableStateFlow(VenueOrdering(OrderField.NAME, Direction.ASCENDING))
    val ordering: StateFlow<VenueOrdering> = _ordering.asStateFlow()

    init {
        loadVenues()
    }

    fun loadVenues() {
        Timber.i("Loading venues")
        viewModelScope.launch {
            _uiState.value = VenueState.Loading
            try {
                val venues = venueRepository.getVenues().sortedBy(Venue::name)
                Timber.i("Loaded ${venues.size} venues")
                _uiState.value = VenueState.Success(venues, venues)
                applyOrdering(_ordering.value)
            } catch (e: Exception) {
                Timber.e(e,"Failed to load venues")
                _uiState.value = VenueState.Error(e.message ?: "unknown error")
            }
        }
    }

    fun resetFilter() {
        if (_uiState.value !is VenueState.Success) return
        _uiState.value = VenueState.Success(
            venues = (_uiState.value as VenueState.Success).venues,
            filtered = (_uiState.value as VenueState.Success).venues
        )
    }

    fun applyFilter(searchTerm: String) {
        when (_uiState.value) {
            is VenueState.Success -> {
                val allVenues = (_uiState.value as VenueState.Success).venues
                _uiState.value = VenueState.Success(
                    allVenues,
                    allVenues.filter { it.hasMatch(searchTerm) }
                )
            }
            else -> return
        }
    }

    fun applyOrdering(ordering: VenueOrdering) {
        Timber.i("Sorting venues events by $ordering")
        when (_uiState.value) {
            is VenueState.Success -> {
                val state = (_uiState.value as VenueState.Success)
                _uiState.value = VenueState.Success(
                    state.venues.sortedWith(Comparator(ordering::compare)),
                    state.filtered.sortedWith(Comparator(ordering::compare))
                )
                _ordering.value = ordering
            }
            else -> {
                Timber.w("Sorting venues by $ordering - not in success state")
                return
            }
        }
    }

    fun addVenue(
        venue: Venue,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                venueRepository.addVenue(venue)
            } catch (e: Exception) {
                Timber.e(e,"Failed to add venue $venue")
                onError("Error saving venue ${venue.name}, try again later")
                return@launch
            }
            onSuccess()
            loadVenues()
        }
    }

    fun updateVenue(
        venue: Venue,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                venueRepository.updateVenue(venue)
            } catch (e: Exception) {
                Timber.e(e,"Failed to update venue $venue")
                onError("Error saving venue ${venue.name}, try again later")
                return@launch
            }
            onSuccess()
            loadVenues()
        }
    }

    fun deleteVenue(
        venue: Venue,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                venueRepository.deleteVenue(venue)
            } catch (e: Exception) {
                Timber.e(e, "Failed to delete venue $venue")
                onError("Error deleting venue ${venue.name}, try again later")
                return@launch
            }
            onSuccess()
            loadVenues()
        }
    }
}