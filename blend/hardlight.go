package blend

import "github.com/0mega24/golor/v2"

// HardLight applies the Hard Light blend mode with Porter-Duff source-over alpha compositing.
// Uses the layer channel to determine which branch applies (Overlay with base/layer swapped).
func HardLight(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		overlayChannel(layer.R, base.R),
		overlayChannel(layer.G, base.G),
		overlayChannel(layer.B, base.B),
	))
}
