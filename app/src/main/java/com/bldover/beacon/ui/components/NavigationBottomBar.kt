package com.bldover.beacon.ui.components

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
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.navigation.NavController
import androidx.navigation.compose.currentBackStackEntryAsState
import com.bldover.beacon.data.model.Screen

@Composable
fun NavigationBottomBar(navController: NavController) {
    val activeScreen = Screen.fromOrDefault(navController.currentBackStackEntryAsState().value?.destination?.route)
    BottomAppBar {
        NavigationBar {
            NavigationItem(
                label = "Planner",
                isSelected = Screen.CONCERT_PLANNER == activeScreen,
                selectedIcon = Icons.Filled.Edit,
                unselectedIcon = Icons.Outlined.Edit,
            ) {
                navigateAndPopAll(navController, Screen.CONCERT_PLANNER)
            }
            NavigationItem(
                label = "Upcoming",
                isSelected = Screen.UPCOMING_EVENTS == activeScreen,
                selectedIcon = Icons.Filled.DateRange,
                unselectedIcon = Icons.Outlined.DateRange,
            ) {
                navigateAndPopAll(navController, Screen.UPCOMING_EVENTS)
            }
            NavigationItem(
                label = "History",
                isSelected = Screen.CONCERT_HISTORY == activeScreen,
                selectedIcon = Icons.Filled.Favorite,
                unselectedIcon = Icons.Outlined.FavoriteBorder,
            ) {
                navigateAndPopAll(navController, Screen.CONCERT_HISTORY)
            }
            NavigationItem(
                label = "Utilities",
                isSelected = Screen.UTILITIES == activeScreen,
                selectedIcon = Icons.Filled.Build,
                unselectedIcon = Icons.Outlined.Build
            ) {
                navigateAndPopAll(navController, Screen.UTILITIES)
            }
        }
    }
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

@Composable
fun RowScope.NavigationItem(
    label: String,
    isSelected: Boolean,
    selectedIcon: ImageVector,
    unselectedIcon: ImageVector,
    onClick: () -> Unit
) {
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