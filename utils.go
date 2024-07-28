package main

import (
	"log"
	"math"
	"math/rand/v2"

	"github.com/adrg/sysfont"
	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
)

func getSystemFontPath() string {
	finder := sysfont.NewFinder(&sysfont.FinderOpts{
		Extensions: []string{".ttf"},
	})
	font := finder.Match("Segoe UI")
	if font == nil {
		log.Fatal("system font not found")
	}
	return font.Filename
}

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

// snapNumber takes a float64 and returns a function that snaps a value to the nearest multiple of the input
func snapNumber[T constraints.Float](point T) func(T) T {
	return func(v T) T {
		return T(math.Round(float64(v/point))) * point
	}
}

// snapResult is a struct to hold the result of the SnapSlice function
type snapResult[T constraints.Float] struct {
	point T
	index int
}

// snapSlice takes a slice of T and returns a function that snaps a value to the nearest point in the slice
func snapSlice[T constraints.Float](points []T) func(T) snapResult[T] {
	return func(v T) snapResult[T] {
		if len(points) == 0 {
			return snapResult[T]{point: 0, index: -1}
		}

		lastDistance := T(math.Abs(float64(points[0] - v)))
		result := snapResult[T]{point: points[0], index: 0}

		for i := 1; i < len(points); i++ {
			distance := T(math.Abs(float64(points[i] - v)))

			if distance == 0 {
				return snapResult[T]{point: points[i], index: i}
			}

			if distance > lastDistance {
				return result
			}

			result = snapResult[T]{point: points[i], index: i}
			lastDistance = distance
		}

		return result // return the last item as the result
	}
}

func mapRange[T constraints.Float](value, fromLow, fromHigh, toLow, toHigh T) T {
	return (value-fromLow)*(toHigh-toLow)/(fromHigh-fromLow) + toLow
}

func drawLinesAroundCircle(center rl.Vector2, radius float32, lineCount int, lineLength float32, color rl.Color) {
	for i := 0; i < lineCount; i++ {
		angle := float32(i) / float32(lineCount) * 2 * math.Pi
		start := rl.Vector2{
			X: center.X + float32(math.Cos(float64(angle)))*radius,
			Y: center.Y + float32(math.Sin(float64(angle)))*radius,
		}
		end := rl.Vector2{
			X: center.X + float32(math.Cos(float64(angle)))*(radius+lineLength),
			Y: center.Y + float32(math.Sin(float64(angle)))*(radius+lineLength),
		}
		rl.DrawLineV(start, end, color)
	}
}
