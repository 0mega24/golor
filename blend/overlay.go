package blend

import "github.com/0mega24/golor/v2"

// Overlay applies the Photoshop Overlay blend mode (channel-wise).
// Uses the base channel to determine which branch applies.
// The result alpha is the source-over composited alpha of base and layer.
func Overlay(base, layer golor.Color) golor.Color {
	return golor.RGBAf(
		overlayChannel(base.R, layer.R),
		overlayChannel(base.G, layer.G),
		overlayChannel(base.B, layer.B),
		blendedAlpha(base, layer),
	)
}

func overlayChannel(base, layer float64) float64 {
	if base <= 0.5 {
		return 2 * base * layer
	}
	return 1 - 2*(1-base)*(1-layer)
}
