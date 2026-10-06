package gradient

import (
	"math"
	"testing"
)

func TestHueRouteTiesAndWrapping(t *testing.T) {
	for _, tc := range []struct {
		a, b  float64
		route HueRoute
		want  float64
	}{
		{0, 180, HueShorter, 90},
		{180, 0, HueShorter, 90},
		{0, 180, HueLonger, 90},
		{180, 0, HueLonger, 90},
		{0, 180, HueIncreasing, 90},
		{180, 0, HueIncreasing, 270},
		{0, 180, HueDecreasing, 270},
		{180, 0, HueDecreasing, 90},
		{30, 30, HueShorter, 30},
		{30, 30, HueLonger, 210},
		{30, 30, HueIncreasing, 30},
		{30, 30, HueDecreasing, 30},
		{710, -350, HueShorter, 0},
		{710, -350, HueLonger, 180},
	} {
		if got := lerpHueRoute(tc.a, tc.b, 0.5, tc.route); math.Abs(got-tc.want) > 1e-12 {
			t.Errorf("%g->%g route %d got %g want %g", tc.a, tc.b, tc.route, got, tc.want)
		}
	}
}
