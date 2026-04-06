package gradient

import "github.com/0mega24/golor"

// GradientRGB returns n evenly-spaced colors interpolated in RGB space between a and b.
// Both endpoints are included. Returns nil if n <= 0.
func GradientRGB(a, b golor.Color, n int) []golor.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []golor.Color{a}
	}
	result := make([]golor.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		result[i] = golor.RGBf(
			a.R+t*(b.R-a.R),
			a.G+t*(b.G-a.G),
			a.B+t*(b.B-a.B),
		)
	}
	return result
}
