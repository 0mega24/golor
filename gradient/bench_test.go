package gradient_test

import (
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/gradient"
)

var sinkColors []golor.Color

func BenchmarkGradient256(b *testing.B) {
	from, to := golor.RGB(255, 0, 0), golor.RGB(0, 0, 255)
	for name, fn := range allSpaces {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sinkColors = fn(from, to, 256)
			}
		})
	}
}

func BenchmarkMultiStop256(b *testing.B) {
	stops := []golor.Color{
		golor.RGB(255, 0, 0), golor.RGB(255, 200, 0), golor.RGB(0, 160, 80),
		golor.RGB(0, 80, 255), golor.RGB(120, 0, 200),
	}
	b.Run("LCH", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkColors = gradient.MultiStopLCH(stops, 256)
		}
	})
	b.Run("OKLCH", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkColors = gradient.MultiStopOKLCH(stops, 256)
		}
	})
}
