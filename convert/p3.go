package convert

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// P3 represents gamma-encoded Display P3 channels in [0,1], with a wider gamut than sRGB.
type P3 struct{ R, G, B float64 }

// ToP3 converts an sRGB Color to Display P3. Alpha is not represented.
func ToP3(c golor.Color) P3 {
	v := space.DelinearizeVec(space.XYZToP3.Mul(space.SRGBToXYZ.Mul(linearRGB(c))))
	return P3{v[0], v[1], v[2]}
}

func (p P3) linearSRGB() space.Vec3 {
	return space.XYZToSRGB.Mul(space.P3ToXYZ.Mul(space.LinearizeVec(space.Vec3{p.R, p.G, p.B})))
}

// FromP3 converts Display P3 to an opaque sRGB Color, clipping out-of-gamut values per
// channel. Use FromP3Mapped to preserve hue and lightness instead.
func FromP3(p P3) golor.Color { return fromLinearClipped(p.linearSRGB()) }

// FromP3Mapped converts Display P3 to an opaque sRGB Color using CSS Color 4 gamut mapping.
func FromP3Mapped(p P3) golor.Color { return fromLinearMapped(p.linearSRGB()) }

// InSRGBGamut reports whether p is representable in sRGB without clipping.
func (p P3) InSRGBGamut() bool { return space.InGamutLinear(p.linearSRGB()) }
