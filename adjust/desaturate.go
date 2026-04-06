package adjust

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// Desaturate decreases the HSL saturation of c by amount (additive, clamped to [0,1]).
func Desaturate(c golor.Color, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.S = clamp01(hsl.S - amount)
	return convert.FromHSL(hsl)
}
