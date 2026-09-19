// Package gradient provides color gradient interpolation in RGB, HSL, L*a*b*, and LCH color spaces.
package gradient

import "github.com/0mega24/golor/v2"

// RGB returns n evenly-spaced colors interpolated in RGB space between a and b.
// Alpha is interpolated linearly alongside RGB.
// Both endpoints are included. Returns nil if n <= 0.
func RGB(a, b golor.Color, n int) []golor.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []golor.Color{a}
	}
	result := make([]golor.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		result[i] = golor.RGBAf(
			a.R+t*(b.R-a.R),
			a.G+t*(b.G-a.G),
			a.B+t*(b.B-a.B),
			a.A+t*(b.A-a.A),
		)
	}
	return result
}
