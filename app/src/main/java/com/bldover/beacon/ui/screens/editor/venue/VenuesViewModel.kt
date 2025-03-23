package com.bldover.beacon.ui.screens.editor.venue

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.Venue
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