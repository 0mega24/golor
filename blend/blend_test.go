package blend_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/blend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	black = golor.RGB(0, 0, 0)
	white = golor.RGB(255, 255, 255)
	red   = golor.RGB(255, 0, 0)
	green = golor.RGB(0, 255, 0)
	blue  = golor.RGB(0, 0, 255)
	gray  = golor.RGBf(0.5, 0.5, 0.5)
)

type blendFunc func(golor.Color, golor.Color) golor.Color

func TestMix(t *testing.T) {
	result := blend.Mix(black, white, 0.5)
	assertColor(t, golor.RGBAf(0.5, 0.5, 0.5, 1), result, 1e-6)

	assert.Equal(t, black, blend.Mix(black, white, 0))
	assert.Equal(t, white, blend.Mix(black, white, 1))
	assert.Equal(t, black, blend.Mix(black, white, -1))
	assert.Equal(t, white, blend.Mix(black, white, 2))

	alpha := blend.Mix(golor.RGBA(255, 0, 0, 0), golor.RGBA(0, 0, 255, 255), 0.5)
	assert.Equal(t, uint8(128), alpha.A8())
}

func TestBlendOpaqueReferenceValues(t *testing.T) {
	base := golor.RGBf(0.25, 0.5, 0.75)
	layer := golor.RGBf(0.75, 0.5, 0.25)
	cases := []struct {
		name string
		fn   blendFunc
		want golor.Color
	}{
		{"Multiply", blend.Multiply, golor.RGBf(0.1875, 0.25, 0.1875)},
		{"Screen", blend.Screen, golor.RGBf(0.8125, 0.75, 0.8125)},
		{"Overlay", blend.Overlay, golor.RGBf(0.375, 0.5, 0.625)},
		{"HardLight", blend.HardLight, golor.RGBf(0.625, 0.5, 0.375)},
		{"Difference", blend.Difference, golor.RGBf(0.5, 0, 0.5)},
		{"Exclusion", blend.Exclusion, golor.RGBf(0.625, 0.5, 0.625)},
		{"ColorDodge", blend.ColorDodge, golor.RGBf(1, 1, 1)},
		{"ColorBurn", blend.ColorBurn, golor.RGBf(0, 0, 0)},
		{"VividLight", blend.VividLight, golor.RGBf(0.5, 0.5, 0.5)},
		{"LinearLight", blend.LinearLight, golor.RGBf(0.75, 0.5, 0.25)},
		{"HardMix", blend.HardMix, golor.RGBf(1, 1, 1)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertColor(t, tc.want, tc.fn(base, layer), 1e-9)
		})
	}
}

func TestSoftLightReferenceValues(t *testing.T) {
	assertColor(t, red, blend.SoftLight(red, gray), 0.01)
	assertColor(t, golor.RGBf(0.25, 0.25, 0.25), blend.SoftLight(gray, black), 1e-9)
	assertColor(t, golor.RGBf(math.Sqrt(0.5), math.Sqrt(0.5), math.Sqrt(0.5)), blend.SoftLight(gray, white), 1e-9)
	assertColor(t, golor.RGBf(0.448, 0.448, 0.448), blend.SoftLight(golor.RGBf(0.2, 0.2, 0.2), white), 1e-9)
}

func TestDodgeBurnExtremeBehavior(t *testing.T) {
	assertColor(t, black, blend.ColorDodge(black, white), 1e-9)
	assertColor(t, white, blend.ColorDodge(white, black), 1e-9)
	assertColor(t, white, blend.ColorDodge(gray, white), 1e-9)

	assertColor(t, white, blend.ColorBurn(white, black), 1e-9)
	assertColor(t, black, blend.ColorBurn(black, white), 1e-9)
	assertColor(t, black, blend.ColorBurn(gray, black), 1e-9)
}

