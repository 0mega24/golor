package deltae_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/deltae"
)

var (
	benchColors = func() []golor.Color {
		r := rand.New(rand.NewSource(42))
		out := make([]golor.Color, 1024)
		for i := range out {
			out[i] = golor.RGB(uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(r.Intn(256)))
		}
		return out
	}()
	sinkF  float64
	sinkCs []golor.Color
)

func BenchmarkDeltaE76(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkF = deltae.DeltaE76(benchColors[i&1023], benchColors[(i+1)&1023])
	}
}

func BenchmarkDeltaE2000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkF = deltae.DeltaE2000(benchColors[i&1023], benchColors[(i+1)&1023])
	}
}

func BenchmarkDeltaEOK(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkF = deltae.DeltaEOK(benchColors[i&1023], benchColors[(i+1)&1023])
	}
}

func BenchmarkEnsureDistinct(b *testing.B) {
	for _, n := range []int{8, 32, 128} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sinkCs = deltae.EnsureDistinct(benchColors[:n], 10)
			}
		})
	}
}
