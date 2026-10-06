package convert_test

import (
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

type roundTrip struct {
	name string
	fn   func(golor.Color) golor.Color
	tol  float64 // Maximum per-channel float error for arbitrary in-gamut input.
}

var roundTrips = []roundTrip{
	{"HSL", func(c golor.Color) golor.Color { return convert.FromHSL(convert.ToHSL(c)) }, 1e-12},
	{"HSV", func(c golor.Color) golor.Color { return convert.FromHSV(convert.ToHSV(c)) }, 1e-12},
	{"CMYK", func(c golor.Color) golor.Color { return convert.FromCMYK(convert.ToCMYK(c)) }, 1e-12},
	{"XYZ", func(c golor.Color) golor.Color { return convert.FromXYZ(convert.ToXYZ(c)) }, 1e-12},
	{"LAB", func(c golor.Color) golor.Color { return convert.FromLAB(convert.ToLAB(c)) }, 1e-12},
	{"LCH", func(c golor.Color) golor.Color { return convert.FromLCH(convert.ToLCH(c)) }, 1e-12},
	{"OKLab", func(c golor.Color) golor.Color { return convert.FromOKLab(convert.ToOKLab(c)) }, 1e-9},
	{"OKLCH", func(c golor.Color) golor.Color { return convert.FromOKLCH(convert.ToOKLCH(c)) }, 1e-9},
	{"P3", func(c golor.Color) golor.Color { return convert.FromP3(convert.ToP3(c)) }, 1e-12},
	{"LABMapped", func(c golor.Color) golor.Color { return convert.FromLABMapped(convert.ToLAB(c)) }, 1e-12},
	{"OKLCHMapped", func(c golor.Color) golor.Color { return convert.FromOKLCHMapped(convert.ToOKLCH(c)) }, 1e-9},
}

func maxErr(a, b golor.Color) float64 {
	return math.Max(math.Abs(a.R-b.R), math.Max(math.Abs(a.G-b.G), math.Abs(a.B-b.B)))
}

// TestRoundTrip8BitLossless samples the 8-bit grid; GOLOR_FULL_SWEEP=1 checks every color.
func TestRoundTrip8BitLossless(t *testing.T) {
	step := 3
	if os.Getenv("GOLOR_FULL_SWEEP") != "" {
		step = 1
	}
	levels := []int{}
	for v := 0; v < 256; v += step {
		levels = append(levels, v)
	}
	if levels[len(levels)-1] != 255 {
		levels = append(levels, 255)
	}
	for _, rt := range roundTrips {
		t.Run(rt.name, func(t *testing.T) {
			check := func(c golor.Color) {
				got := rt.fn(c)
				if got.R8() != c.R8() || got.G8() != c.G8() || got.B8() != c.B8() {
					t.Fatalf("%s -> %s", c, got)
				}
			}
			for _, r := range levels {
				for _, g := range levels {
					for _, b := range levels {
						check(golor.RGB(uint8(r), uint8(g), uint8(b)))
					}
				}
			}
			for v := 0; v < 256; v++ {
				check(golor.RGB(uint8(v), uint8(v), uint8(v)))
			}
		})
	}
}

// TestRoundTripFloatPrecision checks that off-grid float inputs round-trip within each space's tolerance.
func TestRoundTripFloatPrecision(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	inputs := make([]golor.Color, 0, 20002)
	for i := 0; i < 20000; i++ {
		inputs = append(inputs, golor.RGBf(r.Float64(), r.Float64(), r.Float64()))
	}
	// Exclude 0.04045: the sRGB segments do not meet exactly at that breakpoint.
	inputs = append(inputs, golor.RGBf(1e-9, 0.0404, 0.0031308), golor.RGBf(1-1e-12, 1e-6, 0.5))
	for _, rt := range roundTrips {
		t.Run(rt.name, func(t *testing.T) {
			worst := 0.0
			for _, c := range inputs {
				worst = math.Max(worst, maxErr(c, rt.fn(c)))
			}
			if worst > rt.tol {
				t.Fatalf("max error %.3e > %.0e", worst, rt.tol)
			}
		})
	}
}

// TestRoundTripNoDrift checks that repeated conversions do not drift.
func TestRoundTripNoDrift(t *testing.T) {
	c := golor.RGBf(0.8123, 0.3141, 0.0592)
	for _, rt := range roundTrips {
		d := c
		for i := 0; i < 1000; i++ {
			d = rt.fn(d)
		}
		if e := maxErr(c, d); e > rt.tol*10 {
			t.Errorf("%s: drift %.3e after 1000 round trips", rt.name, e)
		}
	}
}

// TestConversionChain checks that a chain through every space returns the original color.
func TestConversionChain(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for i := 0; i < 5000; i++ {
		c := golor.RGBf(r.Float64(), r.Float64(), r.Float64())
		d := convert.FromHSL(convert.ToHSL(c))
		d = convert.FromOKLCH(convert.ToOKLCH(d))
		d = convert.FromLCH(convert.ToLCH(d))
		d = convert.FromP3(convert.ToP3(d))
		d = convert.FromHSV(convert.ToHSV(d))
		d = convert.FromCMYK(convert.ToCMYK(d))
		if e := maxErr(c, d); e > 1e-9 {
			t.Fatalf("chain error %.3e for %v", e, c)
		}
	}
}

func TestSRGBAlwaysInGamut(t *testing.T) {
	for v := 0; v < 1<<24; v += 4099 {
		c := golor.RGB(uint8(v>>16), uint8(v>>8), uint8(v))
		if !convert.ToLAB(c).InGamut() || !convert.ToLCH(c).InGamut() ||
			!convert.ToOKLab(c).InGamut() || !convert.ToOKLCH(c).InGamut() || !convert.ToP3(c).InSRGBGamut() {
			t.Fatalf("%s reported out of gamut", c)
		}
	}
	if (convert.OKLCH{L: 0.7, C: 0.35, H: 150}).InGamut() {
		t.Fatal("vivid OKLCH green should be outside sRGB")
	}
	if (convert.P3{R: 0, G: 1, B: 0}).InSRGBGamut() {
		t.Fatal("P3 green should be outside sRGB")
	}
}

func TestGamutMappingPreservesHue(t *testing.T) {
	in := convert.LCH{L: 70, C: 100, H: 200}
	if in.InGamut() {
		t.Fatal("test color should be out of gamut")
	}
	clipped := convert.ToLCH(convert.FromLCH(in))
	mapped := convert.ToLCH(convert.FromLCHMapped(in))
	hueDiff := func(a, b float64) float64 { d := math.Abs(a - b); return math.Min(d, 360-d) }
	// Mapping keeps lightness and hue much closer than per-channel clipping does.
	if math.Abs(mapped.L-in.L) > 2 || math.Abs(mapped.L-in.L) >= math.Abs(clipped.L-in.L) {
		t.Fatalf("mapped L %.2f, clipped L %.2f, want %.0f", mapped.L, clipped.L, in.L)
	}
	if hueDiff(mapped.H, in.H) > 5 || hueDiff(mapped.H, in.H) >= hueDiff(clipped.H, in.H) {
		t.Fatalf("mapped H %.1f, clipped H %.1f, want %.0f", mapped.H, clipped.H, in.H)
	}
	p3Green := convert.FromP3Mapped(convert.P3{G: 1})
	if p3Green.G8() < 200 || p3Green.R8() > 120 {
		t.Fatalf("P3 green mapped to %s", p3Green)
	}
}

func TestGreysAreAchromatic(t *testing.T) {
	for v := 0; v < 256; v++ {
		c := golor.RGB(uint8(v), uint8(v), uint8(v))
		lch, ok, hsl, hsv := convert.ToLCH(c), convert.ToOKLCH(c), convert.ToHSL(c), convert.ToHSV(c)
		if !lch.IsAchromatic() || lch.H != 0 || !ok.IsAchromatic() || ok.H != 0 ||
			!hsl.IsAchromatic() || !hsv.IsAchromatic() {
			t.Fatalf("%s: LCH %+v OKLCH %+v", c, lch, ok)
		}
		lab := convert.ToLAB(c)
		if math.Abs(lab.A) > 1e-12 || math.Abs(lab.B) > 1e-12 {
			t.Fatalf("%s has a*=%g b*=%g", c, lab.A, lab.B)
		}
	}
}

func TestHueWrapping(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for i := 0; i < 2000; i++ {
		h := r.Float64() * 360
		s, l := r.Float64(), r.Float64()
		for _, k := range []float64{-3, -1, 1, 2, 10} {
			want := convert.FromHSL(convert.HSL{H: h, S: s, L: l})
			if got := convert.FromHSL(convert.HSL{H: h + 360*k, S: s, L: l}); maxErr(got, want) > 1e-9 {
				t.Fatalf("FromHSL(H=%v) = %s, want %s", h+360*k, got, want)
			}
			wantV := convert.FromHSV(convert.HSV{H: h, S: s, V: l})
			if got := convert.FromHSV(convert.HSV{H: h + 360*k, S: s, V: l}); maxErr(got, wantV) > 1e-9 {
				t.Fatalf("FromHSV(H=%v) = %s, want %s", h+360*k, got, wantV)
			}
		}
	}
}

func TestOutOfRangeAndNaNInputs(t *testing.T) {
	hot := golor.Color{R: 2, G: 0.5, B: -1, A: 1}
	if hsl := convert.ToHSL(hot); hsl.S > 1 || hsl.L > 1 {
		t.Fatalf("ToHSL out of range: %+v", hsl)
	}
	if cmyk := convert.ToCMYK(hot); cmyk.K < 0 || cmyk.Y > 1 {
		t.Fatalf("ToCMYK out of range: %+v", cmyk)
	}
	if lab := convert.ToLAB(hot); lab.L > 100+1e-9 {
		t.Fatalf("ToLAB out of range: %+v", lab)
	}
	nan := golor.Color{R: math.NaN(), A: 1}
	if lab := convert.ToLAB(nan); math.IsNaN(lab.L) {
		t.Fatal("NaN propagated through ToLAB")
	}
	if c := convert.FromHSL(convert.HSL{H: math.NaN(), S: 2, L: 0.5}); math.IsNaN(c.R) {
		t.Fatal("NaN propagated through FromHSL")
	}
}

func TestReferenceValues(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	lab := convert.ToLAB(red)
	// Lindbloom's sRGB D65 values differ from the CSS Color 4 matrices in the 4th decimal.
	if math.Abs(lab.L-53.2408) > 0.01 || math.Abs(lab.A-80.0925) > 0.02 || math.Abs(lab.B-67.2032) > 0.02 {
		t.Fatalf("ToLAB(red) = %+v", lab)
	}
	ok := convert.ToOKLCH(red)
	if math.Abs(ok.L-0.62796) > 1e-4 || math.Abs(ok.C-0.25768) > 1e-4 || math.Abs(ok.H-29.2339) > 0.01 {
		t.Fatalf("ToOKLCH(red) = %+v", ok)
	}
	p3 := convert.ToP3(red)
	if math.Abs(p3.R-0.9175) > 1e-3 || math.Abs(p3.G-0.2003) > 1e-3 || math.Abs(p3.B-0.1386) > 1e-3 {
		t.Fatalf("ToP3(red) = %+v", p3)
	}
	white := convert.ToXYZ(golor.RGB(255, 255, 255))
	if math.Abs(white.X-0.3127/0.3290) > 1e-12 || math.Abs(white.Y-1) > 1e-12 {
		t.Fatalf("ToXYZ(white) = %+v", white)
	}
	if d := convert.DeltaEOK(golor.RGB(255, 0, 0), golor.RGB(255, 0, 0)); d != 0 {
		t.Fatalf("DeltaEOK(self) = %v", d)
	}
}
