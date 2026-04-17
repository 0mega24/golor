package contrast_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/contrast"
)

func TestLuminance(t *testing.T) {
	cases := []struct {
		name string
		c    golor.Color
		want float64
		tol  float64
	}{
		{"white", golor.RGB(255, 255, 255), 1.0, 1e-6},
		{"black", golor.RGB(0, 0, 0), 0.0, 1e-6},
		{"mid-gray", golor.RGB(128, 128, 128), 0.2159, 0.001},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := contrast.Luminance(tc.c)
			if math.Abs(got-tc.want) > tc.tol {
				t.Errorf("Luminance(%v) = %f, want %f (±%f)", tc.c, got, tc.want, tc.tol)
			}
		})
	}
}

func TestContrastRatio(t *testing.T) {
	white := golor.RGB(255, 255, 255)
	black := golor.RGB(0, 0, 0)
	ratio := contrast.Ratio(white, black)
	if math.Abs(ratio-21.0) > 0.01 {
		t.Errorf("ContrastRatio(white, black) = %f, want 21.0 (±0.01)", ratio)
	}
	same := contrast.Ratio(white, white)
	if math.Abs(same-1.0) > 1e-6 {
		t.Errorf("ContrastRatio(white, white) = %f, want 1.0", same)
	}
}

func TestMeetsContrast(t *testing.T) {
	white := golor.RGB(255, 255, 255)
	black := golor.RGB(0, 0, 0)
	if !contrast.Meets(white, black, 4.5) {
		t.Error("white/black should meet 4.5:1 ratio")
	}
	if contrast.Meets(white, white, 1.1) {
		t.Error("white/white should not meet 1.1:1 ratio")
	}
}

func TestEnforceContrast(t *testing.T) {
	bg := golor.RGB(0, 0, 0)
	fg := golor.RGB(64, 64, 64)
	adjusted := contrast.EnforceContrast(fg, bg, 4.5)
	if !contrast.Meets(adjusted, bg, 4.5) {
		t.Errorf("EnforceContrast result %v does not meet 4.5:1 against black", adjusted)
	}
}
