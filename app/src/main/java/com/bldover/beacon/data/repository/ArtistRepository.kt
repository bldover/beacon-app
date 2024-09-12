package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.ArtistApi
import com.bldover.beacon.data.model.Artist
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
        return artistApi.getArtists().onEach { it.genreSet = true }
    }

    override suspend fun addArtist(artist: Artist): Artist {
        return artistApi.addArtist(artist)
    }

    override suspend fun updateArtist(artist: Artist): Artist {
        artistApi.updateArtist(artist.id!!, artist)
        return artist
    }

    override suspend fun deleteArtist(artist: Artist) {
        artistApi.deleteArtist(artist.id!!)
    }
}