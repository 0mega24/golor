package adjust

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// SetValue sets the HSV value (brightness) of c to v (clamped to [0,1]).
func SetValue(c golor.Color, v float64) golor.Color {
	hsv := convert.ToHSV(c)
	hsv.V = clamp01(v)
	return convert.FromHSV(hsv)
}
