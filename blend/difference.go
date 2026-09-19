package blend

import (
	"math"

	"github.com/0mega24/golor/v2"
)

// Difference applies the Difference blend mode with Porter-Duff source-over alpha compositing.
func Difference(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		math.Abs(base.R-layer.R),
		math.Abs(base.G-layer.G),
		math.Abs(base.B-layer.B),
	))
}
