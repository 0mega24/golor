package blend

import (
	"math"

	"github.com/0mega24/golor/v2"
)

// SoftLight applies the Soft Light blend mode with Porter-Duff source-over alpha compositing.
func SoftLight(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		softLightChannel(base.R, layer.R),
		softLightChannel(base.G, layer.G),
		softLightChannel(base.B, layer.B),
	))
}

func softLightChannel(base, layer float64) float64 {
	if layer <= 0.5 {
		return base - (1-2*layer)*base*(1-base)
	}
	var d float64
	if base <= 0.25 {
		d = ((16*base-12)*base + 4) * base
	} else {
		d = math.Sqrt(base)
	}
	return base + (2*layer-1)*(d-base)
}
