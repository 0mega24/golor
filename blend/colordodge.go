package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// ColorDodge applies the Color Dodge blend mode with Porter-Duff source-over alpha compositing.
func ColorDodge(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		colorDodgeChannel(base.R, layer.R),
		colorDodgeChannel(base.G, layer.G),
		colorDodgeChannel(base.B, layer.B),
	))
}

func colorDodgeChannel(base, layer float64) float64 {
	if base <= 0 {
		return 0
	}
	if layer >= 1 {
		return 1
	}
	return space.Clamp01(base / (1 - layer))
}
