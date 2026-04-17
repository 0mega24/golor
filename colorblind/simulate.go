// Package colorblind provides color vision deficiency simulation and accessible palette tools.
package colorblind

import (
	"math"

	"github.com/0mega24/golor"
)

// Deficiency represents a type of color vision deficiency.
type Deficiency int

// Deficiency constants enumerate the supported types of color vision deficiency.
const (
	Deuteranopia Deficiency = iota // green-weak
	Protanopia                     // red-weak
	Tritanopia                     // blue-yellow
)

// Simulate returns the approximate appearance of c for someone with the given deficiency.
// Uses the Viénot et al. (1999) matrix method for deuteranopia/protanopia,
// and Brettel et al. (1997) for tritanopia.
func Simulate(c golor.Color, d Deficiency) golor.Color {
	r := cbLinearize(c.R)
	g := cbLinearize(c.G)
	b := cbLinearize(c.B)

	var sr, sg, sb float64
	switch d {
	case Deuteranopia:
		sr = 0.29901*r + 0.58431*g + 0.11668*b
		sg = 0.29901*r + 0.58431*g + 0.11668*b
		sb = b
	case Protanopia:
		sr = 0.10889*r + 0.89111*g
		sg = 0.10889*r + 0.89111*g
		sb = b
	case Tritanopia:
		sr = r
		sg = 0.73407*g + 0.26593*b
		sb = 0.73407*g + 0.26593*b
	}

	return golor.RGBf(cbGamma(sr), cbGamma(sg), cbGamma(sb))
}

func cbLinearize(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func cbGamma(v float64) float64 {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	if v <= 0.0031308 {
		return 12.92 * v
	}
	return 1.055*math.Pow(v, 1/2.4) - 0.055
}
