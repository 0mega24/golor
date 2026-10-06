package deltae_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/deltae"
)

func minPair(cs []golor.Color) float64 {
	m := math.Inf(1)
	for i := range cs {
		for j := i + 1; j < len(cs); j++ {
			m = math.Min(m, deltae.DeltaE76(cs[i], cs[j]))
		}
	}
	return m
}

func TestEnsureDistinctRandomPalettes(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for trial := 0; trial < 500; trial++ {
		cs := make([]golor.Color, 6)
		for i := range cs {
			cs[i] = golor.RGB(uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(r.Intn(256)))
		}
		out := deltae.EnsureDistinct(cs, 15)
		if d := minPair(out); d < 15 {
			t.Fatalf("trial %d: min pair %.2f < 15 for %v -> %v", trial, d, cs, out)
		}
		if out[0] != cs[0] {
			t.Fatal("first color must never move when later colors can")
		}
	}
}

func TestEnsureDistinctIdenticalInputs(t *testing.T) {
	purple := golor.RGBA(120, 60, 200, 200)
	in := []golor.Color{purple, purple, purple, purple, purple}
	out := deltae.EnsureDistinct(in, 20)
	if d := minPair(out); d < 20 {
		t.Fatalf("min pair %.2f < 20: %v", d, out)
	}
	black := 0
	for _, c := range out {
		if c.A != purple.A {
			t.Fatal("alpha not preserved")
		}
		if c.R8() < 8 && c.G8() < 8 && c.B8() < 8 {
			black++
		}
	}
	if black > 1 {
		t.Fatalf("identical inputs collapsed toward black: %v", out)
	}
}

func TestEnsureDistinctUnreachableTerminates(t *testing.T) {
	cs := make([]golor.Color, 40)
	for i := range cs {
		cs[i] = golor.RGB(100, 100, 100)
	}
	out := deltae.EnsureDistinct(cs, 60) // Unreachable: 40 colors separated by 60.
	if len(out) != len(cs) {
		t.Fatal("length changed")
	}
	if minPair(out) <= 0 {
		t.Fatal("expected best-effort spreading")
	}
	if got := deltae.EnsureDistinct(cs[:1], 10); len(got) != 1 {
		t.Fatal("single color")
	}
	if got := deltae.EnsureDistinct(cs[:3], math.NaN()); got[1] != cs[1] {
		t.Fatal("NaN threshold should be a no-op")
	}
}

func TestDeltaEOK(t *testing.T) {
	a, b := golor.RGB(255, 0, 0), golor.RGB(250, 5, 5)
	if d := deltae.DeltaEOK(a, a); d != 0 {
		t.Fatalf("self distance %v", d)
	}
	if d := deltae.DeltaEOK(a, b); d <= 0 || d > 0.05 {
		t.Fatalf("near-red distance %v", d)
	}
}
