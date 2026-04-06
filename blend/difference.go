package blend

import (
	"math"

	"github.com/0mega24/golor"
)

// Difference applies the Photoshop Difference blend mode (channel-wise).
func Difference(base, layer golor.Color) golor.Color {
	return golor.RGBf(
		math.Abs(base.R-layer.R),
		math.Abs(base.G-layer.G),
		math.Abs(base.B-layer.B),
	)
}
