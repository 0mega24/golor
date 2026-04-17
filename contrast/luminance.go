// Package contrast provides WCAG 2.1 luminance, contrast ratio, and enforcement functions.
package contrast

import (
	"math"

	"github.com/0mega24/golor"
)

// Luminance returns the WCAG 2.1 relative luminance of c in [0, 1].
func Luminance(c golor.Color) float64 {
	return 0.2126*wcagLinearize(c.R) + 0.7152*wcagLinearize(c.G) + 0.0722*wcagLinearize(c.B)
}

func wcagLinearize(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}
