package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.ArtistApi
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.data.dto.ArtistDto
import timber.log.Timber

interface ArtistRepository {
    suspend fun getArtists(): List<Artist>
    suspend fun addArtist(artist: Artist): Artist
    suspend fun updateArtist(artist: Artist): Artist
    suspend fun deleteArtist(artist: Artist)
}

class ArtistRepositoryImpl(private val artistApi: ArtistApi) : ArtistRepository {

    override suspend fun getArtists(): List<Artist> {
        Timber.i("Call to getArtists()")
        return artistApi.getArtists().map { Artist(it) }
    }

    override suspend fun addArtist(artist: Artist): Artist {
        val newArtist = artistApi.addArtist(ArtistDto(artist))
        return Artist(newArtist)
    }

    override suspend fun updateArtist(artist: Artist): Artist {
        artistApi.updateArtist(artist.id.primary!!, ArtistDto(artist))
        return artist
    }

    override suspend fun deleteArtist(artist: Artist) {
        artistApi.deleteArtist(artist.id.primary!!)
    }
}