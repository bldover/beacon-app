package com.bldover.beacon.data.api

import com.bldover.beacon.data.dto.RecordDto
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path

interface RecordApi {

    @GET("v1/records")
    suspend fun getRecords(): List<RecordDto>

    @POST("v1/records")
    suspend fun addRecord(@Body record: RecordDto): RecordDto

    @PUT("v1/records/{id}")
    suspend fun updateRecord(@Path("id") id: String, @Body record: RecordDto)

    @DELETE("v1/records/{id}")
    suspend fun deleteRecord(@Path("id") id: String)
}
