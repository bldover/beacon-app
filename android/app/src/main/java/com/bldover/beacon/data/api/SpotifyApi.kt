package com.bldover.beacon.data.api

import com.bldover.beacon.data.dto.SpotifyAuthStatusDto
import com.bldover.beacon.data.dto.SpotifyReauthDto
import retrofit2.http.GET

interface SpotifyApi {

    @GET("v1/spotify/auth/start")
    suspend fun startReauth(): SpotifyReauthDto

    @GET("v1/spotify/auth/status")
    suspend fun getAuthStatus(): SpotifyAuthStatusDto
}
