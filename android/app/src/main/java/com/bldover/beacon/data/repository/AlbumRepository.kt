package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.AlbumApi
import com.bldover.beacon.data.dto.AlbumDto
import com.bldover.beacon.data.model.album.Album
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.toRequestBody

interface AlbumRepository {
    suspend fun getAlbums(): List<Album>
    suspend fun addAlbum(album: Album): Album
    suspend fun updateAlbum(album: Album): Album
    suspend fun deleteAlbum(album: Album)
    suspend fun uploadCoverImage(bytes: ByteArray, contentType: String): String
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

    override suspend fun uploadCoverImage(bytes: ByteArray, contentType: String): String {
        val requestBody = bytes.toRequestBody(contentType.toMediaType())
        val part = MultipartBody.Part.createFormData("image", "cover.jpg", requestBody)
        val response = albumApi.uploadAlbumImage(part)
        return response["url"] ?: error("Missing url in image upload response")
    }
}
