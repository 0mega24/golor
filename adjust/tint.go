package adjust

import (
	"github.com/0mega24/golor"
)

// Tint mixes c toward white by amount (0=original, 1=white).
func Tint(c golor.Color, amount float64) golor.Color {
	t := clamp01(amount)
	return golor.RGBf(
		c.R+t*(1-c.R),
		c.G+t*(1-c.G),
		c.B+t*(1-c.B),
	)
}
