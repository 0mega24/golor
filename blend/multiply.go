package blend

import "github.com/0mega24/golor/v2"

// Multiply applies the Photoshop Multiply blend mode (channel-wise).
// The result alpha is the source-over composited alpha of base and layer.
func Multiply(base, layer golor.Color) golor.Color {
	return golor.RGBAf(
		base.R*layer.R,
		base.G*layer.G,
		base.B*layer.B,
		blendedAlpha(base, layer),
	)
}
