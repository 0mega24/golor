package convert_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

func TestToHSVRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		c    golor.Color
	}{
		{"red", golor.RGB(255, 0, 0)},
		{"green", golor.RGB(0, 255, 0)},
		{"blue", golor.RGB(0, 0, 255)},
		{"white", golor.RGB(255, 255, 255)},
		{"black", golor.RGB(0, 0, 0)},
		{"mid", golor.RGB(128, 64, 200)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := convert.FromHSV(convert.ToHSV(tc.c))
			if math.Abs(got.R-tc.c.R) > 1e-6 || math.Abs(got.G-tc.c.G) > 1e-6 || math.Abs(got.B-tc.c.B) > 1e-6 {
				t.Errorf("HSV round-trip: got %v, want %v", got, tc.c)
			}
		})
	}
}

func TestToHSLRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		c    golor.Color
	}{
		{"red", golor.RGB(255, 0, 0)},
		{"green", golor.RGB(0, 255, 0)},
		{"blue", golor.RGB(0, 0, 255)},
		{"white", golor.RGB(255, 255, 255)},
		{"black", golor.RGB(0, 0, 0)},
		{"mid", golor.RGB(100, 150, 200)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := convert.FromHSL(convert.ToHSL(tc.c))
			if math.Abs(got.R-tc.c.R) > 1e-6 || math.Abs(got.G-tc.c.G) > 1e-6 || math.Abs(got.B-tc.c.B) > 1e-6 {
				t.Errorf("HSL round-trip: got %v, want %v", got, tc.c)
			}
		})
	}
}

func TestToLABRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		c    golor.Color
	}{
		{"red", golor.RGB(255, 0, 0)},
		{"white", golor.RGB(255, 255, 255)},
		{"black", golor.RGB(0, 0, 0)},
		{"mid", golor.RGB(128, 64, 200)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := convert.FromLAB(convert.ToLAB(tc.c))
			if math.Abs(got.R-tc.c.R) > 1e-5 || math.Abs(got.G-tc.c.G) > 1e-5 || math.Abs(got.B-tc.c.B) > 1e-5 {
				t.Errorf("LAB round-trip: got %v, want %v", got, tc.c)
			}
		})
	}
}

func TestToLCHRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		c    golor.Color
	}{
		{"red", golor.RGB(255, 0, 0)},
		{"blue", golor.RGB(0, 0, 255)},
		{"mid", golor.RGB(128, 64, 200)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := convert.FromLCH(convert.ToLCH(tc.c))
			if math.Abs(got.R-tc.c.R) > 1e-5 || math.Abs(got.G-tc.c.G) > 1e-5 || math.Abs(got.B-tc.c.B) > 1e-5 {
				t.Errorf("LCH round-trip: got %v, want %v", got, tc.c)
			}
		})
	}
}

func TestToCMYKRoundTrip(t *testing.T) {
	cases := []golor.Color{
		golor.RGB(255, 0, 0),
		golor.RGB(0, 0, 0),
		golor.RGB(255, 255, 255),
		golor.RGB(128, 64, 200),
	}
	for _, c := range cases {
		got := convert.FromCMYK(convert.ToCMYK(c))
		if math.Abs(got.R-c.R) > 1e-6 || math.Abs(got.G-c.G) > 1e-6 || math.Abs(got.B-c.B) > 1e-6 {
			t.Errorf("CMYK round-trip: got %v, want %v", got, c)
		}
	}
}
