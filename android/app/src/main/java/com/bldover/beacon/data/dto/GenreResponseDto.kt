package com.bldover.beacon.data.dto

data class GenreResponseDto(
    val user: List<String>,
    val spotify: List<String>,
    val lastFm: List<String>,
    val ticketmaster: List<String>
)