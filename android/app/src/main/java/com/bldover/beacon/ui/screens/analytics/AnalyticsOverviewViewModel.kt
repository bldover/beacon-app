package com.bldover.beacon.ui.screens.analytics

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.analytics.AnalyticsSummary
import com.bldover.beacon.data.repository.AnalyticsRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class AnalyticsOverviewState {
    data object Loading : AnalyticsOverviewState()
    data class Success(val summary: AnalyticsSummary) : AnalyticsOverviewState()
    data class Error(val message: String) : AnalyticsOverviewState()
}

@HiltViewModel
class AnalyticsOverviewViewModel @Inject constructor(
    private val analyticsRepository: AnalyticsRepository
) : ViewModel() {

    private val _state = MutableStateFlow<AnalyticsOverviewState>(AnalyticsOverviewState.Loading)
    val state: StateFlow<AnalyticsOverviewState> = _state.asStateFlow()

    init {
        loadData()
    }

    fun loadData() {
        Timber.i("Loading analytics summary")
        viewModelScope.launch {
            _state.value = AnalyticsOverviewState.Loading
            try {
                val summary = analyticsRepository.getSummary()
                _state.value = AnalyticsOverviewState.Success(summary)
            } catch (e: Exception) {
                Timber.e(e, "Failed to load analytics summary")
                _state.value = AnalyticsOverviewState.Error("Failed to load analytics")
            }
        }
    }
}
