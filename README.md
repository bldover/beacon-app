# Beacon

Beacon is a personal-use app utility to manage concerts.
- Keep track of attended concert history and upcoming shows
- Recommend local shows based on listening history

The concert-manager module contains the backend, written in Go.
- GCP Firestore for persistence, with an in-memory cache layer to reduce DB operations so I don't exceed the GCP free tier limits
- CSV upload for batch loading concert history
- Ticketmaster API provides upcoming local concerts and metadata
- Spotify API exposes my music listening history and song frequency metrics to build baseline rankings for known artists
- Tidal API correlates known artists with related artists to enable generating ratings for unknown artists
- Terminal UI framework and TUI frontend, which was the first iteration of the interactive app

The android module contains the Android frontend app, written in Kotlin with Jetpack Compose.
