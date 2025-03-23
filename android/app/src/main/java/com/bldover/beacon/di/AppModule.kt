package com.bldover.beacon.di

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.PreferenceDataStoreFactory
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.preferencesDataStoreFile
import com.bldover.beacon.data.api.ArtistApi
import com.bldover.beacon.data.api.EventApi
import com.bldover.beacon.data.api.RetrofitInstance
import com.bldover.beacon.data.api.VenueApi
import com.bldover.beacon.data.repository.ArtistRepository
import com.bldover.beacon.data.repository.ArtistRepositoryImpl
import com.bldover.beacon.data.repository.EventRepository
import com.bldover.beacon.data.repository.EventRepositoryImpl
import com.bldover.beacon.data.repository.UserSettingsRepository
import com.bldover.beacon.data.repository.VenueRepository
import com.bldover.beacon.data.repository.VenueRepositoryImpl
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import retrofit2.create
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object AppModule {

    @Provides
    @Singleton
    fun providesEventApi(): EventApi {
        return RetrofitInstance.retrofit.create<EventApi>()
    }

    @Provides
    @Singleton
    fun providesEventRepository(eventApi: EventApi): EventRepository {
        return EventRepositoryImpl(eventApi)
    }

    @Provides
    @Singleton
    fun providesVenueApi(): VenueApi {
        return RetrofitInstance.retrofit.create<VenueApi>()
    }

    @Provides
    @Singleton
    fun providesVenueRepository(venueApi: VenueApi): VenueRepository {
        return VenueRepositoryImpl(venueApi)
    }

    @Provides
    @Singleton
    fun providesArtistApi(): ArtistApi {
        return RetrofitInstance.retrofit.create<ArtistApi>()
    }

    @Provides
    @Singleton
    fun providesArtistRepository(artistApi: ArtistApi): ArtistRepository {
        return ArtistRepositoryImpl(artistApi)
    }

    @Provides
    @Singleton
    fun providesDataStore(@ApplicationContext context: Context): DataStore<Preferences> {
        return PreferenceDataStoreFactory.create(
            produceFile = { context.preferencesDataStoreFile("userSettings") }
        )
    }

    @Provides
    @Singleton
    fun providesUserSettingsRepository(dataStore: DataStore<Preferences>): UserSettingsRepository {
        return UserSettingsRepository(dataStore)
    }
}