package com.bldover.beacon.data.api

import com.bldover.beacon.data.model.Artist
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Query

interface ArtistApi {

    @GET("v1/artists")
    suspend fun getArtists(): List<Artist>

    @POST("v1/artists")
    suspend fun addArtist(@Body artist: Artist): Artist

    @DELETE("v1/artists")
    suspend fun deleteArtist(@Query("id") id: String)
}