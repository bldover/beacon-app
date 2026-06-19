package com.bldover.beacon.ui.screens.analytics

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.Direction
import com.bldover.beacon.data.model.OrderField
import com.bldover.beacon.data.model.analytics.AnalyticsCategory
import com.bldover.beacon.data.model.analytics.AnalyticsCount
import com.bldover.beacon.data.model.analytics.AnalyticsOrdering
import com.bldover.beacon.data.repository.AnalyticsRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class AnalyticsListState {
    data object Loading : AnalyticsListState()
    data class Success(
        val all: List<AnalyticsCount>,
        val sorted: List<AnalyticsCount>
    ) : AnalyticsListState()
    data class Error(val message: String) : AnalyticsListState()
}

@HiltViewModel
class AnalyticsListViewModel @Inject constructor(
    private val analyticsRepository: AnalyticsRepository
) : ViewModel() {

    private val states: Map<AnalyticsCategory, MutableStateFlow<AnalyticsListState>> =
        AnalyticsCategory.entries.associateWith { MutableStateFlow<AnalyticsListState>(AnalyticsListState.Loading) }
    private val orderings: Map<AnalyticsCategory, MutableStateFlow<AnalyticsOrdering>> =
        AnalyticsCategory.entries.associateWith {
            MutableStateFlow(AnalyticsOrdering(OrderField.COUNT, Direction.DESCENDING))
        }

    fun stateFor(category: AnalyticsCategory): StateFlow<AnalyticsListState> =
        states.getValue(category).asStateFlow()

    fun orderingFor(category: AnalyticsCategory): StateFlow<AnalyticsOrdering> =
        orderings.getValue(category).asStateFlow()

    fun load(category: AnalyticsCategory) {
        Timber.i("Loading analytics list for $category")
        val state = states.getValue(category)
        val ordering = orderings.getValue(category)
        viewModelScope.launch {
            state.value = AnalyticsListState.Loading
            try {
                val items = analyticsRepository.getCounts(category)
                val sorted = items.sortedWith(Comparator(ordering.value::compare))
                state.value = AnalyticsListState.Success(items, sorted)
            } catch (e: Exception) {
                Timber.e(e, "Failed to load analytics list for $category")
                state.value = AnalyticsListState.Error("Failed to load data")
            }
        }
    }

    fun sort(category: AnalyticsCategory, ordering: AnalyticsOrdering) {
        val stateFlow = states.getValue(category)
        val orderingFlow = orderings.getValue(category)
        val current = stateFlow.value
        orderingFlow.value = ordering
        if (current is AnalyticsListState.Success) {
            stateFlow.value = AnalyticsListState.Success(
                all = current.all,
                sorted = current.all.sortedWith(Comparator(ordering::compare))
            )
        }
    }
}
