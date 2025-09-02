package com.bldover.beacon.data.api

import com.bldover.beacon.data.dto.GenreResponseDto
import retrofit2.http.GET

interface GenreApi {

    @GET("v1/genres")
    suspend fun getGenres(): GenreResponseDto
}