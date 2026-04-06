package blend

import "github.com/0mega24/golor"

// Multiply applies the Photoshop Multiply blend mode (channel-wise).
func Multiply(base, layer golor.Color) golor.Color {
	return golor.RGBf(
		base.R*layer.R,
		base.G*layer.G,
		base.B*layer.B,
	)
}
