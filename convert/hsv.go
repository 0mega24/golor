package convert

import (
	"math"

	"github.com/0mega24/golor"
)

// HSV represents a color in hue-saturation-value space.
// H is in [0, 360), S and V are in [0, 1].
type HSV struct{ H, S, V float64 }

// ToHSV converts an sRGB Color to HSV.
func ToHSV(c golor.Color) HSV {
	r, g, b := c.R, c.G, c.B
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	v := max
	var s, h float64
	if max != 0 {
		s = delta / max
	}
	if delta == 0 {
		return HSV{0, s, v}
	}
	switch max {
	case r:
		h = (g - b) / delta
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/delta + 2
	default:
		h = (r-g)/delta + 4
	}
	h *= 60
	return HSV{h, s, v}
}

// FromHSV converts HSV to an sRGB Color.
func FromHSV(hsv HSV) golor.Color {
	h, s, v := hsv.H, hsv.S, hsv.V
	if s == 0 {
		return golor.RGBf(v, v, v)
	}
	h = math.Mod(h, 360) / 60
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
