package com.bldover.beacon.ui.screens.editor.genre

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AddCircle
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScrollableItemList

@Composable
fun SearchableGenresList(
    genreState: GenreState,
    onSearchGenres: (String) -> Unit,
    onGenreSelected: (String) -> Unit,
    onCustomGenre: (String) -> Unit
) {
    Scaffold(
        topBar = {
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = genreState is GenreState.Success,
                onQueryChange = onSearchGenres
            )
        }
    ) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding)) {
            Spacer(modifier = Modifier.height(16.dp))
            when (genreState) {
                is GenreState.Success -> GenreList(
                    genreState.filtered,
                    onGenreSelected,
                    onCustomGenre
                )
                is GenreState.Error -> LoadErrorMessage()
                is GenreState.Loading -> LoadingSpinner()
            }
        }
    }
}

@Composable
private fun GenreList(
    genres: List<String>,
    onGenreSelected: (String) -> Unit,
    onCustomGenre: (String) -> Unit
) {
    ScrollableItemList(
        items = genres,
        topAnchor = { CustomGenreCard(onCustomGenre) },
        getItemKey = { it }
    ) { genre ->
        GenreCard(
            genre = genre,
            onClick = {
                onGenreSelected(genre)
            }
        )
    }
}

@Composable
private fun GenreCard(
    genre: String,
    onClick: () -> Unit
) {
    BasicCard(
        modifier = Modifier.clickable { onClick() }
    ) {
        Text(
            text = genre,
            style = MaterialTheme.typography.bodyMedium
        )
    }
}

@Composable
private fun CustomGenreCard(
    onCustomGenre: (String) -> Unit
) {
    var customGenreText by remember { mutableStateOf("") }
    
    BasicCard {
        Column {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
            ) {
                Text(text = "Custom Genre")
                Icon(
                    imageVector = Icons.Default.AddCircle,
                    contentDescription = "Custom Genre"
                )
            }
            Spacer(modifier = Modifier.height(8.dp))
            TextField(
                value = customGenreText,
                onValueChange = { customGenreText = it },
                label = { Text("Enter custom genre") },
                modifier = Modifier.fillMaxWidth(),
                singleLine = true
            )
            Spacer(modifier = Modifier.height(8.dp))
            BasicCard(
                modifier = Modifier.clickable {
                    if (customGenreText.isNotBlank()) {
                        onCustomGenre(customGenreText.trim())
                    }
                }
            ) {
                Text("Add Custom Genre")
            }
        }
    }
}