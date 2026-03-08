package com.bldover.beacon.data.model.record

import com.bldover.beacon.data.dto.RecordDto
import com.bldover.beacon.data.model.artist.Artist

data class Record(
    var id: String? = null,
    var name: String = "",
    var artist: Artist = Artist(),
    var year: Int = 0,
    var signed: Boolean = false
) {
    constructor(dto: RecordDto) : this(
        id = dto.id.ifBlank { null },
        name = dto.name,
        artist = Artist(dto.artist),
        year = dto.year,
        signed = dto.signed
    )

    fun isPopulated(): Boolean {
        return name.isNotBlank() && artist.isPopulated() && year in 1000..9999
    }

    fun hasMatch(searchTerm: String): Boolean {
        return name.contains(searchTerm, ignoreCase = true)
            || artist.name.contains(searchTerm, ignoreCase = true)
    }

    fun deepCopy(): Record {
        return Record(
            id = id,
            name = name,
            artist = artist.deepCopy(),
            year = year,
            signed = signed
        )
    }
}
