package blend

import "github.com/0mega24/golor/v2"

func compositeBlend(base, layer, blended golor.Color) golor.Color {
	return golor.RGBAf(
		base.R+layer.A*(blended.R-base.R),
		base.G+layer.A*(blended.G-base.G),
		base.B+layer.A*(blended.B-base.B),
		blendedAlpha(base, layer),
	)
}

func blendedAlpha(base, layer golor.Color) float64 {
	return layer.A + base.A*(1-layer.A)
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
