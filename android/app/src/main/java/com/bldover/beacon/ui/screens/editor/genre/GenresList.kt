package com.bldover.beacon.ui.screens.editor.genre

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.artist.Artist
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.editor.GenreCard
import com.bldover.beacon.ui.components.editor.NewGenreDialogEditCard

@Composable
fun SearchableGenresList(
    artist: Artist?,
    allUserGenres: List<String>,
    filteredUserGenres: List<String>,
    isFiltering: Boolean,
    onSearchGenres: (String) -> Unit,
    onGenreSelected: (String) -> Unit
) {
    Scaffold(
        topBar = {
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = true,
                onQueryChange = onSearchGenres
            )
        }
    ) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding)) {
            Spacer(modifier = Modifier.height(16.dp))
            GenreList(
                artist = artist,
                allUserGenres = allUserGenres,
                filteredUserGenres = filteredUserGenres,
                isFiltering = isFiltering,
                onGenreSelected = onGenreSelected
            )
        }
    }
}

@Composable
private fun GenreList(
    artist: Artist?,
    allUserGenres: List<String>,
    filteredUserGenres: List<String>,
    isFiltering: Boolean,
    onGenreSelected: (String) -> Unit
) {
    val allItems = mutableListOf<GenreItem>()
    
    if (isFiltering) {
        if (filteredUserGenres.isNotEmpty()) {
            allItems.add(GenreItem.SectionHeader("All"))
            filteredUserGenres.forEach { genre ->
                allItems.add(GenreItem.Genre(genre, "All"))
            }
        }
    } else {
        val spotifyGenres = artist?.genres?.spotify ?: emptyList()
        val lastFmGenres = artist?.genres?.lastFm ?: emptyList()
        
        if (spotifyGenres.isNotEmpty()) {
            allItems.add(GenreItem.SectionHeader("Spotify"))
            spotifyGenres.forEach { genre ->
                val hasAccentBorder = allUserGenres.contains(genre)
                allItems.add(GenreItem.Genre(genre, "Spotify", hasAccentBorder))
            }
        }
        
        if (lastFmGenres.isNotEmpty()) {
            allItems.add(GenreItem.SectionHeader("Last.fm"))
            lastFmGenres.forEach { genre ->
                val hasAccentBorder = allUserGenres.contains(genre)
                allItems.add(GenreItem.Genre(genre, "Last.fm", hasAccentBorder))
            }
        }
        
        if (allUserGenres.isNotEmpty()) {
            allItems.add(GenreItem.SectionHeader("All"))
            allUserGenres.forEach { genre ->
                allItems.add(GenreItem.Genre(genre, "All"))
            }
        }
    }
    
    ScrollableItemList(
        items = allItems,
        topAnchor = { NewGenreDialogEditCard(onNewGenre = onGenreSelected) }
    ) { item ->
        when (item) {
            is GenreItem.SectionHeader -> {
                Column {
                    HorizontalDivider()
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(
                        text = item.title,
                        style = MaterialTheme.typography.titleMedium.copy(fontWeight = FontWeight.Bold),
                        modifier = Modifier.padding(horizontal = 16.dp)
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                }
            }
            is GenreItem.Genre -> {
                GenreCard(
                    genre = item.name,
                    onClick = { onGenreSelected(item.name) },
                    hasAccentBorder = item.hasAccentBorder
                )
            }
        }
    }
}

private sealed class GenreItem {
    data class SectionHeader(val title: String) : GenreItem()
    data class Genre(val name: String, val section: String, val hasAccentBorder: Boolean = false) : GenreItem()
}
