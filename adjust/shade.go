package adjust

import (
	"github.com/0mega24/golor"
)

// Shade mixes c toward black by amount (0=original, 1=black).
func Shade(c golor.Color, amount float64) golor.Color {
	t := clamp01(amount)
	return golor.RGBf(
		c.R*(1-t),
		c.G*(1-t),
		c.B*(1-t),
	)
}
