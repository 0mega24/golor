// Package convert provides bidirectional conversions between sRGB and HSL, HSV, L*a*b*, and LCH color spaces.
package convert

import (
	"math"

	"github.com/0mega24/golor"
)

// HSL represents a color in hue-saturation-lightness space.
// H is in [0, 360), S and L are in [0, 1].
type HSL struct{ H, S, L float64 }

// ToHSL converts an sRGB Color to HSL.
func ToHSL(c golor.Color) HSL {
	r, g, b := c.R, c.G, c.B
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
	h *= 60
	return HSL{h, s, l}
}

// FromHSL converts HSL to an sRGB Color.
func FromHSL(hsl HSL) golor.Color {
	h, s, l := hsl.H, hsl.S, hsl.L
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
	r := hueToRGB(p, q, h/360+1.0/3)
	g := hueToRGB(p, q, h/360)
	b := hueToRGB(p, q, h/360-1.0/3)
	return golor.RGBf(r, g, b)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
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
