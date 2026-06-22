package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.SpotifyAuthStatus
import com.google.gson.annotations.SerializedName
import java.time.Instant

data class SpotifyReauthDto(
    @SerializedName("authUrl") val authUrl: String
)

data class SpotifyAuthStatusDto(
    @SerializedName("authenticated") val authenticated: Boolean,
    @SerializedName("expireTs") val expireTs: String
) {
    fun toModel(): SpotifyAuthStatus {
        return SpotifyAuthStatus(
            authenticated = authenticated,
            expireTs = Instant.parse(expireTs)
        )
    }
}
