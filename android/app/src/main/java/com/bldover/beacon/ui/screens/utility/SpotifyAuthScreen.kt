package com.bldover.beacon.ui.screens.utility

import android.net.Uri
import androidx.browser.customtabs.CustomTabsIntent
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.model.SnackbarState
import com.bldover.beacon.data.model.SpotifyAuthStatus
import com.bldover.beacon.data.spotify.SpotifyAuthResult
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import timber.log.Timber
import java.time.ZoneId
import java.time.format.DateTimeFormatter

private val expireTsFormatter: DateTimeFormatter =
    DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm").withZone(ZoneId.systemDefault())

@Composable
fun SpotifyAuthScreen(
    navController: NavController,
    snackbarState: SnackbarState,
    spotifyAuthViewModel: SpotifyAuthViewModel
) {
    Timber.d("composing SpotifyAuthScreen")
    val context = LocalContext.current
    val authStatusState by spotifyAuthViewModel.authStatus.collectAsState()

    LaunchedEffect(Unit) {
        spotifyAuthViewModel.loadAuthStatus()
    }
    LaunchedEffect(Unit) {
        spotifyAuthViewModel.authResults.collect { result ->
            when (result) {
                is SpotifyAuthResult.Success -> snackbarState.showSnackbar("Spotify reconnected")
                is SpotifyAuthResult.Failure -> snackbarState.showSnackbar(
                    "Spotify reconnect failed${result.reason?.let { ": $it" } ?: ""}"
                )
            }
        }
    }

    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Spotify Auth",
                leadingIcon = { BackButton(navController) }
            )
        }
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Top,
            modifier = Modifier.fillMaxSize()
        ) {
            when (val state = authStatusState) {
                is SpotifyAuthStatusState.Loading -> LoadingSpinner()
                is SpotifyAuthStatusState.Error -> Text("Failed to load auth status: ${state.message}")
                is SpotifyAuthStatusState.Loaded -> StatusCard(state.status)
            }
            Spacer(modifier = Modifier.height(24.dp))
            Button(
                onClick = {
                    spotifyAuthViewModel.startReauth(
                        onAuthUrl = { url ->
                            CustomTabsIntent.Builder().build()
                                .launchUrl(context, Uri.parse(url))
                        },
                        onError = { msg ->
                            snackbarState.showSnackbar("Failed to start Spotify reconnect: $msg")
                        }
                    )
                },
                modifier = Modifier.fillMaxWidth()
            ) {
                Text("Reconnect Spotify")
            }
        }
    }
}

@Composable
private fun StatusCard(status: SpotifyAuthStatus) {
    Card(modifier = Modifier.fillMaxWidth()) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            Text(
                text = if (status.authenticated) "Authenticated" else "Not Authenticated",
                style = MaterialTheme.typography.titleMedium
            )
            Spacer(modifier = Modifier.height(8.dp))
            val expireLabel = if (status.expireTs.epochSecond <= 0L) {
                "Unknown"
            } else {
                expireTsFormatter.format(status.expireTs)
            }
            Text(
                text = "Refresh token expires: $expireLabel",
                style = MaterialTheme.typography.bodyMedium
            )
        }
    }
}
