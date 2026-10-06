// Package convert provides sRGB color-space conversions, gamut checks, and gamut mapping.
package convert

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// HSL represents a color in hue-saturation-lightness space.
// H is in [0, 360); S and L are in [0, 1]. Greys have S = 0 and H = 0.
type HSL struct{ H, S, L float64 }

// IsAchromatic reports whether the color is grey, in which case its hue has no effect.
func (hsl HSL) IsAchromatic() bool { return hsl.S < 1e-12 }

// ToHSL converts an sRGB Color to HSL.
func ToHSL(c golor.Color) HSL {
	r, g, b := channels(c)
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l := (max + min) / 2
	delta := max - min
	if delta == 0 {
		return HSL{0, 0, l}
	}
	var s float64
	if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2 - max - min)
	}
	return HSL{hue(r, g, b, max, delta), s, l}
}

// hue returns the hue in degrees shared by HSL and HSV.
func hue(r, g, b, max, delta float64) float64 {
	var h float64
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
	return h * 60
}

// FromHSL converts HSL to an sRGB Color. H may be any angle (it is wrapped into
// [0, 360)); S and L are clamped to [0, 1].
func FromHSL(hsl HSL) golor.Color {
	h, s, l := space.NormalizeHue(hsl.H)/360, space.Clamp01(hsl.S), space.Clamp01(hsl.L)
	if s == 0 {
		return golor.RGBf(l, l, l)
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	return golor.RGBf(hueToRGB(p, q, h+1.0/3), hueToRGB(p, q, h), hueToRGB(p, q, h-1.0/3))
}

// hueToRGB expects t in (-1/3, 4/3), which holds for normalized hues.
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t >= 1 {
		t--
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 1.0/2:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	default:
		return p
	}
}
