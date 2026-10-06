// Package colorblind provides color vision deficiency simulation and accessible palette tools.
package colorblind

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// Deficiency represents a type of color vision deficiency.
type Deficiency int

// Deficiency constants enumerate the supported types of color vision deficiency.
const (
	Deuteranopia Deficiency = iota // Green-blind.
	Protanopia                     // Red-blind.
	Tritanopia                     // Blue-blind.
)

// Parameters: libDaltonLens (public domain), Smith & Pokorny LMS.
// https://github.com/DaltonLens/libDaltonLens/blob/master/libDaltonLens.c.
var (
	// vienotProtan and vienotDeutan use Viénot 1999, a single linear RGB projection accurate for protan/deutan.
	vienotProtan = space.Mat3{
		{0.11238, 0.88762, 0.00000},
		{0.11238, 0.88762, -0.00000},
		{0.00401, -0.00401, 1.00000},
	}
	vienotDeutan = space.Mat3{
		{0.29275, 0.70725, 0.00000},
		{0.29275, 0.70725, -0.00000},
		{-0.02234, 0.02234, 1.00000},
	}
	// Brettel 1997 tritan projections, selected by the separation plane.
	brettelTritan1 = space.Mat3{
		{1.01277, 0.13548, -0.14826},
		{-0.01243, 0.86812, 0.14431},
		{0.07589, 0.80500, 0.11911},
	}
	brettelTritan2 = space.Mat3{
		{0.93678, 0.18979, -0.12657},
		{0.06154, 0.81526, 0.12320},
		{-0.37562, 1.12767, 0.24796},
	}
	brettelTritanNormal = space.Vec3{0.03901, -0.02788, -0.01113}
)

// Simulate uses Viénot 1999 for protan/deutan and Brettel 1997 for tritan.
// Alpha is preserved; unknown deficiencies leave c unchanged.
func Simulate(c golor.Color, d Deficiency) golor.Color {
	rgb := space.LinearizeVec(space.Vec3{space.Clamp01(c.R), space.Clamp01(c.G), space.Clamp01(c.B)})
	var out space.Vec3
	switch d {
	case Deuteranopia:
		out = vienotDeutan.Mul(rgb)
	case Protanopia:
		out = vienotProtan.Mul(rgb)
	case Tritanopia:
		n := brettelTritanNormal
		if rgb[0]*n[0]+rgb[1]*n[1]+rgb[2]*n[2] >= 0 {
			out = brettelTritan1.Mul(rgb)
		} else {
			out = brettelTritan2.Mul(rgb)
		}
	default:
		return c
	}
	return golor.RGBAf(
		space.Delinearize(space.Clamp01(out[0])),
		space.Delinearize(space.Clamp01(out[1])),
		space.Delinearize(space.Clamp01(out[2])),
		c.A,
	)
}
