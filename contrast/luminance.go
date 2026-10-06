// Package contrast provides WCAG 2.1 luminance, contrast ratio, and enforcement functions,
// plus the APCA contrast metric used by WCAG 3 drafts.
package contrast

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// Luminance returns WCAG relative luminance in [0,1] from clamped RGB; Ratio handles alpha.
func Luminance(c golor.Color) float64 {
	return 0.2126*wcagLinearize(c.R) + 0.7152*wcagLinearize(c.G) + 0.0722*wcagLinearize(c.B)
}

// wcagLinearize is the sRGB transfer function (WCAG 2.x uses the 0.04045 breakpoint).
func wcagLinearize(v float64) float64 { return space.Linearize(space.Clamp01(v)) }