// TestNonSeparableBlendModesW3C checks W3C Compositing Level 1 modes (Rec. 601 luma SetLum/SetSat),
// matching CSS mix-blend-mode. Expected values are calculated from the specification formulas.
func TestNonSeparableBlendModesW3C(t *testing.T) {
	// Hue(red, blue) = SetLum(SetSat(blue, 1), Lum(red)=0.3); ClipColor pulls B back to 1.
	hue := blend.Hue(red, blue)
	k := 0.7 / 0.89
	assertColor(t, golor.RGBf(0.3-0.11*k, 0.3-0.11*k, 1), hue, 1e-9)

	saturation := blend.Saturation(golor.RGBf(0.25, 0.5, 0.75), gray)
	assert.InDelta(t, saturation.R, saturation.G, 1e-6)
	assert.InDelta(t, saturation.G, saturation.B, 1e-6)

	// Color(0.2 grey, green) = SetLum(green, 0.2) = clip((-0.39, 0.61, -0.39)).
	color := blend.Color(golor.RGBf(0.2, 0.2, 0.2), green)
	assertColor(t, golor.RGBf(0, 0.2+0.41*0.2/0.59, 0), color, 1e-9)

	luminosity := blend.Luminosity(red, white)
	assertColor(t, white, luminosity, 1e-9)

	// Regression case: HSL produces #3769cd; W3C produces #a4c3ff.
	assert.Equal(t, "#a4c3ff", blend.Luminosity(golor.RGB(0x32, 0x64, 0xc8), golor.RGB(0xf0, 0xc8, 0x14)).String())
	assert.Equal(t, "#7b6400", blend.Color(golor.RGB(0x32, 0x64, 0xc8), golor.RGB(0xf0, 0xc8, 0x14)).String())

	// Every non-separable mode preserves the specified luminosity.
	lum := func(c golor.Color) float64 { return 0.3*c.R + 0.59*c.G + 0.11*c.B }
	b, l := golor.RGBf(0.3, 0.6, 0.2), golor.RGBf(0.9, 0.1, 0.5)
	assert.InDelta(t, lum(l), lum(blend.Luminosity(b, l)), 1e-9)
	assert.InDelta(t, lum(b), lum(blend.Color(b, l)), 1e-9)
	assert.InDelta(t, lum(b), lum(blend.Hue(b, l)), 1e-9)
	assert.InDelta(t, lum(b), lum(blend.Saturation(b, l)), 1e-9)
}

func TestPorterDuffAlphaScenarios(t *testing.T) {
	base := golor.RGBA(128, 64, 32, 128)
	layer := golor.RGBA(64, 128, 255, 0)
	for name, fn := range allPorterDuffBlendModes() {
		t.Run(name+"/transparent_layer", func(t *testing.T) {
			assert.Equal(t, base, fn(base, layer))
		})
	}

	layer = golor.RGBA(64, 128, 255, 128)
	screen := blend.Screen(base, layer)
	opaqueScreen := blend.Screen(golor.RGBA(base.R8(), base.G8(), base.B8(), 255), golor.RGBA(layer.R8(), layer.G8(), layer.B8(), 255))
	ab, as := base.A, layer.A
	ao := as + ab*(1-as)
	assert.InDelta(t, ao, screen.A, 1e-9)
	w3c := func(cb, cs, b float64) float64 { return (as*((1-ab)*cs+ab*b) + ab*cb*(1-as)) / ao }
	assert.InDelta(t, w3c(base.R, layer.R, opaqueScreen.R), screen.R, 1e-9)
	assert.InDelta(t, w3c(base.G, layer.G, opaqueScreen.G), screen.G, 1e-9)
	assert.InDelta(t, w3c(base.B, layer.B, opaqueScreen.B), screen.B, 1e-9)
}

