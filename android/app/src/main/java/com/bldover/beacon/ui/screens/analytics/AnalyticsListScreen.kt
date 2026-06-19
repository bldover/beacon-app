package com.bldover.beacon.ui.screens.analytics

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.data.model.analytics.AnalyticsCategory
import com.bldover.beacon.ui.components.analytics.AnalyticsCountCard
import com.bldover.beacon.ui.components.analytics.AnalyticsSortBar
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.LoadErrorMessage
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.ScrollableItemList
import com.bldover.beacon.ui.components.common.TitleTopBar
import timber.log.Timber

@Composable
fun AnalyticsListScreen(
    navController: NavController,
    category: AnalyticsCategory,
    title: String,
    listViewModel: AnalyticsListViewModel,
    eventsViewModel: AnalyticsEventsViewModel
) {
    Timber.d("composing AnalyticsListScreen for $category")
    LaunchedEffect(category) { listViewModel.load(category) }
    val state by listViewModel.stateFor(category).collectAsState()
    val ordering by listViewModel.orderingFor(category).collectAsState()
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = title,
                leadingIcon = { BackButton(navController = navController) }
            )
        }
    ) {
        Scaffold(
            topBar = {
                AnalyticsSortBar(
                    state = ordering,
                    onChange = { listViewModel.sort(category, it) }
                )
            }
        ) { innerPadding ->
            Column(modifier = Modifier.padding(innerPadding)) {
                Spacer(modifier = Modifier.padding(8.dp))
                when (val s = state) {
                    is AnalyticsListState.Loading -> LoadingSpinner()
                    is AnalyticsListState.Error -> LoadErrorMessage()
                    is AnalyticsListState.Success -> {
                        ScrollableItemList(
                            items = s.sorted,
                            getItemKey = { it.key }
                        ) { item ->
                            AnalyticsCountCard(
                                count = item,
                                onClick = {
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
        }
    }
}

@Composable
fun AnalyticsYearsScreen(
    navController: NavController,
    listViewModel: AnalyticsListViewModel,
    eventsViewModel: AnalyticsEventsViewModel
) = AnalyticsListScreen(
    navController = navController,
    category = AnalyticsCategory.YEARS,
    title = Screen.ANALYTICS_YEARS.title,
    listViewModel = listViewModel,
    eventsViewModel = eventsViewModel
)

@Composable
fun AnalyticsMonthsScreen(
    navController: NavController,
    listViewModel: AnalyticsListViewModel,
    eventsViewModel: AnalyticsEventsViewModel
) = AnalyticsListScreen(
    navController = navController,
    category = AnalyticsCategory.MONTHS,
    title = Screen.ANALYTICS_MONTHS.title,
    listViewModel = listViewModel,
    eventsViewModel = eventsViewModel
)

@Composable
fun AnalyticsArtistsScreen(
    navController: NavController,
    listViewModel: AnalyticsListViewModel,
    eventsViewModel: AnalyticsEventsViewModel
) = AnalyticsListScreen(
    navController = navController,
    category = AnalyticsCategory.ARTISTS,
    title = Screen.ANALYTICS_ARTISTS.title,
    listViewModel = listViewModel,
    eventsViewModel = eventsViewModel
)

@Composable
fun AnalyticsVenuesScreen(
    navController: NavController,
    listViewModel: AnalyticsListViewModel,
    eventsViewModel: AnalyticsEventsViewModel
) = AnalyticsListScreen(
    navController = navController,
    category = AnalyticsCategory.VENUES,
    title = Screen.ANALYTICS_VENUES.title,
    listViewModel = listViewModel,
    eventsViewModel = eventsViewModel
)

@Composable
fun AnalyticsGenresScreen(
    navController: NavController,
    listViewModel: AnalyticsListViewModel,
    eventsViewModel: AnalyticsEventsViewModel
) = AnalyticsListScreen(
    navController = navController,
    category = AnalyticsCategory.GENRES,
    title = Screen.ANALYTICS_GENRES.title,
    listViewModel = listViewModel,
    eventsViewModel = eventsViewModel
)
