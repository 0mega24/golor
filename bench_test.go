package golor_test

import (
	"image"
	"image/color"
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
)

var (
	benchColor  golor.Color
	benchString string
	benchPixels []golor.Color
)

var benchColors = func() []golor.Color {
	r := rand.New(rand.NewSource(42))
	out := make([]golor.Color, 1024)
	for i := range out {
		out[i] = golor.RGB(uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(r.Intn(256)))
	}
	return out
}()

func BenchmarkHex(b *testing.B) {
	hexes := make([]string, len(benchColors))
	for i, c := range benchColors {
		hexes[i] = c.String()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchColor, _ = golor.Hex(hexes[i&1023])
	}
}

func BenchmarkString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchString = benchColors[i&1023].String()
	}
}

func BenchmarkParseCSS(b *testing.B) {
	inputs := []struct{ name, in string }{
		{"hex", "#ff6b35"},
		{"named", "rebeccapurple"},
		{"rgb", "rgb(255 107 53 / 50%)"},
		{"rgbLegacy", "rgba(255, 107, 53, 0.5)"},
		{"hsl", "hsla(9, 100%, 60%, 0.8)"},
		{"oklch", "oklch(62.8% 0.2577 29.23)"},
		{"lab", "lab(54.29 80.82 69.91)"},
		{"color", "color(display-p3 1 0 0)"},
	}
	for _, tc := range inputs {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchColor, _ = golor.ParseCSS(tc.in)
			}
		})
	}
}

func BenchmarkNearestNamed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, benchColor = golor.NearestNamed(benchColors[i&1023])
	}
}

func BenchmarkImagePixels(b *testing.B) {
	img := image.NewNRGBA(image.Rect(0, 0, 512, 512))
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			img.SetNRGBA(x, y, color.NRGBA{uint8(x / 2), uint8(y / 2), uint8((x + y) % 256), 255})
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchPixels = golor.ImagePixels(img)
	}
}
