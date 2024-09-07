package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.EventApi
import com.bldover.beacon.data.model.Event
import com.bldover.beacon.data.model.EventDetail
import com.bldover.beacon.data.model.RawEvent
import com.bldover.beacon.data.model.RecommendationThreshold
import java.time.LocalDate

interface EventRepository {
    suspend fun getPastSavedEvents(): List<Event>
    suspend fun getFutureSavedEvents(): List<Event>
    suspend fun getEvent(eventId: String): Event
    suspend fun saveEvent(event: Event)
    suspend fun updateEvent(event: Event)
    suspend fun deleteEvent(event: Event)
    suspend fun getUpcomingEvents(): List<EventDetail>
    suspend fun getRecommendedEvents(threshold: RecommendationThreshold): List<EventDetail>
}

class EventRepositoryImpl(private val eventApi: EventApi): EventRepository {

    override suspend fun getPastSavedEvents(): List<Event> {
        return eventApi.getSavedEvents()
            .map { Event(it) }
            .filter { LocalDate.now().isAfter(it.date) }
            .toList()
    }

    override suspend fun getEvent(eventId: String): Event {
        return Event(eventApi.getEvent(eventId).first())
    }

    override suspend fun saveEvent(event: Event) {
        eventApi.addEvent(RawEvent(event))
    }

    override suspend fun updateEvent(event: Event) {
        deleteEvent(event)
        saveEvent(event)
    }

    override suspend fun deleteEvent(event: Event) {
        eventApi.deleteEvent(event.id!!)
    }

    override suspend fun getFutureSavedEvents(): List<Event> {
        return eventApi.getSavedEvents()
            .map { Event(it) }
            .filter { LocalDate.now().isBefore(it.date) || LocalDate.now().isEqual(it.date) }
            .toList()
    }

    override suspend fun getUpcomingEvents(): List<EventDetail> {
        return eventApi.getUpcomingEvents()
            .map { EventDetail(it) }
            .toList()
    }

    override suspend fun getRecommendedEvents(threshold: RecommendationThreshold): List<EventDetail> {
        return eventApi.getRecommendations(threshold)
            .map { EventDetail(it) }
            .filter { LocalDate.now().isBefore(it.date) || LocalDate.now().isEqual(it.date) }
            .toList()
    }
}