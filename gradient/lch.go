package gradient

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// LCH returns n evenly-spaced colors interpolated in LCH space between a and b.
// Hue takes the shortest path around the hue wheel. Both endpoints are included.
func LCH(a, b golor.Color, n int) []golor.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []golor.Color{a}
	}
	la := convert.ToLCH(a)
	lb := convert.ToLCH(b)
	result := make([]golor.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		lch := convert.LCH{
			L: la.L + t*(lb.L-la.L),
			C: la.C + t*(lb.C-la.C),
			H: lerpHue(la.H, lb.H, t),
		}
		result[i] = convert.FromLCH(lch)
	}
	return result
}
