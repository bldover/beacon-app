package util

import "testing"

func TestMatch(t *testing.T) {
	dist := getLevenshteinDistance("cat", "cat")
	expected := 0
	if dist != expected {
		t.Errorf("Incorrect distance, expected: %v, actual: %v", expected, dist)
	}
}

func TestMatchDifferentCases(t *testing.T) {
	dist := getLevenshteinDistance("cat", "CAT")
	expected := 0
	if dist != expected {
		t.Errorf("Incorrect distance, expected: %v, actual: %v", expected, dist)
	}
}

func TestSubstitution(t *testing.T) {
	dist := getLevenshteinDistance("cat", "dog")
	expected := 3
	if dist != expected {
		t.Errorf("Incorrect distance, expected: %v, actual: %v", expected, dist)
	}
}

func TestInsertion(t *testing.T) {
	dist := getLevenshteinDistance("cat", "catty")
	expected := 2
	if dist != expected {
		t.Errorf("Incorrect distance, expected: %v, actual: %v", expected, dist)
	}
}

func TestDeletion(t *testing.T) {
	dist := getLevenshteinDistance("cat", "ca")
	expected := 1
	if dist != expected {
		t.Errorf("Incorrect distance, expected: %v, actual: %v", expected, dist)
	}
}

func TestMixed(t *testing.T) {
	dist := getLevenshteinDistance("intention", "execution")
	expected := 5
	if dist != expected {
		t.Errorf("Incorrect distance, expected: %v, actual: %v", expected, dist)
	}
}
