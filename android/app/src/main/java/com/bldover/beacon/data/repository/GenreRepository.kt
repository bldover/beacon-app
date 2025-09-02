package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.GenreApi
import com.bldover.beacon.data.dto.GenreResponseDto

interface GenreRepository {
    suspend fun getGenres(): GenreResponseDto
}

class GenreRepositoryImpl(private val genreApi: GenreApi) : GenreRepository {

    override suspend fun getGenres(): GenreResponseDto {
        return genreApi.getGenres()
    }
}