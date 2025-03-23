package com.bldover.beacon.data.api

import com.bldover.beacon.Config
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.util.concurrent.TimeUnit

object RetrofitInstance {

    private const val BASE_URL = Config.API_URL

    private val client = OkHttpClient.Builder()
        .addInterceptor(RetryInterceptor())
        .readTimeout(2, TimeUnit.MINUTES)
        .build()

    val retrofit: Retrofit = Retrofit.Builder()
        .baseUrl(BASE_URL)
        .client(client)
        .addConverterFactory(GsonConverterFactory.create())
        .build()
}