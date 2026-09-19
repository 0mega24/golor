package convert

import "github.com/0mega24/golor/v2"

// CMYK represents a color in cyan-magenta-yellow-key space.
// C, M, Y, and K are normalized to [0, 1].
type CMYK struct{ C, M, Y, K float64 }

// ToCMYK converts an sRGB Color to CMYK. Alpha is not represented in CMYK.
func ToCMYK(c golor.Color) CMYK {
	k := 1 - max3(c.R, c.G, c.B)
	if k >= 1 {
		return CMYK{K: 1}
	}
	return CMYK{
		C: (1 - c.R - k) / (1 - k),
		M: (1 - c.G - k) / (1 - k),
		Y: (1 - c.B - k) / (1 - k),
		K: k,
	}
}

// FromCMYK converts CMYK to an opaque sRGB Color.
func FromCMYK(cmyk CMYK) golor.Color {
	c := clamp01(cmyk.C)
	m := clamp01(cmyk.M)
	y := clamp01(cmyk.Y)
	k := clamp01(cmyk.K)
	return golor.RGBf((1-c)*(1-k), (1-m)*(1-k), (1-y)*(1-k))
}

func max3(a, b, c float64) float64 {
	if b > a {
		a = b
	}
	if c > a {
		a = c
	}
	return a
}
