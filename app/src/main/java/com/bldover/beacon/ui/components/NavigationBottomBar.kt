package com.bldover.beacon.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.RowScope
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Build
import androidx.compose.material.icons.filled.DateRange
import androidx.compose.material.icons.filled.Favorite
import androidx.compose.material.icons.outlined.Build
import androidx.compose.material.icons.outlined.DateRange
import androidx.compose.material.icons.outlined.Favorite
import androidx.compose.material.icons.outlined.FavoriteBorder
import androidx.compose.material3.BottomAppBar
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.tooling.preview.Preview
import com.bldover.beacon.ActiveScreen
import com.bldover.beacon.ui.theme.BeaconTheme

@Composable
fun NavigationBottomBar(
    activeScreen: ActiveScreen,
    changeScreen: (ActiveScreen) -> Unit
) {
    BottomAppBar {
        NavigationBar {
            NavigationItem(
                navItemScreen = ActiveScreen.CONCERT_HISTORY,
                activeScreen = activeScreen,
                changeScreen = changeScreen,
                selectedIcon = Icons.Filled.Favorite,
                unselectedIcon = Icons.Outlined.FavoriteBorder
            )
            NavigationItem(
                navItemScreen = ActiveScreen.UPCOMING_EVENTS,
                activeScreen = activeScreen,
                changeScreen = changeScreen,
                selectedIcon = Icons.Filled.DateRange,
                unselectedIcon = Icons.Outlined.DateRange
            )
            NavigationItem(
                navItemScreen = ActiveScreen.UTILITIES,
                activeScreen = activeScreen,
                changeScreen = changeScreen,
                selectedIcon = Icons.Filled.Build,
                unselectedIcon = Icons.Outlined.Build
            )
        }
    }
}

@Composable
fun RowScope.NavigationItem(
    navItemScreen: ActiveScreen,
    activeScreen: ActiveScreen,
    changeScreen: (ActiveScreen) -> Unit,
    selectedIcon: ImageVector,
    unselectedIcon: ImageVector
) {
    val isSelected = activeScreen == navItemScreen
    NavigationBarItem(
        selected = isSelected,
        onClick = { changeScreen(navItemScreen) },
        label = { Text(navItemScreen.shortDesc) },
        icon = {
            Icon(
                imageVector = if (isSelected) selectedIcon else unselectedIcon,
                contentDescription = "${navItemScreen.title} Screen Switch"
            )
        }
    )
}

@Preview
@Composable
fun NavigationBottomBarPreview(
    activeScreen: ActiveScreen = ActiveScreen.CONCERT_HISTORY
) {
    BeaconTheme(darkTheme = true, dynamicColor = true) {
        NavigationBottomBar(activeScreen) {}
    }
}