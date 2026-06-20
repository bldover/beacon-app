package com.bldover.beacon.data.api

import com.bldover.beacon.BuildConfig
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.util.concurrent.TimeUnit

object RetrofitInstance {

    private val BASE_URL = BuildConfig.API_URL

    private val client = OkHttpClient.Builder()
        .addInterceptor(AuthInterceptor(BuildConfig.API_KEY))
        .addInterceptor(RetryInterceptor())
        .readTimeout(2, TimeUnit.MINUTES)
        .build()

    val retrofit: Retrofit = Retrofit.Builder()
        .baseUrl(BASE_URL)
        .client(client)
        .addConverterFactory(GsonConverterFactory.create())
        .build()
}