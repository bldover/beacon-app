package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.SpotifyApi
import com.bldover.beacon.data.model.SpotifyAuthStatus

interface SpotifyRepository {
    suspend fun startReauth(): String
    suspend fun getAuthStatus(): SpotifyAuthStatus
}

class SpotifyRepositoryImpl(private val spotifyApi: SpotifyApi) : SpotifyRepository {

    override suspend fun startReauth(): String {
        return spotifyApi.startReauth().authUrl
    }

    override suspend fun getAuthStatus(): SpotifyAuthStatus {
        return spotifyApi.getAuthStatus().toModel()
    }
}
