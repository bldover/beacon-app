package com.bldover.beacon.ui.components.common

import androidx.compose.animation.core.animateOffsetAsState
import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.offset
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.itemsIndexed
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.layout.onGloballyPositioned
import androidx.compose.ui.layout.positionInRoot
import androidx.compose.ui.unit.IntOffset
import androidx.compose.ui.zIndex
import kotlin.math.roundToInt

@Composable
fun <T> ReorderableList(
    items: List<T>,
    onReorder: (from: Int, to: Int) -> Unit,
    itemContent: @Composable (T) -> Unit
) {
    var draggingItemIndex by remember { mutableStateOf<Int?>(null) }
    var dragOffset by remember { mutableStateOf(Offset.Zero) }
    val itemPositions = remember { mutableStateListOf<Offset>() }
    val itemSizes = remember { mutableStateListOf<Int>() }

    LazyColumn(
        modifier = Modifier.fillMaxWidth()
    ) {
        itemsIndexed(items) { index, item ->
            val isDragging = index == draggingItemIndex
            val itemOffset = if (isDragging) {
                dragOffset
            } else {
                val offsetForShuffling = if (draggingItemIndex != null && index > draggingItemIndex!!) {
                    Offset(0f, -itemSizes.getOrElse(draggingItemIndex!!) { 0 }.toFloat())
                } else {
                    Offset.Zero
                }
                offsetForShuffling
            }

            val animatedOffset by animateOffsetAsState(itemOffset)

            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .offset { IntOffset(animatedOffset.x.roundToInt(), animatedOffset.y.roundToInt()) }
                    .zIndex(if (isDragging) 1f else 0f)
                    .onGloballyPositioned { coordinates ->
                        if (itemPositions.size <= index) {
                            itemPositions.add(coordinates.positionInRoot())
                        } else {
                            itemPositions[index] = coordinates.positionInRoot()
                        }
                        if (itemSizes.size <= index) {
                            itemSizes.add(coordinates.size.height)
                        } else {
                            itemSizes[index] = coordinates.size.height
                        }
                    }
                    .pointerInput(Unit) {
                        detectDragGestures(
                            onDragStart = { draggingItemIndex = index },
                            onDrag = { change, offset ->
                                change.consume()
                                dragOffset += offset
                            },
                            onDragEnd = {
                                val fromIndex = draggingItemIndex!!
                                val toIndex = itemPositions.indexOfFirst { pos ->
                                    (dragOffset + (itemPositions.getOrNull(fromIndex) ?: Offset.Zero)).y < pos.y + itemSizes.getOrElse(fromIndex) { 0 } / 2
                                }.let { if (it == -1) items.lastIndex else it }
                                onReorder(fromIndex, toIndex)
                                draggingItemIndex = null
                                dragOffset = Offset.Zero
                            }
                        )
                    }
            ) {
                itemContent(item)
            }
        }
    }
}