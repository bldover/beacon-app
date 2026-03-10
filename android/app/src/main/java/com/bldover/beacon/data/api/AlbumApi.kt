package com.bldover.beacon.data.api

import com.bldover.beacon.data.dto.AlbumDto
import okhttp3.MultipartBody
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.Multipart
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Part
import retrofit2.http.Path

interface AlbumApi {

    @GET("v1/albums")
    suspend fun getAlbums(): List<AlbumDto>

    @POST("v1/albums")
    suspend fun addAlbum(@Body album: AlbumDto): AlbumDto

    @PUT("v1/albums/{id}")
    suspend fun updateAlbum(@Path("id") id: String, @Body album: AlbumDto)

    @DELETE("v1/albums/{id}")
    suspend fun deleteAlbum(@Path("id") id: String)

    @Multipart
    @POST("v1/albums/images")
    suspend fun uploadAlbumImage(@Part image: MultipartBody.Part): Map<String, String>
}
