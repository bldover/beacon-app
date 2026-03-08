package com.bldover.beacon.data.dto

import com.bldover.beacon.data.model.record.Record

data class RecordDto(
    val id: String,
    val name: String,
    val artist: ArtistDto,
    val year: Int,
    val signed: Boolean
) {
    constructor(record: Record) : this(
        id = record.id ?: "",
        name = record.name,
        artist = ArtistDto(record.artist),
        year = record.year,
        signed = record.signed
    )
}
