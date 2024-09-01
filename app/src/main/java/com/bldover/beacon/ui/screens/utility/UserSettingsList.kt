package com.bldover.beacon.ui.screens.utility

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.Card
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.Switch
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.bldover.beacon.data.model.ColorScheme
import com.bldover.beacon.data.model.Screen
import com.bldover.beacon.ui.components.common.BackButton
import com.bldover.beacon.ui.components.common.LoadingSpinner
import com.bldover.beacon.ui.components.common.ScreenFrame
import com.bldover.beacon.ui.components.common.TitleTopBar
import timber.log.Timber

@Composable
fun UserSettingsScreen(
    navController: NavController,
    userSettingsViewModel: UserSettingsViewModel
) {
    Timber.d("composing UserSettingsScreen")
    ScreenFrame(
        topBar = {
            TitleTopBar(
                title = "Settings",
                leadingIcon = { BackButton(navController) }
            )
        }
    ) {
        UserSettingsList(userSettingsViewModel)
    }
}

@Composable
fun UserSettingsList(
    userSettingsViewModel: UserSettingsViewModel
) {
    when (val settings = userSettingsViewModel.userSettings.collectAsState().value) {
        is SettingsState.Loading -> {
            LoadingSpinner()
        }
        is SettingsState.Success -> {
            LazyColumn(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Top,
                modifier = Modifier.fillMaxSize()
            ) {
                item {
                    SelectorSetting(
                        text = "Color Scheme",
                        selectedOption = settings.data.colorScheme.scheme,
                        options = ColorScheme.entries.map { it.scheme }
                    ) {
                        userSettingsViewModel.updateColorScheme(ColorScheme.from(it))
                    }
                }
                item {
                    SelectorSetting(
                        text = "Start Screen",
                        selectedOption = Screen.fromOrDefault(settings.data.startScreen).title,
                        options = Screen.majorScreens().map { it.title }
                    ) {
                        userSettingsViewModel.updateStartScreen(Screen.fromTitle(it))
                    }
                }
            }
        }
    }
}

@Composable
fun ToggleSetting(
    text: String,
    isChecked: Boolean,
    onCheckedChange: (Boolean) -> Unit
) {
    var checked by remember { mutableStateOf(isChecked) }
    SettingCard {
        Row(
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.SpaceBetween,
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 16.dp)
        ) {
            Text(text = text)
            Switch(
                checked = isChecked,
                onCheckedChange = {
                    checked = it
                    onCheckedChange(it)
                }
            )
        }
    }
}

@Composable
fun SelectorSetting(
    text: String,
    selectedOption: String,
    options: List<String>,
    onOptionSelected: (String) -> Unit
) {
    SettingCard {
        Row(
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.SpaceBetween,
            modifier = Modifier
                .fillMaxSize()
                .padding(horizontal = 16.dp)
        ) {
            Text(text = text)
            DropdownOptions(
                options = options,
                selectedOption = selectedOption
            ) {
                onOptionSelected(it)
            }
        }
    }
}

@Composable
fun DropdownOptions(
    options: List<String>,
    selectedOption: String,
    onOptionSelected: (String) -> Unit
) {
    var expanded by remember { mutableStateOf(false) }
    var selectedIndex by remember { mutableIntStateOf(options.indexOf(selectedOption)) }

    Box {
        Text(
            text = options[selectedIndex],
            modifier = Modifier.clickable { expanded = true }
        )
        DropdownMenu(
            expanded = expanded,
            onDismissRequest = { expanded = false }
        ) {
            options.forEachIndexed { index, item ->
                DropdownMenuItem(
                    text = { Text(item) },
                    onClick = {
                        expanded = false
                        selectedIndex = index
                        onOptionSelected(item)
                    }
                )
            }
        }
    }
}

@Composable
fun SettingCard(content: @Composable () -> Unit) {
    Box(modifier = Modifier
        .fillMaxWidth()
        .height(48.dp)
    ) {
        Card(modifier = Modifier.fillMaxSize()) {
            content()
        }
    }
    Spacer(modifier = Modifier.height(16.dp))
}

@Preview
@Composable
fun SelectorPreview() {
    SelectorSetting(
        text = "Start Screen",
        selectedOption = "Upcoming Events",
        options = listOf("Upcoming Events", "Planner", "History"),
        onOptionSelected = {}
    )
}