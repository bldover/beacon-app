package com.bldover.beacon.data.model.ordering

abstract class Ordering<T>(
    val option: OrderField,
    val order: Direction
) : Comparator<T>