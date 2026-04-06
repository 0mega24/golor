package blend

import "github.com/0mega24/golor"

// HardLight applies the Photoshop Hard Light blend mode (channel-wise).
// Uses the layer channel to determine which branch applies (Overlay with base/layer swapped).
func HardLight(base, layer golor.Color) golor.Color {
	return golor.RGBf(
		overlayChannel(layer.R, base.R),
		overlayChannel(layer.G, base.G),
		overlayChannel(layer.B, base.B),
	)
}
