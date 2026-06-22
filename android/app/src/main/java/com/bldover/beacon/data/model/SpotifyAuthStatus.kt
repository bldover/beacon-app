package com.bldover.beacon.data.model

import java.time.Instant

data class SpotifyAuthStatus(
    val authenticated: Boolean,
    val expireTs: Instant
)
