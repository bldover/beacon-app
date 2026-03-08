package com.bldover.beacon.ui.screens.editor.record

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bldover.beacon.data.model.record.Record
import com.bldover.beacon.data.repository.RecordRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

sealed class RecordState {
    data object Loading : RecordState()
    data class Success(
        val records: List<Record>,
        val filtered: List<Record>
    ) : RecordState()
    data class Error(val message: String) : RecordState()
}

@HiltViewModel
class RecordsViewModel @Inject constructor(
    private val recordRepository: RecordRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow<RecordState>(RecordState.Loading)
    val uiState: StateFlow<RecordState> = _uiState.asStateFlow()

    init {
        loadRecords()
    }

    fun loadRecords() {
        Timber.i("Loading records")
        viewModelScope.launch {
            _uiState.value = RecordState.Loading
            try {
                val records = recordRepository.getRecords()
                    .sortedWith(compareBy({ it.artist.name }, { it.year }))
                _uiState.value = RecordState.Success(records, records)
                Timber.i("Loaded ${records.size} records")
            } catch (e: Exception) {
                Timber.e(e, "Failed to load records")
                _uiState.value = RecordState.Error(e.message ?: "unknown error")
            }
        }
    }

    fun resetFilter() {
        if (_uiState.value !is RecordState.Success) return
        val state = _uiState.value as RecordState.Success
        _uiState.value = RecordState.Success(records = state.records, filtered = state.records)
    }

    fun applyFilter(searchTerm: String) {
        if (_uiState.value !is RecordState.Success) return
        val allRecords = (_uiState.value as RecordState.Success).records
        _uiState.value = RecordState.Success(
            allRecords,
            allRecords.filter { it.hasMatch(searchTerm) }
        )
    }

    fun addRecord(
        record: Record,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            if (!record.isPopulated()) {
                onError("Record is missing required fields")
                return@launch
            }
            try {
                recordRepository.addRecord(record)
                onSuccess()
                loadRecords()
            } catch (e: Exception) {
                Timber.e(e, "Failed to add record $record")
                onError("Error saving record ${record.name}, try again later")
            }
        }
    }

    fun updateRecord(
        record: Record,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            if (!record.isPopulated()) {
                onError("Record is missing required fields")
                return@launch
            }
            try {
                recordRepository.updateRecord(record)
                onSuccess()
                loadRecords()
            } catch (e: Exception) {
                Timber.e(e, "Failed to update record $record")
                onError("Error saving record ${record.name}, try again later")
            }
        }
    }

    fun deleteRecord(
        record: Record,
        onSuccess: () -> Unit = {},
        onError: (String) -> Unit = {}
    ) {
        viewModelScope.launch {
            try {
                recordRepository.deleteRecord(record)
                onSuccess()
                loadRecords()
            } catch (e: Exception) {
                Timber.e(e, "Failed to delete record $record")
                onError("Error deleting record ${record.name}, try again later")
            }
        }
    }
}
