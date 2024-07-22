package com.bldover.beacon

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.requiredHeight
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AddCircle
import androidx.compose.material.icons.filled.Build
import androidx.compose.material.icons.filled.Star
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Text
import androidx.compose.material3.VerticalDivider
import androidx.compose.runtime.Composable
import androidx.compose.runtime.MutableState
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.bldover.beacon.ui.theme.BeaconTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            BeaconTheme {
                Application()
            }
        }
    }
}

@Composable
fun Application() {
    var currScreen = remember {
        mutableStateOf("Saved Events")
    }
    Column(
        modifier = Modifier.fillMaxSize(),
        verticalArrangement = Arrangement.Bottom
    ) {
        Box(modifier = Modifier
            .weight(1f)
            .fillMaxWidth()
            .background(Color.DarkGray)
        ) {
            when (currScreen.value) {
                "Saved Events" -> SavedEvents()
                "Recommendations" -> Recommendations()
                "Utilities" -> Utilities()
            }
        }
        BottomNavigationBar(currScreen)
    }
}

@Composable
fun BottomNavigationBar(
    currScreen: MutableState<String>
) {
    Box(
        modifier = Modifier
            .background(Color.Red)
            .requiredHeight(65.dp)
    ) {
        NavigationBar(
            modifier = Modifier.fillMaxWidth()
        ) {
            Row(
                verticalAlignment = Alignment.Bottom
            ) {
                NavigationBarItem(
                    selected = true,
                    onClick = { currScreen.value = "Saved Events" },
                    icon = { Icon(imageVector = Icons.Default.Star, contentDescription = null) })
                VerticalDivider()
                NavigationBarItem(
                    selected = false,
                    onClick = { currScreen.value = "Recommendations" },
                    icon = { Icon(imageVector = Icons.Default.AddCircle, contentDescription = null) })
                VerticalDivider()
                NavigationBarItem(
                    selected = false,
                    onClick = { currScreen.value = "Utilities" },
                    icon = { Icon(imageVector = Icons.Default.Build, contentDescription = null) })
            }
        }
    }
}

@Composable
fun SavedEvents() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = "SavedEvents",
            color = Color.White
        )
    }
}

@Composable
fun Recommendations() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = "Recommendations",
            color = Color.White
        )
    }
}

@Composable
fun Utilities() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = "Utilities",
            color = Color.White
        )
    }
}

@Preview
@Composable
fun DefaultPreview() {
    BeaconTheme {
        Application()
    }
}