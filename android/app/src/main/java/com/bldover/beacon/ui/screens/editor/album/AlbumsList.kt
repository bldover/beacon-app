package com.bldover.beacon.ui.screens.editor.album

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
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bldover.beacon.data.model.album.Album
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.BasicSearchBar
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScrollableItemList

@Composable
fun SearchableAlbumsList(
    albumState: AlbumState,
    onSearchAlbums: (String) -> Unit,
    onAlbumSelected: (Album) -> Unit,
    onNewAlbum: () -> Unit
) {
    Scaffold(
        topBar = {
            BasicSearchBar(
                modifier = Modifier.fillMaxWidth(),
                enabled = albumState is AlbumState.Success,
                onQueryChange = onSearchAlbums
            )
        }
    ) { innerPadding ->
        Column(modifier = Modifier.padding(innerPadding)) {
            Spacer(modifier = Modifier.height(16.dp))
            when (albumState) {
                is AlbumState.Success -> AlbumList(
                    albumState.filtered,
                    onAlbumSelected,
                    onNewAlbum
                )
                is AlbumState.Error -> LoadErrorMessage()
                is AlbumState.Loading -> LoadingSpinner()
            }
        }
    }
}

@Composable
private fun AlbumList(
    albums: List<Album>,
    onAlbumSelected: (Album) -> Unit,
    onNewAlbum: () -> Unit
) {
    ScrollableItemList(
        items = albums,
        topAnchor = { NewAlbumCard(onNewAlbum) },
        getItemKey = { it.id!! }
    ) { album ->
        AlbumDetailsCard(
            album = album,
            onClick = { onAlbumSelected(album) }
        )
    }
}

@Composable
private fun AlbumDetailsCard(
    album: Album,
    onClick: () -> Unit
) {
    BasicCard(modifier = Modifier.clickable { onClick() }) {
        Text(
            text = "${album.artists.joinToString(", ") { it.name }} - ${album.name}",
            style = MaterialTheme.typography.bodyLarge
        )
        Text(
            text = album.year.toString(),
            style = MaterialTheme.typography.bodySmall
        )
    }
}

@Composable
private fun NewAlbumCard(onNewAlbum: () -> Unit) {
    BasicCard(modifier = Modifier.clickable { onNewAlbum() }) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            Text(text = "New Album")
            Icon(
                imageVector = Icons.Default.AddCircle,
                contentDescription = "New Album"
            )
        }
    }
}
