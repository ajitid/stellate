package main

import (
	"math"
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

// snapResult is a struct to hold the result of the SnapSlice function
type snapResult struct {
	Point float64
	Index int
}

// snapNumber takes a float64 and returns a function that snaps a value to the nearest multiple of the input
func snapNumber(point float64) func(float64) float64 {
	return func(v float64) float64 {
		return math.Round(v/point) * point
	}
}

// snapSlice takes a slice of float64 and returns a function that snaps a value to the nearest point in the slice
func snapSlice(points []float64) func(float64) snapResult {
	return func(v float64) snapResult {
		if len(points) == 0 {
			return snapResult{Point: 0, Index: -1}
		}

		lastDistance := math.Abs(points[0] - v)
		result := snapResult{Point: points[0], Index: 0}

		for i := 1; i < len(points); i++ {
			distance := math.Abs(points[i] - v)

			if distance == 0 {
				return snapResult{Point: points[i], Index: i}
			}

			if distance > lastDistance {
				return result
			}

			result = snapResult{Point: points[i], Index: i}
			lastDistance = distance
		}

		return result // return the last item as the result
	}
}
