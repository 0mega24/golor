package adjust_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/adjust"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/deltae"
)

func TestSaturateLeavesGreyAlone(t *testing.T) {
	grey := golor.RGBA(128, 128, 128, 200)
	for name, got := range map[string]golor.Color{
		"Saturate":        adjust.Saturate(grey, 0.5),
		"SetSaturation":   adjust.SetSaturation(grey, 0.8),
		"SaturateOKLCH":   adjust.SaturateOKLCH(grey, 0.5),
		"DesaturateOKLCH": adjust.DesaturateOKLCH(grey, 0.5),
		"ShiftHueOKLCH":   adjust.ShiftHueOKLCH(grey, 90),
	} {
		if got != grey {
			t.Errorf("%s(grey) = %s, want unchanged", name, got)
		}
	}
}

func TestLightenOKLCHIsPerceptuallyEven(t *testing.T) {
	// HSL lightening moves yellow and blue by very different perceived amounts; OKLCH
	// lightening moves them by similar amounts.
	yellow, blue := golor.RGB(200, 200, 0), golor.RGB(0, 0, 200)
	dy := deltae.DeltaE2000(yellow, adjust.LightenOKLCH(yellow, 0.08))
	db := deltae.DeltaE2000(blue, adjust.LightenOKLCH(blue, 0.08))
	if math.Max(dy, db)/math.Min(dy, db) > 1.6 {
		t.Fatalf("OKLCH lighten dE00 yellow %.2f vs blue %.2f", dy, db)
	}
	hy := deltae.DeltaE2000(yellow, adjust.Lighten(yellow, 0.08))
	hb := deltae.DeltaE2000(blue, adjust.Lighten(blue, 0.08))
	if math.Max(hy, hb)/math.Min(hy, hb) < math.Max(dy, db)/math.Min(dy, db) {
		t.Fatalf("expected HSL to be less even than OKLCH")
	}
}

func TestOKLCHAdjustKeepsHueAndAlpha(t *testing.T) {
	c := golor.RGBA(40, 120, 200, 90)
	want := convert.ToOKLCH(c).H
	for name, got := range map[string]golor.Color{
		"Lighten":    adjust.LightenOKLCH(c, 0.2),
		"Darken":     adjust.DarkenOKLCH(c, 0.2),
		"Saturate":   adjust.SaturateOKLCH(c, 0.5), // Exceeds the gamut, then is gamut-mapped.
		"Desaturate": adjust.DesaturateOKLCH(c, 0.1),
	} {
		if got.A != c.A {
			t.Errorf("%s changed alpha", name)
		}
		if h := convert.ToOKLCH(got).H; math.Abs(h-want) > 8 { // CSS gamut mapping may clip within 1 JND.
			t.Errorf("%s hue %.1f, want %.1f", name, h, want)
		}
	}
	if got := convert.ToOKLCH(adjust.ShiftHueOKLCH(c, 120)).H; math.Abs(got-math.Mod(want+120, 360)) > 4 {
		t.Errorf("ShiftHueOKLCH hue %.1f", got)
	}
	if got := adjust.DesaturateOKLCH(c, 1); !convert.ToOKLCH(got).IsAchromatic() && convert.ToOKLCH(got).C > 1e-4 {
		t.Errorf("DesaturateOKLCH(1) = %s, want grey", got)
	}
}
