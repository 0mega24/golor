package convert

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// LCH holds D65 lightness/chroma and hue in [0,360); C < 1e-4 makes hue zero.
type LCH struct{ L, C, H float64 }

// IsAchromatic reports whether the color is grey, in which case its hue has no effect
// and interpolation should borrow the other endpoint's hue.
func (lch LCH) IsAchromatic() bool { return lch.C < space.AchromaticLCH }

// ToLCH converts an sRGB Color to CIE L*C*h° (D65 illuminant).
func ToLCH(c golor.Color) LCH {
	lab := ToLAB(c)
	v := space.ToPolar(space.Vec3{lab.L, lab.A, lab.B}, space.AchromaticLCH)
	return LCH{v[0], v[1], v[2]}
}

func (lch LCH) lab() LAB {
	v := space.FromPolar(space.Vec3{lch.L, lch.C, lch.H})
	return LAB{v[0], v[1], v[2]}
}

// FromLCH converts D65 LCH to sRGB with clipping; FromLCHMapped preserves hue/lightness.
func FromLCH(lch LCH) golor.Color {
	return FromLAB(lch.lab())
}

// FromLCHMapped converts CIE L*C*h° to an sRGB Color, bringing out-of-gamut values into
// gamut with the CSS Color 4 gamut-mapping algorithm (OKLCH chroma reduction).
func FromLCHMapped(lch LCH) golor.Color {
	return FromLABMapped(lch.lab())
}

// InGamut reports whether lch is representable in sRGB without clipping.
func (lch LCH) InGamut() bool { return lch.lab().InGamut() }
