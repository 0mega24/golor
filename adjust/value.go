package adjust

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// SetValue sets the HSV value (brightness) of c to v (clamped to [0,1]) and preserves alpha.
func SetValue(c golor.Color, v float64) golor.Color {
	hsv := convert.ToHSV(c)
	hsv.V = clamp01(v)
	out := convert.FromHSV(hsv)
	out.A = c.A
	return out
}
