package convert

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// CMYK holds normalized channels in [0,1]; it is an unprofiled conversion without alpha.
type CMYK struct{ C, M, Y, K float64 }

// ToCMYK converts an sRGB Color to CMYK. Alpha is not represented in CMYK.
func ToCMYK(c golor.Color) CMYK {
	r, g, b := channels(c)
	k := 1 - max(r, g, b)
	if k >= 1 {
		return CMYK{K: 1}
	}
	return CMYK{
		C: (1 - r - k) / (1 - k),
		M: (1 - g - k) / (1 - k),
		Y: (1 - b - k) / (1 - k),
		K: k,
	}
}

// FromCMYK converts CMYK to an opaque sRGB Color.
func FromCMYK(cmyk CMYK) golor.Color {
	c := space.Clamp01(cmyk.C)
	m := space.Clamp01(cmyk.M)
	y := space.Clamp01(cmyk.Y)
	k := space.Clamp01(cmyk.K)
	return golor.RGBf((1-c)*(1-k), (1-m)*(1-k), (1-y)*(1-k))
}
