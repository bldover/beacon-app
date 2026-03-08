package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.AlbumApi
import com.bldover.beacon.data.dto.AlbumDto
import com.bldover.beacon.data.model.album.Album

interface AlbumRepository {
    suspend fun getAlbums(): List<Album>
    suspend fun addAlbum(album: Album): Album
    suspend fun updateAlbum(album: Album): Album
    suspend fun deleteAlbum(album: Album)
}

class AlbumRepositoryImpl(private val albumApi: AlbumApi) : AlbumRepository {

    override suspend fun getAlbums(): List<Album> {
        return albumApi.getAlbums().map { Album(it) }
    }

    override suspend fun addAlbum(album: Album): Album {
        val newAlbum = albumApi.addAlbum(AlbumDto(album))
        return Album(newAlbum)
    }

    override suspend fun updateAlbum(album: Album): Album {
        albumApi.updateAlbum(album.id!!, AlbumDto(album))
        return album
    }

    override suspend fun deleteAlbum(album: Album) {
        albumApi.deleteAlbum(album.id!!)
    }
}
