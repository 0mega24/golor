package convert

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// XYZ represents a color in CIE 1931 XYZ space relative to D65, with Y = 1 for white.
type XYZ struct{ X, Y, Z float64 }

// ToXYZ converts an sRGB Color to CIE XYZ (D65).
func ToXYZ(c golor.Color) XYZ {
	v := space.SRGBToXYZ.Mul(linearRGB(c))
	return XYZ{v[0], v[1], v[2]}
}

// FromXYZ converts CIE XYZ (D65) to an sRGB Color, clipping out-of-gamut values per channel.
func FromXYZ(xyz XYZ) golor.Color {
	return fromLinearClipped(space.XYZToSRGB.Mul(space.Vec3{xyz.X, xyz.Y, xyz.Z}))
}
