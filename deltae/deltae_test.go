package deltae_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/deltae"
)

func TestDeltaE76Identical(t *testing.T) {
	cases := []golor.Color{
		golor.RGB(255, 0, 0),
		golor.RGB(0, 255, 0),
		golor.RGB(128, 128, 128),
	}
	for _, c := range cases {
		if d := deltae.DeltaE76(c, c); d != 0 {
			t.Errorf("DeltaE76(%v, %v) = %f, want 0", c, c, d)
		}
	}
}

func TestDeltaE76Symmetry(t *testing.T) {
	c1 := golor.RGB(255, 0, 0)
	c2 := golor.RGB(0, 0, 255)
	if d1, d2 := deltae.DeltaE76(c1, c2), deltae.DeltaE76(c2, c1); math.Abs(d1-d2) > 1e-9 {
		t.Errorf("DeltaE76 not symmetric: %f vs %f", d1, d2)
	}
}

func TestDeltaE2000Identical(t *testing.T) {
	c := golor.RGB(200, 100, 50)
	if d := deltae.DeltaE2000(c, c); d != 0 {
		t.Errorf("DeltaE2000(%v, %v) = %f, want 0", c, c, d)
	}
}

func TestEnsureDistinct(t *testing.T) {
	c1 := golor.RGB(100, 100, 100)
	c2 := golor.RGB(101, 101, 101)
	result := deltae.EnsureDistinct([]golor.Color{c1, c2}, 5.0)
	if len(result) != 2 {
		t.Fatalf("expected 2 colors, got %d", len(result))
	}
	if deltae.DeltaE76(result[0], result[1]) < 5.0 {
		t.Errorf("colors not distinct enough: ΔE76 = %f", deltae.DeltaE76(result[0], result[1]))
	}
}
