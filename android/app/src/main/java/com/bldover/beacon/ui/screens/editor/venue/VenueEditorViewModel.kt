package com.bldover.beacon.ui.screens.editor.venue

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.venue.Venue
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import javax.inject.Inject

@HiltViewModel
class VenueEditorViewModel @Inject constructor() : ViewModel() {

    private val _venueState = MutableStateFlow(Venue())
    val venueState = _venueState.asStateFlow()

    private var onSave: (Venue) -> Unit = {}

    fun launchEditor(
        navController: NavController,
        venue: Venue? = null,
        onSave: (Venue) -> Unit,
    ) {
        this.onSave = onSave
        _venueState.value = venue?.deepCopy() ?: Venue()
        navController.navigate(Screen.EDIT_VENUE.name)
    }

    fun updateName(name: String) {
        _venueState.value = _venueState.value.copy(name = name)
    }

    fun updateCity(city: String) {
        _venueState.value = _venueState.value.copy(city = city)
    }

    fun updateState(state: String) {
        _venueState.value = _venueState.value.copy(state = state)
    }

    fun onSave() {
        onSave(_venueState.value)
    }
}