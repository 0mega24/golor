// Package blend provides photographic blend mode functions for compositing two colors.
package blend

import "github.com/0mega24/golor"

// Mix linearly interpolates between a and b. t=0 returns a, t=1 returns b.
func Mix(a, b golor.Color, t float64) golor.Color {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return golor.RGBf(
		a.R+t*(b.R-a.R),
		a.G+t*(b.G-a.G),
		a.B+t*(b.B-a.B),
	)
}
