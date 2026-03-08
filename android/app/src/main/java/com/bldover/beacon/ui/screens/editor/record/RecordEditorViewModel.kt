package com.bldover.beacon.ui.screens.editor.record

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.model.record.Record
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import java.time.LocalDate
import javax.inject.Inject

@HiltViewModel
class RecordEditorViewModel @Inject constructor() : ViewModel() {

    private val _recordState = MutableStateFlow(Record())
    val recordState = _recordState.asStateFlow()

    private var onSave: (Record) -> Unit = {}
    private var onDelete: (Record) -> Unit = {}
    var showDelete: Boolean = false
        private set

    fun launchEditor(
        navController: NavController,
        record: Record? = null,
        onSave: (Record) -> Unit,
        onDelete: ((Record) -> Unit)? = null
    ) {
        this.onSave = onSave
        this.onDelete = onDelete ?: {}
        this.showDelete = record != null
        _recordState.value = record?.deepCopy() ?: Record(year = LocalDate.now().year)
        navController.navigate(Screen.EDIT_RECORD.name)
    }

    fun updateName(name: String) {
        _recordState.value = _recordState.value.copy(name = name)
    }

    fun updateArtist(artist: Artist) {
        _recordState.value = _recordState.value.copy(artist = artist)
    }

    fun updateYear(year: Int) {
        _recordState.value = _recordState.value.copy(year = year)
    }

    fun updateSigned(signed: Boolean) {
        _recordState.value = _recordState.value.copy(signed = signed)
    }

    fun onSave() {
        onSave(_recordState.value)
    }

    fun onDelete() {
        onDelete(_recordState.value)
    }
}
