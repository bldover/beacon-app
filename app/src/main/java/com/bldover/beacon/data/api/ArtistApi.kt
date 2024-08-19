package com.bldover.beacon.data.api

import com.bldover.beacon.data.model.Artist
import retrofit2.http.GET

interface ArtistApi {

    @GET("v1/artists")
    suspend fun getArtists(): List<Artist>
}