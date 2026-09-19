package blend

import "github.com/0mega24/golor/v2"

func blendedAlpha(base, layer golor.Color) float64 {
	return layer.A + base.A*(1-layer.A)
}
