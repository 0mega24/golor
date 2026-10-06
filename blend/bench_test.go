package blend_test

import (
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/blend"
)

var benchColors = func() []golor.Color {
	r := rand.New(rand.NewSource(42))
	out := make([]golor.Color, 1024)
	for i := range out {
		out[i] = golor.RGB(uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(r.Intn(256)))
	}
	return out
}()

var sink golor.Color

func BenchmarkBlend(b *testing.B) {
	modes := []struct {
		name string
		fn   func(a, b golor.Color) golor.Color
	}{
		{"Multiply", blend.Multiply},
		{"SoftLight", blend.SoftLight},
		{"Hue", blend.Hue},
		{"Luminosity", blend.Luminosity},
		{"Mix", func(a, c golor.Color) golor.Color { return blend.Mix(a, c, 0.3) }},
	}
	for _, m := range modes {
		b.Run(m.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sink = m.fn(benchColors[i&1023], benchColors[(i+7)&1023])
			}
		})
	}
}
