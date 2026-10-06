package contrast_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/contrast"
)

var white = golor.RGB(255, 255, 255)

func hex(t *testing.T, s string) golor.Color {
	t.Helper()
	c, err := golor.Hex(s)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// TestAPCAReferenceValues checks values from the apca-w3 README (APCA-W3 0.0.98G-4g).
func TestAPCAReferenceValues(t *testing.T) {
	cases := []struct {
		text, bg string
		want     float64
	}{
		{"#888888", "#ffffff", 63.056469930209424},
		{"#ffffff", "#888888", -68.54146436644962},
		{"#000000", "#aaaaaa", 58.146262578561334},
		{"#aaaaaa", "#000000", -56.24113336839742},
	}
	for _, tc := range cases {
		got := contrast.APCA(hex(t, tc.text), hex(t, tc.bg))
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("APCA(%s on %s) = %.12f, want %.12f", tc.text, tc.bg, got, tc.want)
		}
	}
	if got := contrast.APCA(white, white); got != 0 {
		t.Errorf("APCA(white, white) = %v, want 0", got)
	}
	// Low-contrast pairs clip to zero.
	if got := contrast.APCA(hex(t, "#fefefe"), white); got != 0 {
		t.Errorf("APCA(#fefefe on white) = %v, want 0", got)
	}
	// Transparent text over a background has no contrast.
	if got := contrast.APCA(golor.RGBA(0, 0, 0, 0), white); got != 0 {
		t.Errorf("APCA(transparent on white) = %v, want 0", got)
	}
}

func TestEnforceContrastProbeRepros(t *testing.T) {
	cases := []struct {
		fg, bg string
		min    float64
	}{
		{"#969696", "#808080", 4.5}, // Mid-tone background: only black can pass.
		{"#969696", "#777777", 4.5},
		{"#1e1ec8", "#767676", 4.5},   // Regression case: the previous result falls just below the threshold.
		{"#f0f0f0", "#00000014", 4.5}, // Nearly transparent black background composites to near-white.
		{"#40404080", "#000000", 3},   // Translucent foreground.
		{"#c8c8c8", "#141414", 18.4},  // Threshold is just reachable.
		{"#336699", "#336699", 4.5},   // Identical colors.
		{"#ff0000", "#00ff00", 3},     // Saturated pair.
	}
	for _, tc := range cases {
		fg, bg := hex(t, tc.fg), hex(t, tc.bg)
		got, ok := contrast.EnforceContrastOK(fg, bg, tc.min)
		if !ok || !contrast.Meets(got, bg, tc.min) {
			t.Errorf("EnforceContrast(%s, %s, %.1f) = %s ratio %.3f ok=%v", tc.fg, tc.bg, tc.min, got, contrast.Ratio(got, bg), ok)
		}
		if back := hex(t, got.String()); !contrast.Meets(back, bg, tc.min) {
			t.Errorf("EnforceContrast(%s, %s) = %s fails after hex round trip", tc.fg, tc.bg, got)
		}
		if got.A8() != fg.A8() {
			t.Errorf("alpha changed: %s -> %s", fg, got)
		}
	}
}

func TestEnforceContrastUnreachable(t *testing.T) {
	bg := hex(t, "#141414")
	got, ok := contrast.EnforceContrastOK(hex(t, "#c8c8c8"), bg, 22)
	if ok {
		t.Fatal("22:1 is unreachable")
	}
	if got.String() != "#ffffff" {
		t.Errorf("unreachable fallback = %s, want the highest-contrast endpoint #ffffff", got)
	}
	// A nearly transparent fg cannot reach 4.5:1 in either direction; alpha is preserved.
	faint := golor.RGBA(0, 0, 0, 40)
	got, ok = contrast.EnforceContrastOK(faint, white, 4.5)
	if ok || got.A8() != 40 {
		t.Errorf("faint fg: got %s ok=%v", got, ok)
	}
	if contrast.EnforceContrast(faint, white, 4.5) != got {
		t.Error("EnforceContrast and EnforceContrastOK disagree")
	}
}

func TestEnforceContrastAlreadyMeetsIsUnchanged(t *testing.T) {
	fg := golor.RGBf(0.123456789, 0.1, 0.1)
	got, ok := contrast.EnforceContrastOK(fg, white, 4.5)
	if !ok || got != fg {
		t.Errorf("got %v ok=%v, want fg unchanged", got, ok)
	}
}

// TestEnforceContrastSweep checks that EnforceContrast passes whenever black or white passes, including
// after rounding to hex.
func TestEnforceContrastSweep(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	const n = 20000
	for i := 0; i < n; i++ {
		fg := golor.RGBA(uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(128+r.Intn(128)))
		bg := golor.RGBA(uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(r.Intn(256)))
		minRatio := []float64{3, 4.5, 7}[i%3]
		fgW := golor.RGBA(255, 255, 255, fg.A8())
		fgK := golor.RGBA(0, 0, 0, fg.A8())
		possible := contrast.Meets(fgW, bg, minRatio) || contrast.Meets(fgK, bg, minRatio) || contrast.Meets(fg, bg, minRatio)
		got, ok := contrast.EnforceContrastOK(fg, bg, minRatio)
		if ok != possible {
			t.Fatalf("fg %s bg %s min %.1f: ok=%v possible=%v (got %s)", fg, bg, minRatio, ok, possible, got)
		}
		if !ok {
			continue
		}
		if !contrast.Meets(got, bg, minRatio) {
			t.Fatalf("fg %s bg %s min %.1f: %s ratio %.4f", fg, bg, minRatio, got, contrast.Ratio(got, bg))
		}
		back, err := golor.Hex(got.String())
		if err != nil || !contrast.Meets(back, bg, minRatio) {
			t.Fatalf("fg %s bg %s: %s fails after hex rounding", fg, bg, got)
		}
	}
}

func TestLuminanceClampsInput(t *testing.T) {
	if l := contrast.Luminance(golor.Color{R: 2, G: 2, B: 2, A: 1}); l != 1 {
		t.Errorf("Luminance(out of range) = %v, want 1", l)
	}
	if l := contrast.Luminance(golor.Color{R: math.NaN(), A: 1}); l != 0 {
		t.Errorf("Luminance(NaN) = %v, want 0", l)
	}
}
