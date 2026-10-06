package harmony_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/harmony"
)

func TestOKLCHHarmonyKeepsLightness(t *testing.T) {
	c := golor.RGB(230, 180, 40)
	base := convert.ToOKLCH(c)
	tri := harmony.TriadicOKLCH(c)
	tet := harmony.TetradicOKLCH(c)
	comp := harmony.ComplementaryOKLCH(c)
	split := harmony.SplitComplementaryOKLCH(c)
	ana := harmony.AnalogousOKLCH(c, 30)
	schemes := [][]golor.Color{tri[:], tet[:], comp[:], split[:], ana[:]}
	for _, s := range schemes {
		for _, col := range s {
			if l := convert.ToOKLCH(col).L; math.Abs(l-base.L) > 0.03 {
				t.Fatalf("OKLCH harmony changed lightness %.3f -> %.3f (%s)", base.L, l, col)
			}
		}
	}
	h := convert.ToOKLCH(comp[1]).H
	if d := math.Abs(math.Mod(h-base.H+360, 360) - 180); d > 8 {
		t.Fatalf("complement hue %.1f vs base %.1f", h, base.H)
	}
	// HSL rotation of the same yellow changes perceived lightness far more.
	hsl := harmony.Triadic(c)
	if math.Abs(convert.ToOKLCH(hsl[2]).L-base.L) < 0.1 {
		t.Fatalf("expected HSL triadic to shift lightness")
	}
}

func TestOKLCHHarmonyGreyAndAlpha(t *testing.T) {
	grey := golor.RGBA(90, 90, 90, 40)
	if got := harmony.ComplementaryOKLCH(grey); got[1] != grey {
		t.Fatalf("grey complement = %s", got[1])
	}
	c := golor.RGBA(10, 120, 200, 77)
	for _, col := range harmony.TriadicOKLCH(c) {
		if col.A != c.A {
			t.Fatal("alpha not preserved")
		}
	}
	steps := harmony.ExpandOKLCH(c, 5)
	for i := 1; i < len(steps); i++ {
		if convert.ToOKLCH(steps[i]).L <= convert.ToOKLCH(steps[i-1]).L {
			t.Fatal("ExpandOKLCH not increasing in lightness")
		}
	}
	if harmony.ExpandOKLCH(c, 0) != nil || len(harmony.ExpandOKLCH(c, 1)) != 1 {
		t.Fatal("ExpandOKLCH size contract")
	}
}
