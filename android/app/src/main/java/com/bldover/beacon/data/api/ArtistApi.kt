package com.bldover.beacon.data.api

import com.bldover.beacon.data.model.artist.RawArtist
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path

interface ArtistApi {

    @GET("v1/artists")
    suspend fun getArtists(): List<RawArtist>

    @POST("v1/artists")
    suspend fun addArtist(@Body artist: RawArtist): RawArtist

    @PUT("v1/artists/{id}")
    suspend fun updateArtist(@Path("id") id: String, @Body artist: RawArtist)

    @DELETE("v1/artists/{id}")
    suspend fun deleteArtist(@Path("id") id: String)
}