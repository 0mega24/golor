package convert_test

import (
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
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
	sinkColor golor.Color
	sinkHSL   convert.HSL
	sinkHSV   convert.HSV
	sinkLAB   convert.LAB
	sinkLCH   convert.LCH
	sinkOKLab convert.OKLab
	sinkOKLCH convert.OKLCH
	sinkCMYK  convert.CMYK
	sinkP3    convert.P3
	sinkBool  bool
)

func benchTo[T any](b *testing.B, sink *T, fn func(golor.Color) T) {
	for i := 0; i < b.N; i++ {
		*sink = fn(benchColors[i&1023])
	}
}

func benchFrom[T any](b *testing.B, to func(golor.Color) T, from func(T) golor.Color) {
	in := make([]T, 1024)
	for i := range in {
		in[i] = to(benchColors[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkColor = from(in[i&1023])
	}
}

func BenchmarkToHSL(b *testing.B)   { benchTo(b, &sinkHSL, convert.ToHSL) }
func BenchmarkToHSV(b *testing.B)   { benchTo(b, &sinkHSV, convert.ToHSV) }
func BenchmarkToCMYK(b *testing.B)  { benchTo(b, &sinkCMYK, convert.ToCMYK) }
func BenchmarkToLAB(b *testing.B)   { benchTo(b, &sinkLAB, convert.ToLAB) }
func BenchmarkToLCH(b *testing.B)   { benchTo(b, &sinkLCH, convert.ToLCH) }
func BenchmarkToOKLab(b *testing.B) { benchTo(b, &sinkOKLab, convert.ToOKLab) }
func BenchmarkToOKLCH(b *testing.B) { benchTo(b, &sinkOKLCH, convert.ToOKLCH) }
func BenchmarkToP3(b *testing.B)    { benchTo(b, &sinkP3, convert.ToP3) }

// BenchmarkToLABOffGrid uses off-grid floats to bypass the 8-bit linearization lookup table.
func BenchmarkToLABOffGrid(b *testing.B) {
	in := make([]golor.Color, 1024)
	for i := range in {
		c := benchColors[i]
		in[i] = golor.RGBf(c.R*0.999+0.0003, c.G*0.999+0.0003, c.B*0.999+0.0003)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkLAB = convert.ToLAB(in[i&1023])
	}
}

func BenchmarkFromHSL(b *testing.B)   { benchFrom(b, convert.ToHSL, convert.FromHSL) }
func BenchmarkFromHSV(b *testing.B)   { benchFrom(b, convert.ToHSV, convert.FromHSV) }
func BenchmarkFromCMYK(b *testing.B)  { benchFrom(b, convert.ToCMYK, convert.FromCMYK) }
func BenchmarkFromLAB(b *testing.B)   { benchFrom(b, convert.ToLAB, convert.FromLAB) }
func BenchmarkFromLCH(b *testing.B)   { benchFrom(b, convert.ToLCH, convert.FromLCH) }
func BenchmarkFromOKLab(b *testing.B) { benchFrom(b, convert.ToOKLab, convert.FromOKLab) }
func BenchmarkFromOKLCH(b *testing.B) { benchFrom(b, convert.ToOKLCH, convert.FromOKLCH) }
func BenchmarkFromP3(b *testing.B)    { benchFrom(b, convert.ToP3, convert.FromP3) }

// BenchmarkFromOKLCHMappedOutOfGamut maps colors far outside sRGB (worst case: full bisection).
func BenchmarkFromOKLCHMappedOutOfGamut(b *testing.B) {
	in := make([]convert.OKLCH, 1024)
	for i := range in {
		in[i] = convert.OKLCH{L: 0.3 + 0.6*float64(i%7)/7, C: 0.4, H: float64(i) * 360 / 1024}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkColor = convert.FromOKLCHMapped(in[i&1023])
	}
}

func BenchmarkOKLCHInGamut(b *testing.B) {
	o := convert.OKLCH{L: 0.7, C: 0.2, H: 150}
	for i := 0; i < b.N; i++ {
		sinkBool = o.InGamut()
	}
}
