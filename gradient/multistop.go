package gradient

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// MultiStopLCH interpolates through all stops in LCH space, returning n total colors.
// Both outer endpoints are included. Returns nil if n <= 0 or stops is empty.
func MultiStopLCH(stops []golor.Color, n int) []golor.Color {
	return multiStop(stops, n, func(a, b golor.Color, t float64) golor.Color {
		la := convert.ToLCH(a)
		lb := convert.ToLCH(b)
		lch := convert.LCH{
			L: la.L + t*(lb.L-la.L),
			C: la.C + t*(lb.C-la.C),
			H: lerpHue(la.H, lb.H, t),
		}
		return convert.FromLCH(lch)
	})
}

// MultiStopRGB interpolates through all stops in RGB space, returning n total colors.
// Both outer endpoints are included. Returns nil if n <= 0 or stops is empty.
func MultiStopRGB(stops []golor.Color, n int) []golor.Color {
	return multiStop(stops, n, func(a, b golor.Color, t float64) golor.Color {
		return golor.RGBf(
			a.R+t*(b.R-a.R),
			a.G+t*(b.G-a.G),
			a.B+t*(b.B-a.B),
		)
	})
}

func multiStop(stops []golor.Color, n int, lerp func(a, b golor.Color, t float64) golor.Color) []golor.Color {
	if n <= 0 || len(stops) == 0 {
		return nil
	}
	if len(stops) == 1 || n == 1 {
		return []golor.Color{stops[0]}
	}
	numSegs := len(stops) - 1
	result := make([]golor.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		scaled := t * float64(numSegs)
		seg := int(scaled)
		if seg >= numSegs {
			seg = numSegs - 1
		}
		lt := scaled - float64(seg)
		result[i] = lerp(stops[seg], stops[seg+1], lt)
	}
	return result
}
