package blend

import "github.com/0mega24/golor/v2"

// Screen applies the Photoshop Screen blend mode (channel-wise).
// The result alpha is the source-over composited alpha of base and layer.
func Screen(base, layer golor.Color) golor.Color {
	return golor.RGBAf(
		1-(1-base.R)*(1-layer.R),
		1-(1-base.G)*(1-layer.G),
		1-(1-base.B)*(1-layer.B),
		blendedAlpha(base, layer),
	)
}
