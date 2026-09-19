// Package blend provides photographic blend mode functions for compositing two colors.
package blend

import "github.com/0mega24/golor/v2"

// Mix linearly interpolates between a and b, including alpha. t=0 returns a, t=1 returns b.
func Mix(a, b golor.Color, t float64) golor.Color {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return golor.RGBAf(
		a.R+t*(b.R-a.R),
		a.G+t*(b.G-a.G),
		a.B+t*(b.B-a.B),
		a.A+t*(b.A-a.A),
	)
}
