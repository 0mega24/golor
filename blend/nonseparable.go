package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// Non-separable modes use W3C Compositing Level 1 with Rec. 601 luma.

type rgb [3]float64

func toRGB(c golor.Color) rgb { return rgb{space.Clamp01(c.R), space.Clamp01(c.G), space.Clamp01(c.B)} }

func lum(c rgb) float64 { return 0.3*c[0] + 0.59*c[1] + 0.11*c[2] }

func clipColor(c rgb) rgb {
	l := lum(c)
	n := min(c[0], c[1], c[2])
	x := max(c[0], c[1], c[2])
	if n < 0 && l-n > 0 {
		k := l / (l - n)
		c = rgb{l + (c[0]-l)*k, l + (c[1]-l)*k, l + (c[2]-l)*k}
	}
	if x > 1 && x-l > 0 {
		k := (1 - l) / (x - l)
		c = rgb{l + (c[0]-l)*k, l + (c[1]-l)*k, l + (c[2]-l)*k}
	}
	return c
}

func setLum(c rgb, l float64) rgb {
	d := l - lum(c)
	return clipColor(rgb{c[0] + d, c[1] + d, c[2] + d})
}

func sat(c rgb) float64 { return max(c[0], c[1], c[2]) - min(c[0], c[1], c[2]) }

// setSat scales c to saturation s while preserving channel order.
func setSat(c rgb, s float64) rgb {
	maxI, midI, minI := 0, 1, 2
	if c[maxI] < c[midI] {
		maxI, midI = midI, maxI
	}
	if c[midI] < c[minI] {
		midI, minI = minI, midI
	}
	if c[maxI] < c[midI] {
		maxI, midI = midI, maxI
	}
	var out rgb
	if c[maxI] > c[minI] {
		out[midI] = (c[midI] - c[minI]) * s / (c[maxI] - c[minI])
		out[maxI] = s
	}
	return out
}

func fromRGB(c rgb) golor.Color { return golor.RGBf(c[0], c[1], c[2]) }

// Hue applies the Hue blend mode with Porter-Duff source-over alpha compositing.
// The result has the hue of layer and the saturation and luminosity of base.
func Hue(base, layer golor.Color) golor.Color {
	cb, cs := toRGB(base), toRGB(layer)
	return compositeBlend(base, layer, fromRGB(setLum(setSat(cs, sat(cb)), lum(cb))))
}

// Saturation applies the Saturation blend mode with Porter-Duff source-over alpha compositing.
// The result has the saturation of layer and the hue and luminosity of base.
func Saturation(base, layer golor.Color) golor.Color {
	cb, cs := toRGB(base), toRGB(layer)
	return compositeBlend(base, layer, fromRGB(setLum(setSat(cb, sat(cs)), lum(cb))))
}

// Color applies the layer's hue/saturation and base's luminosity, with source-over alpha.
func Color(base, layer golor.Color) golor.Color {
	cb, cs := toRGB(base), toRGB(layer)
	return compositeBlend(base, layer, fromRGB(setLum(cs, lum(cb))))
}

// Luminosity applies the Luminosity blend mode with Porter-Duff source-over alpha compositing.
// The result has the luminosity of layer and the hue and saturation of base.
func Luminosity(base, layer golor.Color) golor.Color {
	cb, cs := toRGB(base), toRGB(layer)
	return compositeBlend(base, layer, fromRGB(setLum(cb, lum(cs))))
}
