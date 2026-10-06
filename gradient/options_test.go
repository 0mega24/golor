package gradient_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/gradient"
)

func TestConfiguredDefaultsMatchExisting(t *testing.T) {
	a, b := golor.RGBA(15, 200, 80, 90), golor.RGB(230, 10, 120)
	cases := []struct {
		space gradient.Space
		fn    gradFn
		multi func([]golor.Color, int) []golor.Color
	}{
		{gradient.SpaceRGB, gradient.RGB, gradient.MultiStopRGB},
		{gradient.SpaceLinearRGB, gradient.LinearRGB, gradient.MultiStopLinearRGB},
		{gradient.SpaceHSL, gradient.HSL, gradient.MultiStopHSL},
		{gradient.SpaceHSV, gradient.HSV, gradient.MultiStopHSV},
		{gradient.SpaceLAB, gradient.LAB, gradient.MultiStopLAB},
		{gradient.SpaceLCH, gradient.LCH, gradient.MultiStopLCH},
		{gradient.SpaceOKLab, gradient.OKLab, gradient.MultiStopOKLab},
		{gradient.SpaceOKLCH, gradient.OKLCH, gradient.MultiStopOKLCH},
	}
	for _, tc := range cases {
		got, err := gradient.Between(a, b, 17, gradient.Options{Space: tc.space})
		if err != nil {
			t.Fatal(err)
		}
		want := tc.fn(a, b, 17)
		for i := range got {
			if got[i].String() != want[i].String() {
				t.Fatalf("space %d sample %d: %s != %s", tc.space, i, got[i], want[i])
			}
		}
		mid := golor.RGB(50, 120, 230)
		got, err = gradient.Generate([]gradient.Stop{{Color: a, Position: 0}, {Color: mid, Position: 0.5}, {Color: b, Position: 1}}, 17, gradient.Options{Space: tc.space})
		if err != nil {
			t.Fatal(err)
		}
		want = tc.multi([]golor.Color{a, mid, b}, 17)
		for i := range got {
			if got[i].String() != want[i].String() {
				t.Fatalf("space %d multi sample %d mismatch", tc.space, i)
			}
		}
	}
}

func TestHueRoutes(t *testing.T) {
	for _, sp := range []gradient.Space{gradient.SpaceHSL, gradient.SpaceHSV, gradient.SpaceLCH, gradient.SpaceOKLCH} {
		// Low chroma keeps the entire route inside sRGB in perceptual spaces.
		var makeColor func(float64) golor.Color
		var hue func(golor.Color) float64
		switch sp {
		case gradient.SpaceRGB, gradient.SpaceLinearRGB, gradient.SpaceLAB, gradient.SpaceOKLab:
			t.Fatal("hue test requires a polar space")
		case gradient.SpaceHSL:
			makeColor = func(h float64) golor.Color { return convert.FromHSL(convert.HSL{H: h, S: 0.5, L: 0.5}) }
			hue = func(c golor.Color) float64 { return convert.ToHSL(c).H }
		case gradient.SpaceHSV:
			makeColor = func(h float64) golor.Color { return convert.FromHSV(convert.HSV{H: h, S: 0.5, V: 0.5}) }
			hue = func(c golor.Color) float64 { return convert.ToHSV(c).H }
		case gradient.SpaceLCH:
			makeColor = func(h float64) golor.Color { return convert.FromLCH(convert.LCH{L: 60, C: 10, H: h}) }
			hue = func(c golor.Color) float64 { return convert.ToLCH(c).H }
		case gradient.SpaceOKLCH:
			makeColor = func(h float64) golor.Color { return convert.FromOKLCH(convert.OKLCH{L: 0.6, C: 0.03, H: h}) }
			hue = func(c golor.Color) float64 { return convert.ToOKLCH(c).H }
		}
		for _, tc := range []struct {
			from, to float64
			route    gradient.HueRoute
			want     float64
		}{
			{350, 10, gradient.HueShorter, 0},
			{350, 10, gradient.HueLonger, 180},
			{350, 10, gradient.HueIncreasing, 0},
			{350, 10, gradient.HueDecreasing, 180},
			{10, 350, gradient.HueIncreasing, 180},
			{10, 350, gradient.HueDecreasing, 0},
			{30, 30, gradient.HueLonger, 210},
			{30, 30, gradient.HueShorter, 30},
		} {
			got, err := gradient.Between(makeColor(tc.from), makeColor(tc.to), 3, gradient.Options{Space: sp, Hue: tc.route})
			if err != nil {
				t.Fatal(err)
			}
			if h := hue(got[1]); hueDist(h, tc.want) > 1e-3 {
				t.Fatalf("space %d route %d %g->%g midpoint %g want %g", sp, tc.route, tc.from, tc.to, h, tc.want)
			}
		}
	}
}

