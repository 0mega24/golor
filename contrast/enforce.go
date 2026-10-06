package contrast

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// EnforceContrast moves fg toward black/white, choosing the nearest passing color in OKLab.
// Adjusted colors preserve alpha at 8-bit precision and meet minRatio after rounding.
// Unchanged if already passing; best effort if unreachable (see EnforceContrastOK).
func EnforceContrast(fg, bg golor.Color, minRatio float64) golor.Color {
	c, _ := EnforceContrastOK(fg, bg, minRatio)
	return c
}

// EnforceContrastOK applies EnforceContrast and reports whether the result meets minRatio.
func EnforceContrastOK(fg, bg golor.Color, minRatio float64) (golor.Color, bool) {
	if Meets(fg, bg, minRatio) {
		return fg, true
	}
	backdrop := composite(bg, golor.RGB(255, 255, 255))
	bgLum := Luminance(backdrop)
	ratio := func(c golor.Color) float64 {
		l1, l2 := Luminance(composite(c, backdrop)), bgLum
		if l1 < l2 {
			l1, l2 = l2, l1
		}
		return (l1 + 0.05) / (l2 + 0.05)
	}
	alpha := fg.A8()
	point := func(target, t float64) golor.Color {
		return golor.RGBAf(fg.R+t*(target-fg.R), fg.G+t*(target-fg.G), fg.B+t*(target-fg.B), fg.A)
	}
	quantize := func(c golor.Color) golor.Color { return golor.RGBA(c.R8(), c.G8(), c.B8(), alpha) }

	var best golor.Color
	bestDistance, bestRatio := -1.0, -1.0
	var fallback golor.Color
	for _, target := range [2]float64{1, 0} {
		end := quantize(point(target, 1))
		if r := ratio(end); r < minRatio {
			if r > bestRatio {
				fallback, bestRatio = end, r
			}
			continue
		}
		// Along a line toward white or black, luminance is monotonic, so "meets" flips
		// from false to true exactly once and bisection finds the boundary.
		lo, hi := 0.0, 1.0
		for range 24 {
			mid := (lo + hi) / 2
			if ratio(point(target, mid)) >= minRatio {
				hi = mid
			} else {
				lo = mid
			}
		}
		c := quantize(point(target, hi))
		for t := hi; ratio(c) < minRatio; {
			t = min(1, t+1.0/255)
			c = quantize(point(target, t))
		}
		if d := convert.DeltaEOK(c, fg); bestDistance < 0 || d < bestDistance {
			best, bestDistance = c, d
		}
	}
	if bestDistance >= 0 {
		return best, true
	}
	return fallback, false
}
