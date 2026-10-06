package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// LinearLight applies the Linear Light blend mode with Porter-Duff source-over alpha compositing.
func LinearLight(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		linearLightChannel(base.R, layer.R),
		linearLightChannel(base.G, layer.G),
		linearLightChannel(base.B, layer.B),
	))
}

func linearLightChannel(base, layer float64) float64 {
	return space.Clamp01(base + 2*layer - 1)
}