func TestLinearLightAndHSV(t *testing.T) {
	black, white := golor.RGB(0, 0, 0), golor.RGB(255, 255, 255)
	if got := gradient.LinearRGB(black, white, 3)[1]; got.R8() != 188 {
		t.Fatalf("linear-light midpoint %s want #bcbcbc", got)
	}
	if got := gradient.HSV(golor.RGB(255, 255, 0), golor.RGB(0, 0, 255), 3)[1]; got.String() != "#00ff80" {
		t.Fatalf("HSV yellow->blue midpoint %s", got)
	}
}

func TestPositionedStopsAndEasing(t *testing.T) {
	black, red, white := golor.RGB(0, 0, 0), golor.RGB(255, 0, 0), golor.RGB(255, 255, 255)
	stops := []gradient.Stop{{Color: black, Position: 0}, {Color: red, Position: 0.25}, {Color: white, Position: 1}}
	original := append([]gradient.Stop(nil), stops...)
	got, err := gradient.Generate(stops, 9, gradient.Options{Easing: gradient.EaseIn})
	if err != nil {
		t.Fatal(err)
	}
	if got[1].R != 0.25 || got[2] != red || got[5].G != 0.25 || got[8] != white {
		t.Fatalf("positioned/eased samples: %v", got)
	}
	if !reflect.DeepEqual(stops, original) {
		t.Fatal("input stops changed")
	}
	hard := []gradient.Stop{{Color: black, Position: 0.25}, {Color: red, Position: 0.5}, {Color: white, Position: 0.5}}
	got, err = gradient.Generate(hard, 5, gradient.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if got[0] != black || got[1] != black || got[2] != white || got[4] != white {
		t.Fatalf("hard transition/held edges: %v", got)
	}
	got, err = gradient.Generate([]gradient.Stop{{Color: red, Position: 0.5}}, 4, gradient.Options{})
	if err != nil || len(got) != 4 {
		t.Fatalf("one stop: %v %v", got, err)
	}
	for _, c := range got {
		if c != red {
			t.Fatal("one stop should fill every sample")
		}
	}
}

func TestCurvesAndSafeCallbacks(t *testing.T) {
	black, white := golor.RGB(0, 0, 0), golor.RGB(255, 255, 255)
	for _, tc := range []struct {
		curve gradient.Curve
		want  float64
	}{
		{gradient.Linear, 0.5},
		{gradient.EaseIn, 0.25},
		{gradient.EaseOut, 0.75},
		{gradient.EaseInOut, 0.5},
		{gradient.Logarithmic(3), math.Log(2.5) / math.Log(4)},
		{gradient.Logarithmic(-0.5), math.Log(0.75) / math.Log(0.5)},
		{gradient.Logarithmic(0), 0.5},
		{gradient.Logarithmic(-1), 0.5},
		{gradient.Logarithmic(math.Inf(1)), 0.5},
		{func(float64) float64 { return math.NaN() }, 0.5},
		{func(float64) float64 { return 2 }, 1},
		{func(float64) float64 { return -1 }, 0},
	} {
		got, err := gradient.Between(black, white, 3, gradient.Options{Easing: tc.curve})
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(got[1].R-tc.want) > 1e-12 || got[0] != black || got[2] != white {
			t.Fatalf("curve samples %v want midpoint %g and exact endpoints", got, tc.want)
		}
	}
	clear := golor.RGBA(0, 0, 0, 0)
	got, err := gradient.Between(clear, white, 3, gradient.Options{Space: gradient.SpaceHSV, Easing: gradient.EaseIn})
	if err != nil || got[1].R != 1 || got[1].A != 0.25 {
		t.Fatalf("eased premultiplied fade %v %v", got, err)
	}
}

func TestChromaControls(t *testing.T) {
	for _, sp := range []gradient.Space{gradient.SpaceLCH, gradient.SpaceOKLCH} {
		scale := 1.0
		var a, b golor.Color
		var chroma func(golor.Color) float64
		if sp == gradient.SpaceLCH {
			scale = 100
			a = convert.FromLCH(convert.LCH{L: 60, C: 3, H: 40})
			b = convert.FromLCH(convert.LCH{L: 60, C: 3, H: 80})
			chroma = func(c golor.Color) float64 { return convert.ToLCH(c).C }
		} else {
			a = convert.FromOKLCH(convert.OKLCH{L: 0.6, C: 0.03, H: 40})
			b = convert.FromOKLCH(convert.OKLCH{L: 0.6, C: 0.03, H: 80})
			chroma = func(c golor.Color) float64 { return convert.ToOKLCH(c).C }
		}
		got, err := gradient.Between(a, b, 3, gradient.Options{Space: sp, MinChroma: 0.05 * scale, Chroma: func(t, c float64) float64 { return c + 0.01*scale }})
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(chroma(got[1])-0.05*scale) > 1e-5 || got[0] != a || got[2] != b {
			t.Fatalf("space %d chroma control %v", sp, got)
		}
		got, err = gradient.Between(a, b, 3, gradient.Options{Space: sp, Chroma: func(t, c float64) float64 { return -1 }})
		if err != nil || chroma(got[1]) > 1e-5 {
			t.Fatalf("negative chroma: %v %v", got, err)
		}
		got, err = gradient.Between(a, b, 3, gradient.Options{Space: sp, MinChroma: 1000})
		if err != nil {
			t.Fatal(err)
		}
		for _, c := range got {
			for _, v := range []float64{c.R, c.G, c.B, c.A} {
				if math.IsNaN(v) || v < 0 || v > 1 {
					t.Fatalf("gamut mapping failed: %v", c)
				}
			}
		}
		got, err = gradient.Between(a, b, 3, gradient.Options{Space: sp, Chroma: func(t, c float64) float64 { return math.NaN() }})
		if err != nil || math.Abs(chroma(got[1])-0.03*scale) > 1e-5 {
			t.Fatalf("non-finite chroma: %v %v", got, err)
		}
	}
}

func TestGenerateValidationAndSizes(t *testing.T) {
	for _, opts := range []gradient.Options{
		{Space: 255},
		{Hue: 255},
		{MinChroma: -1},
		{Space: gradient.SpaceOKLCH, MinChroma: math.NaN()},
		{Space: gradient.SpaceOKLCH, MinChroma: math.Inf(1)},
		{MinChroma: 0.1},
		{Chroma: func(t, c float64) float64 { return c }},
	} {
		if _, err := gradient.Generate(nil, 0, opts); err == nil {
			t.Fatalf("accepted invalid options %+v", opts)
		}
	}
	for _, stops := range [][]gradient.Stop{
		{{Position: -0.1}}, {{Position: 1.1}}, {{Position: math.NaN()}}, {{Position: math.Inf(1)}}, {{Position: 0.8}, {Position: 0.2}},
	} {
		if _, err := gradient.Generate(stops, 0, gradient.Options{}); err == nil {
			t.Fatalf("accepted invalid stops %v", stops)
		}
	}
	for _, n := range []int{-1, 0, 1, 5} {
		got, err := gradient.Generate(nil, n, gradient.Options{})
		if err != nil || got != nil {
			t.Fatalf("empty stops %v %v", got, err)
		}
		got, err = gradient.Between(golor.RGB(255, 0, 0), golor.RGB(0, 0, 255), n, gradient.Options{})
		if err != nil || len(got) != max(n, 0) {
			t.Fatalf("sample count %d: %v %v", n, got, err)
		}
	}
}
