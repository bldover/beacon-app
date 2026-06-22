package com.bldover.beacon.data.spotify

import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.SharedFlow
import kotlinx.coroutines.flow.asSharedFlow

sealed class SpotifyAuthResult {
    data object Success : SpotifyAuthResult()
    data class Failure(val reason: String?) : SpotifyAuthResult()
}

object SpotifyAuthResultBus {
    private val _results = MutableSharedFlow<SpotifyAuthResult>(extraBufferCapacity = 1)
    val results: SharedFlow<SpotifyAuthResult> = _results.asSharedFlow()

    fun post(result: SpotifyAuthResult) {
        _results.tryEmit(result)
    }
}
