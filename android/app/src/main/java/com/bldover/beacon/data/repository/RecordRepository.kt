package com.bldover.beacon.data.repository

import com.bldover.beacon.data.api.RecordApi
import com.bldover.beacon.data.dto.RecordDto
import com.bldover.beacon.data.model.record.Record

interface RecordRepository {
    suspend fun getRecords(): List<Record>
    suspend fun addRecord(record: Record): Record
    suspend fun updateRecord(record: Record): Record
    suspend fun deleteRecord(record: Record)
}

class RecordRepositoryImpl(private val recordApi: RecordApi) : RecordRepository {

    override suspend fun getRecords(): List<Record> {
        return recordApi.getRecords().map { Record(it) }
    }

    override suspend fun addRecord(record: Record): Record {
        val newRecord = recordApi.addRecord(RecordDto(record))
        return Record(newRecord)
    }

    override suspend fun updateRecord(record: Record): Record {
        recordApi.updateRecord(record.id!!, RecordDto(record))
        return record
    }

    override suspend fun deleteRecord(record: Record) {
        recordApi.deleteRecord(record.id!!)
    }
}
