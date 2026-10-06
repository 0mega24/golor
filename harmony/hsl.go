// Package harmony provides color harmony generators based on hue relationships.
package harmony

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/adjust"
)

// Complementary returns c and its complement (hue rotated 180°).
func Complementary(c golor.Color) [2]golor.Color {
	return [2]golor.Color{c, adjust.ShiftHue(c, 180)}
}

// Triadic returns c and two colors equally spaced around the hue wheel (120° apart).
func Triadic(c golor.Color) [3]golor.Color {
	return [3]golor.Color{c, adjust.ShiftHue(c, 120), adjust.ShiftHue(c, 240)}
}

// Analogous returns three colors: c shifted left by angle, c itself, and c shifted right by angle.
// angle is the spread in degrees (e.g. 30).
func Analogous(c golor.Color, angle float64) [3]golor.Color {
	return [3]golor.Color{adjust.ShiftHue(c, -angle), c, adjust.ShiftHue(c, angle)}
}

// SplitComplementary returns c and two colors flanking its complement (±30° from 180°).
func SplitComplementary(c golor.Color) [3]golor.Color {
	return [3]golor.Color{c, adjust.ShiftHue(c, 150), adjust.ShiftHue(c, 210)}
}

// Tetradic returns four colors equally spaced around the hue wheel (90° apart).
func Tetradic(c golor.Color) [4]golor.Color {
	return [4]golor.Color{c, adjust.ShiftHue(c, 90), adjust.ShiftHue(c, 180), adjust.ShiftHue(c, 270)}
}
