package adjust

import (
	"math"

	"github.com/0mega24/golor/v2/internal/space"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// OKLCH adjustments preserve alpha and gamut-map results into sRGB.

// maxOKLCHChroma scales saturation amounts to the usable sRGB chroma range.
const maxOKLCHChroma = 0.4

func fromOKLCH(o convert.OKLCH, alpha float64) golor.Color {
	out := convert.FromOKLCHMapped(o)
	out.A = alpha
	return out
}

// LightenOKLCH adds amount to OKLCH lightness (L is in [0, 1]).
func LightenOKLCH(c golor.Color, amount float64) golor.Color {
	o := convert.ToOKLCH(c)
	o.L = space.Clamp01(o.L + amount)
	return fromOKLCH(o, c.A)
}

// DarkenOKLCH subtracts amount from OKLCH lightness (L is in [0, 1]).
func DarkenOKLCH(c golor.Color, amount float64) golor.Color {
	return LightenOKLCH(c, -amount)
}

// SaturateOKLCH adds amount × 0.4 to chroma; greys are unchanged.
func SaturateOKLCH(c golor.Color, amount float64) golor.Color {
	o := convert.ToOKLCH(c)
	if o.IsAchromatic() {
		return c
	}
	o.C = math.Max(0, o.C+amount*maxOKLCHChroma)
	return fromOKLCH(o, c.A)
}

// DesaturateOKLCH decreases OKLCH chroma by amount × 0.4 (amount 1 always yields grey).
func DesaturateOKLCH(c golor.Color, amount float64) golor.Color {
	o := convert.ToOKLCH(c)
	if o.IsAchromatic() {
		return c
	}
	o.C = math.Max(0, o.C-amount*maxOKLCHChroma)
	return fromOKLCH(o, c.A)
}

// ShiftHueOKLCH rotates hue at constant lightness/chroma; greys are unchanged.
func ShiftHueOKLCH(c golor.Color, degrees float64) golor.Color {
	o := convert.ToOKLCH(c)
	if o.IsAchromatic() {
		return c
	}
	o.H += degrees
	return fromOKLCH(o, c.A)
}
