package com.bldover.beacon.ui.screens.analytics

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.analytics.AnalyticsCategory
import com.bldover.beacon.data.model.analytics.AnalyticsCount
import com.bldover.beacon.data.model.analytics.AnalyticsSummary
import com.bldover.beacon.ui.components.analytics.AnalyticsTopSection
import com.bldover.beacon.ui.components.common.BasicCard
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import timber.log.Timber

@Composable
fun AnalyticsOverviewScreen(
    navController: NavController,
    overviewViewModel: AnalyticsOverviewViewModel,
    eventsViewModel: AnalyticsEventsViewModel
) {
    Timber.d("composing AnalyticsOverviewScreen")
    val state by overviewViewModel.state.collectAsState()
    ScreenFrame(
        topBar = { TitleTopBar(title = Screen.ANALYTICS.title) }
    ) {
        when (state) {
            is AnalyticsOverviewState.Loading -> LoadingSpinner()
            is AnalyticsOverviewState.Error -> LoadErrorMessage()
            is AnalyticsOverviewState.Success -> {
                OverviewContent(
                    summary = (state as AnalyticsOverviewState.Success).summary,
                    onSeeAll = { category ->
                        navController.navigate(category.toScreen().name)
                    },
                    onItemClick = { category, item ->
                        eventsViewModel.launch(
                            navController = navController,
                            category = category,
                            key = item.key,
                            name = item.name
                        )
                    }
                )
            }
        }
    }
}

@Composable
private fun OverviewContent(
    summary: AnalyticsSummary,
    onSeeAll: (AnalyticsCategory) -> Unit,
    onItemClick: (AnalyticsCategory, AnalyticsCount) -> Unit
) {
    val sections = listOf(
        Triple(AnalyticsCategory.YEARS, "Top Years", summary.topYears),
        Triple(AnalyticsCategory.MONTHS, "Top Months", summary.topMonths),
        Triple(AnalyticsCategory.ARTISTS, "Top Artists", summary.topArtists),
        Triple(AnalyticsCategory.VENUES, "Top Venues", summary.topVenues),
        Triple(AnalyticsCategory.GENRES, "Top Genres", summary.topGenres)
    )
    LazyColumn(
        modifier = Modifier.fillMaxWidth(),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        item { TotalEventsCard(summary.totalEvents) }
        items(items = sections, key = { it.first.name }) { (category, title, items) ->
            AnalyticsTopSection(
                title = title,
                items = items,
                onSeeAll = { onSeeAll(category) },
                onItemClick = { item -> onItemClick(category, item) }
            )
        }
    }
}

@Composable
private fun TotalEventsCard(total: Int) {
    BasicCard(modifier = Modifier.fillMaxWidth()) {
        Column(modifier = Modifier.padding(vertical = 8.dp)) {
            Text(
                text = "Total Concerts",
                style = MaterialTheme.typography.titleMedium
            )
            Text(
                text = total.toString(),
                style = MaterialTheme.typography.headlineMedium
            )
        }
    }
}

private fun AnalyticsCategory.toScreen(): Screen = when (this) {
    AnalyticsCategory.YEARS -> Screen.ANALYTICS_YEARS
    AnalyticsCategory.MONTHS -> Screen.ANALYTICS_MONTHS
    AnalyticsCategory.ARTISTS -> Screen.ANALYTICS_ARTISTS
    AnalyticsCategory.VENUES -> Screen.ANALYTICS_VENUES
    AnalyticsCategory.GENRES -> Screen.ANALYTICS_GENRES
}
