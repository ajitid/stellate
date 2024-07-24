package main

import (
	"math/rand/v2"
	"strconv"
)

// [min, max)
func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func main() {
	brightness := randRange(0, 101)
	setBrightness(strconv.Itoa(brightness))
}
