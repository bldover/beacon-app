package com.bldover.beacon.ui.components.common

import androidx.compose.foundation.layout.RowScope
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Build
import androidx.compose.material.icons.filled.DateRange
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.filled.Favorite
import androidx.compose.material.icons.outlined.Build
import androidx.compose.material.icons.outlined.DateRange
import androidx.compose.material.icons.outlined.Edit
import androidx.compose.material.icons.outlined.FavoriteBorder
import androidx.compose.material3.BottomAppBar
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.derivedStateOf
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.navigation.NavController
import androidx.navigation.compose.currentBackStackEntryAsState
import com.bldover.beacon.data.model.Screen
import timber.log.Timber

@Composable
fun NavigationBottomBar(navController: NavController) {
    Timber.d("composing NavigationBottomBar")
    val currentBackStackEntry by navController.currentBackStackEntryAsState()
    val activeScreen by remember(currentBackStackEntry) {
        derivedStateOf {
            Screen.fromOrDefault(currentBackStackEntry?.destination?.route)
        }
    }

    BottomAppBar {
        NavigationBar {
            val navigationItems = remember {
                listOf(
                    NavigationItemData("Planner", Screen.CONCERT_PLANNER, Icons.Filled.Edit, Icons.Outlined.Edit),
                    NavigationItemData("Upcoming", Screen.UPCOMING_EVENTS, Icons.Filled.DateRange, Icons.Outlined.DateRange),
                    NavigationItemData("History", Screen.CONCERT_HISTORY, Icons.Filled.Favorite, Icons.Outlined.FavoriteBorder),
                    NavigationItemData("Utilities", Screen.UTILITIES, Icons.Filled.Build, Icons.Outlined.Build)
                )
            }

            navigationItems.forEach { item ->
                NavigationItem(
                    label = item.label,
                    isSelected = item.screen == activeScreen,
                    selectedIcon = item.selectedIcon,
                    unselectedIcon = item.unselectedIcon
                ) {
                    Timber.d("NavigationItem ${item.label.lowercase()} selected")
                    navigateAndPopAll(navController, item.screen)
                }
            }
        }
    }
}

data class NavigationItemData(
    val label: String,
    val screen: Screen,
    val selectedIcon: ImageVector,
    val unselectedIcon: ImageVector
)

@Composable
fun RowScope.NavigationItem(
    label: String,
    isSelected: Boolean,
    selectedIcon: ImageVector,
    unselectedIcon: ImageVector,
    onClick: () -> Unit
) {
    Timber.d("composing NavigationItem - $label")
    NavigationBarItem(
        selected = isSelected,
        onClick = { if (!isSelected) onClick() },
        label = { Text(label) },
        icon = {
            Icon(
                imageVector = if (isSelected) selectedIcon else unselectedIcon,
                contentDescription = "$label Screen Switch"
            )
        }
    )
}

fun navigateAndPopAll(
    navController: NavController,
    screen: Screen
) {
    navController.navigate(screen.name) {
        popUpTo(navController.graph.startDestinationId)
        launchSingleTop = true
    }
}