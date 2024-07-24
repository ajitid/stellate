package main

import (
	"math/rand/v2"
)

// [min, max)
func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func clamp(min, max, val int) int {
	if val > max {
		return max
	} else if val < min {
		return min
	}
	return val
}
