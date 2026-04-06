package convert

import (
	"math"

	"github.com/0mega24/golor"
)

// LAB represents a color in CIE L*a*b* space (D65 illuminant).
type LAB struct{ L, A, B float64 }

// D65 white point for XYZ conversions.
const (
	d65X = 0.95047
	d65Y = 1.00000
	d65Z = 1.08883
)

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func linearize(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func delinearize(v float64) float64 {
	if v <= 0.0031308 {
		return 12.92 * v
	}
	return 1.055*math.Pow(v, 1/2.4) - 0.055
}

func rgbToXYZ(c golor.Color) (x, y, z float64) {
	r := linearize(c.R)
	g := linearize(c.G)
	b := linearize(c.B)
	x = 0.4124564*r + 0.3575761*g + 0.1804375*b
	y = 0.2126729*r + 0.7151522*g + 0.0721750*b
	z = 0.0193339*r + 0.1191920*g + 0.9503041*b
	return
}

func xyzToRGB(x, y, z float64) golor.Color {
	r := 3.2404542*x - 1.5371385*y - 0.4985314*z
	g := -0.9692660*x + 1.8760108*y + 0.0415560*z
	b := 0.0556434*x - 0.2040259*y + 1.0572252*z
	return golor.RGBf(delinearize(clamp01(r)), delinearize(clamp01(g)), delinearize(clamp01(b)))
}

func labF(t float64) float64 {
	if t > 0.008856 {
		return math.Cbrt(t)
	}
	return t/0.128418 + 4.0/29
}

func labFInv(t float64) float64 {
	if t > 6.0/29 {
		return t * t * t
	}
	return 0.128418 * (t - 4.0/29)
}

// ToLAB converts an sRGB Color to CIE L*a*b* (D65 illuminant).
func ToLAB(c golor.Color) LAB {
	x, y, z := rgbToXYZ(c)
	fx := labF(x / d65X)
	fy := labF(y / d65Y)
	fz := labF(z / d65Z)
	return LAB{
		L: 116*fy - 16,
		A: 500 * (fx - fy),
		B: 200 * (fy - fz),
	}
}

// FromLAB converts CIE L*a*b* to an sRGB Color (D65 illuminant).
// Out-of-gamut values are clamped.
func FromLAB(lab LAB) golor.Color {
	fy := (lab.L + 16) / 116
	fx := lab.A/500 + fy
	fz := fy - lab.B/200
	x := labFInv(fx) * d65X
	y := labFInv(fy) * d65Y
	z := labFInv(fz) * d65Z
	return xyzToRGB(x, y, z)
}
