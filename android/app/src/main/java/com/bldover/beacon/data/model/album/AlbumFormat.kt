package com.bldover.beacon.data.model.album

enum class AlbumFormat(val displayName: String) {
    SINGLE("Single"),
    EP("EP"),
    LP("LP"),
    DOUBLE_LP("2xLP");

    companion object {
        fun fromString(value: String): AlbumFormat {
            return entries.find { it.displayName == value } ?: LP
        }
    }
}
