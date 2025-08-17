package com.bldover.beacon.data.model

abstract class Ordering<T>(
    val option: OrderField,
    val order: Direction
) : Comparator<T>