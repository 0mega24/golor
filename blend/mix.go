// Package blend implements W3C blend modes; base is the backdrop, layer is the source.
package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// Mix interpolates premultiplied color/alpha; t is clamped to [0,1], with NaN treated as 0.
func Mix(a, b golor.Color, t float64) golor.Color {
	t = space.Clamp01(t)
	if t == 0 {
		return a
	}
	if t == 1 {
		return b
	}
	aa, ba := space.Clamp01(a.A), space.Clamp01(b.A)
	alpha := aa + t*(ba-aa)
	if alpha <= 0 {
		return golor.RGBAf(0, 0, 0, 0)
	}
	ch := func(x, y float64) float64 {
		return (space.Clamp01(x)*aa + t*(space.Clamp01(y)*ba-space.Clamp01(x)*aa)) / alpha
	}
	return golor.RGBAf(ch(a.R, b.R), ch(a.G, b.G), ch(a.B, b.B), alpha)
}
