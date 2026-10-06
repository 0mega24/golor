package adjust

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// Tint mixes c toward white by amount (0=original, 1=white) and preserves alpha.
func Tint(c golor.Color, amount float64) golor.Color {
	t := space.Clamp01(amount)
	return golor.RGBAf(
		c.R+t*(1-c.R),
		c.G+t*(1-c.G),
		c.B+t*(1-c.B),
		c.A,
	)
}

// Shade mixes c toward black by amount (0=original, 1=black) and preserves alpha.
func Shade(c golor.Color, amount float64) golor.Color {
	t := space.Clamp01(amount)
	return golor.RGBAf(
		c.R*(1-t),
		c.G*(1-t),
		c.B*(1-t),
		c.A,
	)
}
