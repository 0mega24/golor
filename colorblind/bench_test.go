package colorblind_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/colorblind"
)

var benchColors = func() []golor.Color {
	r := rand.New(rand.NewSource(42))
	out := make([]golor.Color, 1024)
	for i := range out {
		out[i] = golor.RGB(uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(r.Intn(256)))
	}
	return out
}()

var (
	sink  golor.Color
	sinkS []golor.Color
)

func BenchmarkSimulate(b *testing.B) {
	for _, d := range []struct {
		name string
		d    colorblind.Deficiency
	}{{"Deuteranopia", colorblind.Deuteranopia}, {"Protanopia", colorblind.Protanopia}, {"Tritanopia", colorblind.Tritanopia}} {
		b.Run(d.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sink = colorblind.Simulate(benchColors[i&1023], d.d)
			}
		})
	}
}

func BenchmarkAccessiblePalette(b *testing.B) {
	for _, n := range []int{8, 32, 128} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sinkS = colorblind.AccessiblePalette(benchColors[:n], colorblind.Deuteranopia)
			}
		})
	}
}
