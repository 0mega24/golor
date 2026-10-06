package convert

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// OKLab represents perceptual lightness L in [0,1] and A/B roughly in ±0.4 for sRGB.
type OKLab struct{ L, A, B float64 }

// OKLCH is the polar form of OKLab: lightness, chroma, and hue in degrees [0, 360).
// For achromatic colors (C below 1e-6), H is 0 and has no effect; see IsAchromatic.
type OKLCH struct{ L, C, H float64 }

// IsAchromatic reports whether the color is grey, in which case its hue has no effect.
func (o OKLCH) IsAchromatic() bool { return o.C < space.AchromaticOKLCH }

// ToOKLab converts an sRGB Color to OKLab.
func ToOKLab(c golor.Color) OKLab {
	v := space.LinearSRGBToOKLab(linearRGB(c))
	return OKLab{v[0], v[1], v[2]}
}

func (o OKLab) linear() space.Vec3 { return space.OKLabToLinearSRGB(space.Vec3{o.L, o.A, o.B}) }

// FromOKLab converts OKLab to an sRGB Color, clipping out-of-gamut values per channel.
// Use FromOKLabMapped to preserve hue and lightness instead.
func FromOKLab(o OKLab) golor.Color { return fromLinearClipped(o.linear()) }

// FromOKLabMapped converts OKLab to an sRGB Color using CSS Color 4 gamut mapping.
func FromOKLabMapped(o OKLab) golor.Color { return fromLinearMapped(o.linear()) }

// InGamut reports whether o is representable in sRGB without clipping.
func (o OKLab) InGamut() bool { return space.InGamutLinear(o.linear()) }

// ToOKLCH converts an sRGB Color to OKLCH.
func ToOKLCH(c golor.Color) OKLCH {
	lab := ToOKLab(c)
	v := space.ToPolar(space.Vec3{lab.L, lab.A, lab.B}, space.AchromaticOKLCH)
	return OKLCH{v[0], v[1], v[2]}
}

func (o OKLCH) lab() OKLab {
	v := space.FromPolar(space.Vec3{o.L, o.C, o.H})
	return OKLab{v[0], v[1], v[2]}
}

// FromOKLCH converts OKLCH to an sRGB Color, clipping out-of-gamut values per channel.
// Use FromOKLCHMapped to preserve hue and lightness instead.
func FromOKLCH(o OKLCH) golor.Color { return FromOKLab(o.lab()) }

// FromOKLCHMapped converts OKLCH to an sRGB Color using CSS Color 4 gamut mapping
// (chroma reduction at constant lightness and hue).
func FromOKLCHMapped(o OKLCH) golor.Color {
	m := space.GamutMapOKLCH(space.Vec3{o.L, o.C, o.H})
	return golor.RGBf(space.Delinearize(m[0]), space.Delinearize(m[1]), space.Delinearize(m[2]))
}

// InGamut reports whether o is representable in sRGB without clipping.
func (o OKLCH) InGamut() bool { return o.lab().InGamut() }

// DeltaEOK returns the Euclidean distance between two colors in OKLab. A difference of
// about 0.02 is just noticeable.
func DeltaEOK(a, b golor.Color) float64 {
	x, y := ToOKLab(a), ToOKLab(b)
	return space.Vec3Distance(space.Vec3{x.L, x.A, x.B}, space.Vec3{y.L, y.A, y.B})
}
