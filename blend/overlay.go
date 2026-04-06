package blend

import "github.com/0mega24/golor"

// Overlay applies the Photoshop Overlay blend mode (channel-wise).
// Uses the base channel to determine which branch applies.
func Overlay(base, layer golor.Color) golor.Color {
	return golor.RGBf(
		overlayChannel(base.R, layer.R),
		overlayChannel(base.G, layer.G),
		overlayChannel(base.B, layer.B),
	)
}

func overlayChannel(base, layer float64) float64 {
	if base <= 0.5 {
		return 2 * base * layer
	}
	return 1 - 2*(1-base)*(1-layer)
}
