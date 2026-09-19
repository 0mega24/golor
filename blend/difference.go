package blend

import (
	"math"

	"github.com/0mega24/golor/v2"
)

// Difference applies the Photoshop Difference blend mode (channel-wise).
// The result alpha is the source-over composited alpha of base and layer.
func Difference(base, layer golor.Color) golor.Color {
	return golor.RGBAf(
		math.Abs(base.R-layer.R),
		math.Abs(base.G-layer.G),
		math.Abs(base.B-layer.B),
		blendedAlpha(base, layer),
	)
}
