package com.bldover.beacon.data.model.artist

import com.bldover.beacon.data.dto.GenresDto

data class Genres(
    var spotify: List<String> = emptyList(),
    var lastFm: List<String> = emptyList(),
    var ticketmaster: List<String> = emptyList(), // TM genres are terrible, don't use these
    var user: List<String> = emptyList()
) : Comparable<Genres> {
    constructor(genresDto: GenresDto) : this(
        spotify = genresDto.spotify ?: emptyList(),
        lastFm = genresDto.lastFm ?: emptyList(),
        ticketmaster = genresDto.ticketmaster ?: emptyList(),
        user = genresDto.user ?: emptyList()
    )

    fun hasGenre(genre: String) : Boolean = getGenres().any { it.contains(genre, ignoreCase = true) }
    fun getGenres() : List<String> = user.ifEmpty { spotify }.ifEmpty { lastFm }
    fun getTopGenre() : String? = getGenres().firstOrNull()
    fun hasUserOverride() : Boolean = user.isNotEmpty()

    override fun compareTo(other: Genres): Int {
        val genres = getGenres()
        val otherGenres = other.getGenres()
        for (i in 0 until genres.size.coerceAtMost(otherGenres.size)) {
            if (genres[i] != otherGenres[i]) {
                return genres[i].compareTo(otherGenres[i])
            }
        }
        return if (genres.size < otherGenres.size) 1 else -1
    }

    fun deepCopy() : Genres {
        return Genres(
            spotify = spotify.toMutableList(),
            lastFm = lastFm.toMutableList(),
            ticketmaster = ticketmaster.toMutableList(),
            user = user.toMutableList()
        )
    }
}