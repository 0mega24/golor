package convert

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// LAB represents a color in CIE L*a*b* space (D65 illuminant).
type LAB struct{ L, A, B float64 }

// channels returns c's RGB channels clamped to [0, 1], with NaN mapped to 0, so that
// struct literals with out-of-range fields cannot produce out-of-range results.
func channels(c golor.Color) (r, g, b float64) {
	return space.Clamp01(c.R), space.Clamp01(c.G), space.Clamp01(c.B)
}

func linearRGB(c golor.Color) space.Vec3 {
	r, g, b := channels(c)
	return space.LinearizeVec(space.Vec3{r, g, b})
}

func fromLinearClipped(rgb space.Vec3) golor.Color {
	return golor.RGBf(
		space.Delinearize(space.Clamp01(rgb[0])),
		space.Delinearize(space.Clamp01(rgb[1])),
		space.Delinearize(space.Clamp01(rgb[2])),
	)
}

func fromLinearMapped(rgb space.Vec3) golor.Color {
	m := space.GamutMapLinear(rgb)
	return golor.RGBf(space.Delinearize(m[0]), space.Delinearize(m[1]), space.Delinearize(m[2]))
}

func labToLinear(lab LAB) space.Vec3 {
	return space.XYZToSRGB.Mul(space.LabToXYZ(space.Vec3{lab.L, lab.A, lab.B}, space.D65))
}

// ToLAB converts an sRGB Color to CIE L*a*b* (D65 illuminant).
func ToLAB(c golor.Color) LAB {
	v := space.XYZToLab(space.SRGBToXYZ.Mul(linearRGB(c)), space.D65)
	return LAB{v[0], v[1], v[2]}
}

// FromLAB converts D65 LAB to sRGB with clipping; FromLABMapped preserves hue/lightness.
func FromLAB(lab LAB) golor.Color {
	return fromLinearClipped(labToLinear(lab))
}

// FromLABMapped converts CIE L*a*b* to an sRGB Color, bringing out-of-gamut values into
// gamut with the CSS Color 4 gamut-mapping algorithm (OKLCH chroma reduction).
func FromLABMapped(lab LAB) golor.Color {
	return fromLinearMapped(labToLinear(lab))
}

// InGamut reports whether lab is representable in sRGB without clipping.
func (lab LAB) InGamut() bool {
	return space.InGamutLinear(labToLinear(lab))
}
