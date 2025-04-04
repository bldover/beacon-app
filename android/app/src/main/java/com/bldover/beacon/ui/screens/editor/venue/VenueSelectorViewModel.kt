package com.bldover.beacon.ui.screens.editor.venue

import androidx.lifecycle.ViewModel
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.venue.Venue
import dagger.hilt.android.lifecycle.HiltViewModel
import timber.log.Timber
import javax.inject.Inject

@HiltViewModel
class VenueSelectorViewModel @Inject constructor() : ViewModel() {

    private var onSelect: (Venue) -> Unit = {}

    fun launchSelector(
        navController: NavController,
        onSelect: (Venue) -> Unit
    ) {
        Timber.d("launching venue selector")
        this.onSelect = onSelect
        navController.navigate(Screen.SELECT_VENUE.name)
    }

    fun selectVenue(venue: Venue) {
        Timber.d("venue selector - selected venue $venue")
        onSelect(venue)
    }
}