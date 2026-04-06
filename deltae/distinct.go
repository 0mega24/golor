package deltae

import (
	"math"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// EnsureDistinct adjusts colors in the slice so each pair has at least minDeltaE
// perceptual difference (using DeltaE76). Returns the adjusted slice.
func EnsureDistinct(colors []golor.Color, minDeltaE float64) []golor.Color {
	result := make([]golor.Color, len(colors))
	copy(result, colors)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			for DeltaE76(result[i], result[j]) < minDeltaE {
				lchi := convert.ToLCH(result[i])
				lchj := convert.ToLCH(result[j])
				if lchi.L >= lchj.L {
					lchj.L = math.Max(0, lchj.L-minDeltaE)
				} else {
					lchj.L = math.Min(100, lchj.L+minDeltaE)
				}
				prev := result[j]
				result[j] = convert.FromLCH(lchj)
				if result[j] == prev {
					break
				}
			}
		}
	}
	return result
}
