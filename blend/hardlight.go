package blend

import "github.com/0mega24/golor/v2"

// HardLight applies the Photoshop Hard Light blend mode (channel-wise).
// Uses the layer channel to determine which branch applies (Overlay with base/layer swapped).
// The result alpha is the source-over composited alpha of base and layer.
func HardLight(base, layer golor.Color) golor.Color {
	return golor.RGBAf(
		overlayChannel(layer.R, base.R),
		overlayChannel(layer.G, base.G),
		overlayChannel(layer.B, base.B),
		blendedAlpha(base, layer),
	)
}
