package main

import "math/rand/v2"

// [min, max)
func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
