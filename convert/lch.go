package convert

import (
	"math"

	"github.com/0mega24/golor"
)

// LCH represents a color in CIE L*C*h° space (D65 illuminant).
type LCH struct{ L, C, H float64 }

// ToLCH converts an sRGB Color to CIE L*C*h° (D65 illuminant).
func ToLCH(c golor.Color) LCH {
	lab := ToLAB(c)
	h := math.Atan2(lab.B, lab.A) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	return LCH{
		L: lab.L,
		C: math.Sqrt(lab.A*lab.A + lab.B*lab.B),
		H: h,
	}
}

// FromLCH converts CIE L*C*h° to an sRGB Color (D65 illuminant).
// Out-of-gamut values are clamped.
func FromLCH(lch LCH) golor.Color {
	hRad := lch.H * math.Pi / 180
	lab := LAB{
		L: lch.L,
		A: lch.C * math.Cos(hRad),
		B: lch.C * math.Sin(hRad),
	}
	return FromLAB(lab)
}
