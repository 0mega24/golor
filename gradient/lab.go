package gradient

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// LAB returns n evenly-spaced colors interpolated in LAB space between a and b.
// Both endpoints are included.
func LAB(a, b golor.Color, n int) []golor.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []golor.Color{a}
	}
	la := convert.ToLAB(a)
	lb := convert.ToLAB(b)
	result := make([]golor.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		lab := convert.LAB{
			L: la.L + t*(lb.L-la.L),
			A: la.A + t*(lb.A-la.A),
			B: la.B + t*(lb.B-la.B),
		}
		result[i] = convert.FromLAB(lab)
	}
	return result
}
