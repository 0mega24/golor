package blend

import "github.com/0mega24/golor/v2"

// Screen applies the Screen blend mode with Porter-Duff source-over alpha compositing.
func Screen(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		1-(1-base.R)*(1-layer.R),
		1-(1-base.G)*(1-layer.G),
		1-(1-base.B)*(1-layer.B),
	))
}
