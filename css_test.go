package golor_test

import (
	"image"
	stdcolor "image/color"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustCSS(t *testing.T, s string) golor.Color {
	t.Helper()
	c, err := golor.ParseCSS(s)
	require.NoError(t, err, s)
	return c
}

func TestParseCSSValid(t *testing.T) {
	cases := []struct{ in, want string }{
		// Hex colors.
		{"#fff", "#ffffff"},
		{"#FFF", "#ffffff"},
		{"#ffff", "#ffffff"},
		{"#f008", "#ff000088"},
		{"#FF6B35", "#ff6b35"},
		{" #ff6b35 ", "#ff6b35"},
		{"#ff6b3580", "#ff6b3580"},
		// Named colors.
		{"RebeccaPurple", "#663399"},
		{"Transparent", "#00000000"},
		// rgb(): modern and legacy syntax, aliases, clamping, and none.
		{"rgb(255 0 0)", "#ff0000"},
		{"RGB(255 0 0)", "#ff0000"},
		{"rgb(300 0 0)", "#ff0000"},
		{"rgb(-5 0 0)", "#000000"},
		{"rgb(50% 0 0)", "#800000"},
		{"rgb(50% 128 0)", "#808000"},
		{"rgb(150% 0 0)", "#ff0000"},
		{"rgb(255 0 0 / 150%)", "#ff0000"},
		{"rgb(255 0 0 / 1.5)", "#ff0000"},
		{"rgb(255 0 0 / -1)", "#ff000000"},
		{"rgb(255 0 0 / 50%)", "#ff000080"},
		{"rgb(255, 0, 0, 0.5)", "#ff000080"},
		{"rgba(255, 0, 0, 50%)", "#ff000080"},
		{"rgba(255 0 0)", "#ff0000"},
		{"rgba(255,0,0)", "#ff0000"},
		{"rgb(100%, 0%, 0%)", "#ff0000"},
		{"rgb(255 0 none)", "#ff0000"},
		{"rgb(none none none / none)", "#00000000"},
		{"rgb(1e2 0 0)", "#640000"},
		{"rgb(+127.5 .5 0)", "#800100"},
		{"rgb(\t255\n0  0 )", "#ff0000"},
		// hsl(): hue units, wrapping, clamping, and numeric saturation and lightness.
		{"hsl(120 100% 50%)", "#00ff00"},
		{"hsl(120deg 100% 50%)", "#00ff00"},
		{"HSL(120DEG 100% 50%)", "#00ff00"},
		{"hsl(360 100% 50%)", "#ff0000"},
		{"hsl(400 100% 50%)", "#ffaa00"},
		{"hsl(-30 100% 50%)", "#ff0080"},
		{"hsl(0.5turn 100% 50%)", "#00ffff"},
		{"hsl(200grad 100% 50%)", "#00ffff"},
		{"hsl(3.14159265rad 100% 50%)", "#00ffff"},
		{"hsl(120 100 50)", "#00ff00"},
		{"hsl(120 120% 50%)", "#00ff00"},
		{"hsl(120 100% 150%)", "#ffffff"},
		{"hsl(none 0% 50%)", "#808080"},
		{"hsla(120, 100%, 50%, 0.5)", "#00ff0080"},
		{"hsl(120, 100%, 50%)", "#00ff00"},
		{"hsla(120 100% 50%)", "#00ff00"},
		{"hsl(9 100% 60% / 80%)", "#ff5233cc"},
		// hwb() colors.
		{"hwb(120 0% 0%)", "#00ff00"},
		{"hwb(0 50% 50%)", "#808080"},
		{"hwb(0 60% 60%)", "#808080"},
		{"hwb(0 0 0 / 0.5)", "#ff000080"},
		// lab(), lch(), oklab(), and oklch(): CSS reference values for sRGB red.
		{"lab(54.29 80.82 69.91)", "#ff0000"},
		{"lab(54.29% 64.656% 55.928%)", "#ff0000"},
		{"lch(54.29 106.84 40.85)", "#ff0000"},
		{"oklab(0.62796 0.22486 0.12585)", "#ff0000"},
		{"oklab(62.796% 56.215% 31.4625%)", "#ff0000"},
		{"oklch(62.8% 0.2577 29.23)", "#ff0000"},
		{"OKLCH(0.628 64.425% 29.23deg)", "#ff0000"},
		{"lab(100 0 0)", "#ffffff"},
		{"lab(150 0 0)", "#ffffff"},
		{"lab(0 0 0)", "#000000"},
		{"oklch(1 0.3 120)", "#ffffff"},
		{"oklch(0 0.3 120)", "#000000"},
		{"oklch(0.5 -0.1 120)", "#636363"},
		{"lab(50 0 0 / 25%)", "#77777740"},
		{"oklch(0.5 0 none)", "#636363"},
		// color() colors.
		{"color(srgb 1 0.5 0)", "#ff8000"},
		{"color(srgb 100% 50% 0%)", "#ff8000"},
		{"color(srgb-linear 1 0.2158605 0)", "#ff8000"},
		{"color(display-p3 0.9175 0.2003 0.1386)", "#ff0000"},
		{"color(xyz 0.4124 0.2126 0.0193)", "#ff0000"},
		{"color(xyz-d65 0.4124 0.2126 0.0193)", "#ff0000"},
		{"color(xyz-d50 0.43607 0.22249 0.01392)", "#ff0000"},
		{"color(srgb 1 0 0 / 0.5)", "#ff000080"},
		{"COLOR(SRGB none 1 none)", "#00ff00"},
	}
	for _, tc := range cases {
		got := mustCSS(t, tc.in)
		assert.Equal(t, tc.want, got.String(), tc.in)
	}
}

func TestParseCSSInvalid(t *testing.T) {
	cases := []string{
		"", "   ", "#ff", "#fffff", "#fffffff", "#ggg", "##fff", "#ff00zz",
		"fff", "notacolor", "currentcolor", "CurrentColor",
		"rgb(255,0,0,)", "rgb(,255,0,0)", "rgb(255,,0,0)", "rgb(255, 0 0)", "rgb(255,0,0,1,1)",
		"rgb(255 0 0 / 0.5 / 0.2)", "rgb(255 0 0 /)", "rgb(255 0 0 / 0.5 0.2)", "rgb(255, 0, 0 / 0.5)",
		"rgb(255,0,none)", "rgb(50%, 128, 0)", "rgb(255 0)", "rgb(1 2 3 4)", "rgb()", "rgb( )",
		"rgb(1. 2 3)", "rgb(0x10 0 0)", "rgb(inf 0 0)", "rgb(nan 0 0)", "rgb(1e999 0 0)", "rgb(1_0 0 0)",
		"rgb(10deg 0 0)", "rgb (255 0 0)", "rgb(255 0 0", "rgb 255 0 0)", "rgb(calc(1) 0 0)",
		"hsl(999deg)", "hsl(120, 100, 50)", "hsl(10% 100% 50%)", "hsl(none, 100%, 50%)",
		"hsl(120 100% 50%, 1)", "hsl(120deg, 100%, 50%, 0.5, 1)", "hsl(1e999deg 100% 50%)",
		"hsl(10foo 100% 50%)", "hsl(deg 100% 50%)",
		"hwb(120, 0%, 0%)", "hwb(120 0%)", "lab(50, 0, 0)", "lab(50 0)", "lch(50 0 10%)",
		"oklab(0.5 0 0 0)", "oklch(0.5, 0.1, 30)",
		"color(rec2020 1 0 0)", "color(srgb 1 0)", "color(1 0 0)", "color(srgb 1 0 0 0)", "color(srgb 1deg 0 0)",
		"hsv(120 100% 50%)",
	}
	for _, in := range cases {
		_, err := golor.ParseCSS(in)
		assert.Error(t, err, "%q should not parse", in)
	}
	_, err := golor.ParseCSS("currentcolor")
	assert.ErrorContains(t, err, "context")
}

func TestParseRGBAndHSLStillScoped(t *testing.T) {
	c, err := golor.ParseRGB("rgba(255 0 0)")
	require.NoError(t, err)
	assert.Equal(t, "#ff0000", c.String())
	_, err = golor.ParseRGB("hsl(0 100% 50%)")
	assert.Error(t, err)
	c, err = golor.ParseHSL("HSLA(0.25turn 100% 50%)")
	require.NoError(t, err)
	assert.Equal(t, "#80ff00", c.String())
	_, err = golor.ParseHSL("rgb(0 0 0)")
	assert.Error(t, err)
}

func TestParseCSSGamutMapsWideColors(t *testing.T) {
	inGamut := func(c golor.Color) bool {
		return c.R >= 0 && c.R <= 1 && c.G >= 0 && c.G <= 1 && c.B >= 0 && c.B <= 1
	}
	p3Red := mustCSS(t, "color(display-p3 1 0 0)")
	assert.True(t, inGamut(p3Red))
	assert.Equal(t, uint8(255), p3Red.R8(), p3Red.String())
	assert.Less(t, int(p3Red.G8()), 60, p3Red.String())
	assert.Less(t, int(p3Red.B8()), 60, p3Red.String())

	// A vivid OKLCH green is outside sRGB; mapping keeps it green rather than clipping
	// to a different hue.
	green := mustCSS(t, "oklch(0.7 0.35 145)")
	assert.True(t, inGamut(green))
	assert.Greater(t, green.G, green.R)
	assert.Greater(t, green.G, green.B)
	ref := mustCSS(t, "oklch(0.7 0.05 145)")
	_, _, _ = ref.R, ref.G, ref.B
}

// formats lists every CSS string method that preserves 8-bit colors.
var formats = []struct {
	name string
	fn   func(golor.Color) string
}{
	{"CSSRGBString", golor.Color.CSSRGBString},
	{"CSSHSLString", golor.Color.CSSHSLString},
	{"CSSOKLCHString", golor.Color.CSSOKLCHString},
	{"String", golor.Color.String},
}

// TestCSSStringsRoundTrip samples CSS round trips; GOLOR_FULL_SWEEP=1 checks every opaque color.
func TestCSSStringsRoundTrip(t *testing.T) {
	step := 5
	if os.Getenv("GOLOR_FULL_SWEEP") != "" {
		step = 1
	}
	var levels []int
	for v := 0; v < 256; v += step {
		levels = append(levels, v)
	}
	if levels[len(levels)-1] != 255 {
		levels = append(levels, 255)
	}
	for _, f := range formats {
		t.Run(f.name, func(t *testing.T) {
			check := func(c golor.Color) {
				s := f.fn(c)
				back, err := golor.ParseCSS(s)
				if err != nil {
					t.Fatalf("%s -> %q -> %v", c, s, err)
				}
				if back.String() != c.String() {
					t.Fatalf("%s -> %q -> %s", c, s, back)
				}
			}
			for _, r := range levels {
				for _, g := range levels {
					for _, b := range levels {
						check(golor.RGB(uint8(r), uint8(g), uint8(b)))
					}
				}
			}
			for a := 0; a < 256; a++ {
				check(golor.RGBA(200, 30, 90, uint8(a)))
				check(golor.RGBA(0, 0, 0, uint8(a)))
			}
			// Off-grid floats round to the same 8-bit color.
			check(golor.RGBAf(0.123456, 0.987654, 0.5, 0.3333))
		})
	}
}

func TestCSSStringFormatting(t *testing.T) {
	assert.Equal(t, "rgb(255 0 0 / 0.502)", golor.RGBA(255, 0, 0, 128).CSSRGBString())
	assert.Equal(t, "hsl(0 100% 50%)", golor.RGB(255, 0, 0).CSSHSLString())
	assert.Equal(t, "hsl(16.04 100% 60.39%)", golor.RGB(255, 107, 53).CSSHSLString())
	assert.Equal(t, "hsl(0 100% 23.73%)", golor.RGB(0x79, 0, 0).CSSHSLString())
	assert.Equal(t, "hsl(359.5 100% 23.73%)", golor.RGB(0x79, 0, 1).CSSHSLString())
	assert.Equal(t, "oklch(62.796% 0.25768 29.234)", golor.RGB(255, 0, 0).CSSOKLCHString())
	assert.Equal(t, "oklch(100% 0 0)", golor.RGB(255, 255, 255).CSSOKLCHString())
	assert.Equal(t, "oklch(0% 0 0 / 0.502)", golor.RGBA(0, 0, 0, 128).CSSOKLCHString())
}

func TestHexForms(t *testing.T) {
	cases := map[string]string{
		"fff": "#ffffff", "#abc": "#aabbcc", "abcd": "#aabbccdd", "#ABCDEF": "#abcdef",
		"\t#00ff0080\n": "#00ff0080", "#000f": "#000000",
	}
	for in, want := range cases {
		c, err := golor.Hex(in)
		require.NoError(t, err, in)
		assert.Equal(t, want, c.String(), in)
	}
	for _, in := range []string{"", "#", "#f", "#ff", "#fffff", "#fffffff", "#fffffffff", "##fff", "#+ff", "#0xffff", "# fff"} {
		_, err := golor.Hex(in)
		assert.Error(t, err, in)
	}
	_, err := golor.Hex("#zzz")
	assert.EqualError(t, err, `golor: invalid hex color "#zzz"`)
}

func TestNaNAndOutOfRange(t *testing.T) {
	c := golor.RGBf(math.NaN(), 0.5, math.Inf(1))
	assert.Equal(t, 0.0, c.R)
	assert.Equal(t, 1.0, c.B)
	assert.Equal(t, 0.0, golor.RGBAf(0, 0, 0, math.NaN()).A)

	lit := golor.Color{R: math.NaN(), G: 2, B: -1, A: math.NaN()}
	assert.Equal(t, "#00ff0000", lit.String())
	assert.Equal(t, uint8(0), lit.R8())
	r, g, b, a := lit.RGBA()
	assert.Equal(t, [4]uint32{0, 0, 0, 0}, [4]uint32{r, g, b, a})
}

func TestStringFormat(t *testing.T) {
	assert.Equal(t, "#0a0b0c", golor.RGB(10, 11, 12).String())
	assert.Equal(t, "#0a0b0cfe", golor.RGBA(10, 11, 12, 254).String())
	assert.Equal(t, "#ffffff00", golor.RGBA(255, 255, 255, 0).String())
}

func TestNearestNamedPerceptual(t *testing.T) {
	// RGB distance selects darkslategray; darkslateblue is perceptually much closer.
	name, _ := golor.NearestNamed(golor.RGB(0x44, 0x3b, 0x64))
	assert.Equal(t, "darkslateblue", name)
	name, _ = golor.NearestNamed(golor.RGB(128, 128, 128))
	assert.Equal(t, "gray", name) // The gray/grey tie is resolved by name.
	name, _ = golor.NearestNamed(golor.RGBA(0, 0, 0, 0))
	assert.Equal(t, "black", name) // Transparent is never returned.
	for _, n := range []string{"tomato", "rebeccapurple", "white", "teal"} {
		c, err := golor.Named(n)
		require.NoError(t, err)
		got, _ := golor.NearestNamed(c)
		assert.Equal(t, n, got)
	}
}

// genericImage hides the concrete type so ImagePixels takes its img.At fallback.
type genericImage struct{ image.Image }

func TestImagePixelsFastPathsMatchFromStdColor(t *testing.T) {
	nrgba := image.NewNRGBA(image.Rect(-3, 5, 253, 261))
	rgba := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for a := 0; a < 256; a++ {
		for v := 0; v < 256; v++ {
			x, y := v, a
			nrgba.SetNRGBA(x-3, y+5, stdcolor.NRGBA{uint8(v), uint8(255 - v), uint8(v / 3), uint8(a)})
			// Include invalid premultiplied values (channel > alpha).
			rgba.SetRGBA(x, y, stdcolor.RGBA{uint8(v), uint8(v / 2), uint8(255 - v), uint8(a)})
		}
	}
	for _, img := range []image.Image{
		nrgba, rgba,
		nrgba.SubImage(image.Rect(10, 20, 50, 30)),
		rgba.SubImage(image.Rect(100, 3, 101, 200)),
	} {
		fast := golor.ImagePixels(img)
		slow := golor.ImagePixels(genericImage{img})
		require.Equal(t, len(slow), len(fast))
		for i := range fast {
			if fast[i] != slow[i] {
				t.Fatalf("%T pixel %d: fast %+v, generic %+v", img, i, fast[i], slow[i])
			}
		}
	}
}

func TestFromStdColorStraightAlphaIsLossless(t *testing.T) {
	for a := 1; a < 256; a++ {
		for v := 0; v < 256; v++ {
			src := stdcolor.NRGBA{uint8(v), uint8(255 - v), uint8(v / 2), uint8(a)}
			if got := golor.FromStdColor(src).NRGBA(); got != src {
				t.Fatalf("%v -> %v", src, got)
			}
		}
	}
	c := golor.RGBAf(0.123456789, 0.5, 1, 0.001)
	assert.Equal(t, c, golor.FromStdColor(c))
	assert.Equal(t, golor.Color{}, golor.FromStdColor(stdcolor.NRGBA{R: 9, A: 0}))
	assert.Equal(t, golor.RGBA(255, 0, 0, 128), golor.FromStdColor(stdcolor.NRGBA64{R: 0xffff, A: 0x8080}))
}

func TestPaintAndFill(t *testing.T) {
	img := image.NewNRGBA(image.Rect(2, 2, 5, 4))
	golor.Fill(img, golor.RGBA(1, 2, 3, 4))
	assert.Equal(t, stdcolor.NRGBA{1, 2, 3, 4}, img.NRGBAAt(4, 3))
	golor.PaintPixels(img, []golor.Color{golor.RGB(9, 9, 9), golor.RGB(8, 8, 8), golor.RGB(7, 7, 7), golor.RGB(6, 6, 6)})
	assert.Equal(t, stdcolor.NRGBA{9, 9, 9, 255}, img.NRGBAAt(2, 2))
	assert.Equal(t, stdcolor.NRGBA{6, 6, 6, 255}, img.NRGBAAt(2, 3))
	assert.Equal(t, stdcolor.NRGBA{1, 2, 3, 4}, img.NRGBAAt(3, 3))
	sub := img.SubImage(image.Rect(3, 3, 5, 4)).(*image.NRGBA)
	golor.Fill(sub, golor.RGB(0, 0, 0))
	assert.Equal(t, stdcolor.NRGBA{0, 0, 0, 255}, img.NRGBAAt(4, 3))
	assert.Equal(t, stdcolor.NRGBA{6, 6, 6, 255}, img.NRGBAAt(2, 3))
}

func FuzzParseCSS(f *testing.F) {
	seeds := []string{
		"#fff", "#ff6b35", "#ff6b3580", "tomato", "transparent", "rgb(255 0 0 / 50%)",
		"rgba(255, 0, 0, 0.5)", "hsl(120deg 100% 50%)", "hsla(0.5turn, 10%, 20%, 30%)",
		"hwb(10 20% 30%)", "lab(54.29 80.82 69.91)", "lch(50 150 300)", "oklab(0.5 0.1 -0.1)",
		"oklch(0.7 0.4 145 / none)", "color(display-p3 1 0 0)", "color(xyz-d50 0.9 1 0.8)",
		"rgb(1e308 -1e308 0)", "hsl(1e300 0 0)", "lab(-10 1000 -1000)", "color(srgb-linear -5 5 0.5)",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		c, err := golor.ParseCSS(s)
		if err != nil {
			return
		}
		for _, v := range []float64{c.R, c.G, c.B, c.A} {
			if math.IsNaN(v) || v < 0 || v > 1 {
				t.Fatalf("ParseCSS(%q) = %+v out of range", s, c)
			}
		}
		for _, emit := range formats {
			out := emit.fn(c)
			back, err := golor.ParseCSS(out)
			if err != nil || back.String() != c.String() {
				t.Fatalf("ParseCSS(%q) = %s; %s %q re-parsed as %s (%v)", s, c, emit.name, out, back, err)
			}
		}
	})
}

func FuzzHex(f *testing.F) {
	for _, s := range []string{"#fff", "ffff", "#ff6b35", "#ff6b3580", "  #ABC ", "#", "#zz"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		c, err := golor.Hex(s)
		if err != nil {
			return
		}
		h := strings.TrimPrefix(strings.TrimSpace(s), "#")
		if n := len(h); n != 3 && n != 4 && n != 6 && n != 8 {
			t.Fatalf("Hex(%q) accepted length %d", s, n)
		}
		back, err := golor.Hex(c.String())
		if err != nil || back != c {
			t.Fatalf("Hex(%q) = %s, re-parsed %s (%v)", s, c, back, err)
		}
	})
}
