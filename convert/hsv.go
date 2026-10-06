package convert

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// HSV represents a color in hue-saturation-value space.
// H is in [0, 360); S and V are in [0, 1]. Greys have S = 0 and H = 0.
type HSV struct{ H, S, V float64 }

// IsAchromatic reports whether the color is grey, in which case its hue has no effect.
func (hsv HSV) IsAchromatic() bool { return hsv.S < 1e-12 }

// ToHSV converts an sRGB Color to HSV.
func ToHSV(c golor.Color) HSV {
	r, g, b := channels(c)
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min
	if delta == 0 {
		return HSV{0, 0, max}
	}
	return HSV{hue(r, g, b, max, delta), delta / max, max}
}

// FromHSV converts HSV to an sRGB Color. H may be any angle (it is wrapped into
// [0, 360)); S and V are clamped to [0, 1].
func FromHSV(hsv HSV) golor.Color {
	h, s, v := space.NormalizeHue(hsv.H)/60, space.Clamp01(hsv.S), space.Clamp01(hsv.V)
	if s == 0 {
		return golor.RGBf(v, v, v)
	}
	i := math.Floor(h)
	f := h - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))
	switch int(i) {
	case 0:
		return golor.RGBf(v, t, p)
	case 1:
		return golor.RGBf(q, v, p)
	case 2:
		return golor.RGBf(p, v, t)
	case 3:
		return golor.RGBf(p, q, v)
	case 4:
		return golor.RGBf(t, p, v)
	default:
		return golor.RGBf(v, p, q)
	}
}