func TestBlendOverTransparentBackdropReturnsLayer(t *testing.T) {
	clear := golor.RGBA(0, 0, 0, 0)
	halfRed := golor.RGBA(255, 0, 0, 128)
	for name, fn := range allPorterDuffBlendModes() {
		t.Run(name, func(t *testing.T) {
			assertColor(t, red, fn(clear, red), 1e-9)
			assert.Equal(t, "#ff000080", fn(clear, halfRed).String())
		})
	}
}

func TestBlendOpaqueUnchangedByAlphaModel(t *testing.T) {
	// For opaque layers over opaque backdrops, the W3C model reduces to B(Cb, Cs).
	base, layer := golor.RGBf(0.25, 0.5, 0.75), golor.RGBf(0.75, 0.5, 0.25)
	assertColor(t, golor.RGBf(0.1875, 0.25, 0.1875), blend.Multiply(base, layer), 1e-12)
	// Opaque backdrop, half-transparent layer: base + αs·(B - base).
	half := golor.RGBAf(0.75, 0.5, 0.25, 0.5)
	assertColor(t, golor.RGBf(0.25+0.5*(0.1875-0.25), 0.5+0.5*(0.25-0.5), 0.75+0.5*(0.1875-0.75)), blend.Multiply(base, half), 1e-12)
}

func TestMixPremultiplied(t *testing.T) {
	clear := golor.RGBA(0, 0, 0, 0)
	assert.Equal(t, "#ffffff80", blend.Mix(clear, white, 0.5).String())
	assert.Equal(t, "#ff000080", blend.Mix(golor.RGBA(0, 0, 255, 0), red, 0.5).String())
	assert.Equal(t, golor.RGBAf(0, 0, 0, 0), blend.Mix(clear, golor.RGBA(255, 255, 255, 0), 0.5))
	assert.Equal(t, black, blend.Mix(black, white, math.NaN()))
}

func TestNoBlendModeLeaksInvalidChannels(t *testing.T) {
	base := golor.RGBAf(0, 1, 0.5, 0.4)
	layer := golor.RGBAf(1, 0, 0.5, 0.6)
	for name, fn := range allPorterDuffBlendModes() {
		t.Run(name, func(t *testing.T) {
			got := fn(base, layer)
			require.False(t, math.IsNaN(got.R))
			require.False(t, math.IsNaN(got.G))
			require.False(t, math.IsNaN(got.B))
			require.False(t, math.IsInf(got.R, 0))
			require.False(t, math.IsInf(got.G, 0))
			require.False(t, math.IsInf(got.B, 0))
			assert.GreaterOrEqual(t, got.R, 0.0)
			assert.LessOrEqual(t, got.R, 1.0)
			assert.GreaterOrEqual(t, got.G, 0.0)
			assert.LessOrEqual(t, got.G, 1.0)
			assert.GreaterOrEqual(t, got.B, 0.0)
			assert.LessOrEqual(t, got.B, 1.0)
		})
	}
}

func allPorterDuffBlendModes() map[string]blendFunc {
	return map[string]blendFunc{
		"Color":       blend.Color,
		"ColorBurn":   blend.ColorBurn,
		"ColorDodge":  blend.ColorDodge,
		"Difference":  blend.Difference,
		"Exclusion":   blend.Exclusion,
		"HardLight":   blend.HardLight,
		"HardMix":     blend.HardMix,
		"Hue":         blend.Hue,
		"LinearLight": blend.LinearLight,
		"Luminosity":  blend.Luminosity,
		"Multiply":    blend.Multiply,
		"Overlay":     blend.Overlay,
		"Saturation":  blend.Saturation,
		"Screen":      blend.Screen,
		"SoftLight":   blend.SoftLight,
		"VividLight":  blend.VividLight,
	}
}

func assertColor(t *testing.T, want, got golor.Color, delta float64) {
	t.Helper()
	assert.InDelta(t, want.R, got.R, delta)
	assert.InDelta(t, want.G, got.G, delta)
	assert.InDelta(t, want.B, got.B, delta)
	assert.InDelta(t, want.A, got.A, delta)
}
