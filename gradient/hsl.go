package gradient

import (
	"math"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// HSL returns n evenly-spaced colors interpolated in HSL space between a and b.
// Hue takes the shortest path around the hue wheel. Both endpoints are included.
func HSL(a, b golor.Color, n int) []golor.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []golor.Color{a}
	}
	ha := convert.ToHSL(a)
	hb := convert.ToHSL(b)
	result := make([]golor.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		hsl := convert.HSL{
			H: lerpHue(ha.H, hb.H, t),
			S: ha.S + t*(hb.S-ha.S),
			L: ha.L + t*(hb.L-ha.L),
		}
		result[i] = convert.FromHSL(hsl)
	}
	return result
}

// lerpHue interpolates between two hue angles taking the shortest path (±180°).
func lerpHue(h1, h2, t float64) float64 {
	diff := h2 - h1
	for diff > 180 {
		diff -= 360
	}
	for diff < -180 {
		diff += 360
	}
	return math.Mod(h1+t*diff+360*10, 360)
}
