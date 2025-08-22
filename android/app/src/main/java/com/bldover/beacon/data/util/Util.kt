package com.bldover.beacon.data.util

import java.time.format.DateTimeFormatter

val dateFormatter: DateTimeFormatter = DateTimeFormatter.ofPattern("M/d/yyyy")

fun toCommaSeparatedString(list: List<String>): String {
    return if (list.isEmpty()) ""
    else list.reduce { acc, s -> "$acc, $s" }
}

fun fromCommaSeparatedString(string: String): List<String> {
    return string.replace(" ", "").split(",")
}