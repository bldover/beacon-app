package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.ArtistApi
import com.bldover.beacon.data.model.Artist

interface ArtistRepository {
    suspend fun getArtists(): List<Artist>
}

class ArtistRepositoryImpl(private val artistApi: ArtistApi) : ArtistRepository {

    override suspend fun getArtists(): List<Artist> {
        return artistApi.getArtists()
    }
}