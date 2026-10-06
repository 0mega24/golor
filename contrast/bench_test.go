package contrast_test

import (
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/contrast"
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
	sinkF float64
	sinkC golor.Color
)

func BenchmarkRatio(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkF = contrast.Ratio(benchColors[i&1023], benchColors[(i+1)&1023])
	}
}

func BenchmarkEnforceContrast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkC = contrast.EnforceContrast(benchColors[i&1023], benchColors[(i+1)&1023], 7)
	}
}

func BenchmarkAPCA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkF = contrast.APCA(benchColors[i&1023], benchColors[(i+1)&1023])
	}
}
