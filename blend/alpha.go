package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// compositeBlend applies W3C backdrop-alpha blending followed by source-over compositing.
func compositeBlend(base, layer, blended golor.Color) golor.Color {
	ab, as := space.Clamp01(base.A), space.Clamp01(layer.A)
	switch {
	case as == 0:
		return golor.RGBAf(base.R, base.G, base.B, ab)
	case ab == 1 && as == 1:
		return golor.RGBf(blended.R, blended.G, blended.B)
	}
	ao := as + ab*(1-as)
	return golor.RGBAf(
		compositeChannel(space.Clamp01(base.R), space.Clamp01(layer.R), space.Clamp01(blended.R), ab, as, ao),
		compositeChannel(space.Clamp01(base.G), space.Clamp01(layer.G), space.Clamp01(blended.G), ab, as, ao),
		compositeChannel(space.Clamp01(base.B), space.Clamp01(layer.B), space.Clamp01(blended.B), ab, as, ao),
		ao,
	)
}

func compositeChannel(cb, cs, b, ab, as, ao float64) float64 {
	mixed := (1-ab)*cs + ab*b
	return (as*mixed + ab*cb*(1-as)) / ao
}
