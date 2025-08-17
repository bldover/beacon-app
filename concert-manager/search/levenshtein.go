package search

import (
	"math"
	"strings"
)

// Two-row optimized Levenshtein distance algorithm to find the minimum number of
// inserts, deletions, or substitutions required to make two case-insensitive strings match
func getLevenshteinDistance(x string, y string) int {
	xRunes := []rune(x)
	yRunes := []rune(y)
	xLen, yLen := len(xRunes), len(yRunes)

	currRow := make([]int, yLen+1)
	prevRow := make([]int, yLen+1)
	for j := 0; j <= yLen; j++ {
		prevRow[j] = j
	}

	for i := 0; i < xLen; i++ {
		currRow[0] = i + 1
		for j := 0; j < yLen; j++ {
			ind := 0
			if !strings.EqualFold(string(xRunes[i]), string(yRunes[j])) {
				ind = 1
			}
			delDist := prevRow[j+1] + 1
			insDist := currRow[j] + 1
			subDist := prevRow[j] + ind

			levDist := int(math.Min(math.Min(float64(delDist), float64(insDist)), float64(subDist)))
			currRow[j+1] = levDist
		}
		prevRow, currRow = currRow, prevRow
	}
	return prevRow[yLen]
}
