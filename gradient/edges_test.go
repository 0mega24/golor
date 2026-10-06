package gradient_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/gradient"
)

type gradFn func(a, b golor.Color, n int) []golor.Color

var allSpaces = map[string]gradFn{
	"RGB": gradient.RGB, "HSL": gradient.HSL, "LAB": gradient.LAB,
	"LCH": gradient.LCH, "OKLab": gradient.OKLab, "OKLCH": gradient.OKLCH,
	"HSV": gradient.HSV, "LinearRGB": gradient.LinearRGB,
}

func hueDist(a, b float64) float64 {
	d := math.Abs(a - b)
	return math.Min(d, 360-d)
}

// TestAchromaticEndpointBorrowsHue checks that greys do not introduce spurious hues.
func TestAchromaticEndpointBorrowsHue(t *testing.T) {
	blue := golor.RGB(0, 0, 255)
	for _, grey := range []golor.Color{golor.RGB(255, 255, 255), golor.RGB(128, 128, 128), golor.RGB(0, 0, 0)} {
		// Each space is checked against its own hue (CIE LCH blue drifts purple in OKLCH
		// terms; this is a property of CIELAB, rather than the interpolation).
		spaces := map[string]struct {
			fn  gradFn
			hue func(golor.Color) (h float64, chromatic bool)
		}{
			"LCH":   {gradient.LCH, func(c golor.Color) (float64, bool) { l := convert.ToLCH(c); return l.H, l.C > 3 }},
			"OKLCH": {gradient.OKLCH, func(c golor.Color) (float64, bool) { l := convert.ToOKLCH(c); return l.H, l.C > 0.02 }},
			"HSV":   {gradient.HSV, func(c golor.Color) (float64, bool) { h := convert.ToHSV(c); return h.H, h.S > 0.05 }},
			"HSL":   {gradient.HSL, func(c golor.Color) (float64, bool) { h := convert.ToHSL(c); return h.H, h.S > 0.05 }},
		}
		for name, sp := range spaces {
			want, _ := sp.hue(blue)
			for _, c := range sp.fn(grey, blue, 9)[1:8] {
				if h, chromatic := sp.hue(c); chromatic && hueDist(h, want) > 6 {
					t.Fatalf("%s %s->blue passes through %s (hue %.0f, want %.0f)", name, grey, c, h, want)
				}
			}
		}
	}
}

// TestPremultipliedAlpha checks that fading to transparent preserves color without darkening.
func TestPremultipliedAlpha(t *testing.T) {
	clear := golor.RGBA(0, 0, 0, 0)
	white := golor.RGB(255, 255, 255)
	for name, fn := range allSpaces {
		mid := fn(clear, white, 3)[1]
		if mid.String() != "#ffffff80" {
			t.Errorf("%s: transparent->white midpoint = %s, want #ffffff80", name, mid)
		}
	}
	// Equal alphas behave exactly like straight interpolation.
	a, b := golor.RGBA(255, 0, 0, 128), golor.RGBA(0, 0, 255, 128)
	if got := gradient.RGB(a, b, 3)[1]; got.R8() != 128 || got.B8() != 128 || got.A8() != 128 {
		t.Errorf("equal-alpha RGB midpoint = %s", got)
	}
}

// TestEndpointsAndRange checks that every space returns exact endpoints and in-range channels.
func TestEndpointsAndRange(t *testing.T) {
	pairs := [][2]golor.Color{
		{golor.RGB(255, 0, 0), golor.RGB(0, 255, 255)},
		{golor.RGB(12, 200, 80), golor.RGBA(250, 250, 20, 30)},
		{golor.RGB(0, 0, 0), golor.RGB(255, 255, 255)},
	}
	for name, fn := range allSpaces {
		for _, p := range pairs {
			out := fn(p[0], p[1], 17)
			if out[0].String() != p[0].String() || out[16].String() != p[1].String() {
				t.Fatalf("%s endpoints %s..%s, want %s..%s", name, out[0], out[16], p[0], p[1])
			}
			for _, c := range out {
				for _, v := range []float64{c.R, c.G, c.B, c.A} {
					if !(v >= 0 && v <= 1) {
						t.Fatalf("%s out-of-range channel in %v", name, c)
					}
				}
			}
		}
	}
}

// TestGamutMappedPolarGradients checks that LCH/OKLCH interpolation gamut-maps vivid colors, so
// lightness stays monotonic between endpoints of different lightness.
func TestGamutMappedPolarGradients(t *testing.T) {
	a, b := golor.RGB(0, 0, 255), golor.RGB(255, 255, 0)
	for name, fn := range map[string]gradFn{"LCH": gradient.LCH, "OKLCH": gradient.OKLCH} {
		prev := -1.0
		for _, c := range fn(a, b, 32) {
			l := convert.ToOKLab(c).L
			if l < prev-0.01 {
				t.Fatalf("%s lightness not monotonic: %v after %v", name, l, prev)
			}
			prev = l
		}
	}
}

func TestMultiStopMatchesPairwise(t *testing.T) {
	stops := []golor.Color{golor.RGB(255, 0, 0), golor.RGB(0, 255, 0), golor.RGB(0, 0, 255)}
	got := gradient.MultiStopOKLCH(stops, 5)
	first := gradient.OKLCH(stops[0], stops[1], 3)
	second := gradient.OKLCH(stops[1], stops[2], 3)
	want := []golor.Color{first[0], first[1], first[2], second[1], second[2]}
	for i := range want {
		if got[i].String() != want[i].String() {
			t.Fatalf("MultiStopOKLCH[%d] = %s, want %s", i, got[i], want[i])
		}
	}
	for _, fn := range []func([]golor.Color, int) []golor.Color{
		gradient.MultiStopRGB, gradient.MultiStopHSL, gradient.MultiStopLAB,
		gradient.MultiStopLCH, gradient.MultiStopOKLab, gradient.MultiStopOKLCH,
		gradient.MultiStopHSV, gradient.MultiStopLinearRGB,
	} {
		if fn(stops, 0) != nil || len(fn(stops[:1], 4)) != 1 || len(fn(stops, 10)) != 10 {
			t.Fatal("multistop size contract")
		}
	}
}

func TestLerpHueLargeAngles(t *testing.T) {
	// Hue values far outside [0,360) in struct literals still take the shorter arc.
	a := convert.FromHSL(convert.HSL{H: 350, S: 1, L: 0.5})
	b := convert.FromHSL(convert.HSL{H: 10, S: 1, L: 0.5})
	mid := convert.ToHSL(gradient.HSL(a, b, 3)[1])
	if hueDist(mid.H, 0) > 0.5 {
		t.Fatalf("350->10 midpoint hue %v, want 0", mid.H)
	}
}
