package com.bldover.beacon.ui.screens.utility

import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import java.time.Duration
import java.time.Instant

private val REAUTH_PROMPT_WINDOW: Duration = Duration.ofDays(14)

@Composable
fun SpotifyAuthNotification(
    spotifyAuthViewModel: SpotifyAuthViewModel
) {
    val authStatus by spotifyAuthViewModel.authStatus.collectAsState()
    val initialPromptResolved by spotifyAuthViewModel.initialPromptResolved.collectAsState()

    if (initialPromptResolved || !needsSpotifyReauthPrompt(authStatus)) return

    AlertDialog(
        onDismissRequest = { spotifyAuthViewModel.markInitialPromptResolved() },
        title = { Text("Spotify Reauth Required") },
        text = {
            Text(
                "Your Spotify connection is missing or expires soon. " +
                    "Open Utilities > Spotify Auth to reconnect."
            )
        },
        confirmButton = {
            Button(onClick = { spotifyAuthViewModel.markInitialPromptResolved() }) {
                Text("Close")
            }
        }
    )
}

private fun needsSpotifyReauthPrompt(state: SpotifyAuthStatusState): Boolean {
    if (state !is SpotifyAuthStatusState.Loaded) return false
    val status = state.status
    if (!status.authenticated) return true
    return Duration.between(Instant.now(), status.expireTs) <= REAUTH_PROMPT_WINDOW
}
