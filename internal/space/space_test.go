package space

import (
	"math"
	"testing"
)

func near(t *testing.T, name string, got, want Vec3, tol float64) {
	t.Helper()
	for i := range got {
		if math.Abs(got[i]-want[i]) > tol {
			t.Fatalf("%s = %v, want %v (tol %g)", name, got, want, tol)
		}
	}
}

func TestMatricesMapWhiteToWhite(t *testing.T) {
	near(t, "sRGB white", SRGBToXYZ.Mul(Vec3{1, 1, 1}), D65, 1e-15)
	near(t, "P3 white", P3ToXYZ.Mul(Vec3{1, 1, 1}), D65, 1e-15)
	near(t, "D50->D65", D50ToD65.Mul(D50), D65, 1e-12)
	near(t, "D65->D50", D65ToD50.Mul(D65), D50, 1e-12)
	near(t, "inverse", XYZToSRGB.Mul(SRGBToXYZ.Mul(Vec3{0.2, 0.5, 0.9})), Vec3{0.2, 0.5, 0.9}, 1e-15)
}

func TestGreyHasNoChroma(t *testing.T) {
	for _, v := range []float64{1, 0.5, 0.2, 0.01} {
		lin := Linearize(v)
		lab := XYZToLab(SRGBToXYZ.Mul(Vec3{lin, lin, lin}), D65)
		if math.Hypot(lab[1], lab[2]) > 1e-12 {
			t.Fatalf("grey %v has Lab chroma %g", v, math.Hypot(lab[1], lab[2]))
		}
		ok := LinearSRGBToOKLab(Vec3{lin, lin, lin})
		if math.Hypot(ok[1], ok[2]) > 1e-7 {
			t.Fatalf("grey %v has OKLab chroma %g", v, math.Hypot(ok[1], ok[2]))
		}
	}
}

func TestOKLabReference(t *testing.T) {
	// CSS Color 4 / Ottosson reference values.
	near(t, "red", LinearSRGBToOKLab(Vec3{1, 0, 0}), Vec3{0.627955, 0.224863, 0.125846}, 1e-5)
	near(t, "green", LinearSRGBToOKLab(Vec3{0, 1, 0}), Vec3{0.866440, -0.233888, 0.179498}, 1e-5)
	near(t, "blue", LinearSRGBToOKLab(Vec3{0, 0, 1}), Vec3{0.452014, -0.032457, -0.311528}, 1e-5)
	near(t, "round trip", OKLabToLinearSRGB(LinearSRGBToOKLab(Vec3{0.3, 0.7, 0.1})), Vec3{0.3, 0.7, 0.1}, 1e-12)
}

func TestLinearizeLUTMatchesFormula(t *testing.T) {
	for i := 0; i < 256; i++ {
		v := float64(i) / 255
		if math.Abs(Linearize(v)-linearizeSlow(v)) > 1e-15 {
			t.Fatalf("LUT mismatch at %d", i)
		}
		if math.Abs(Delinearize(Linearize(v))-v) > 1e-15 {
			t.Fatalf("transfer round trip at %d", i)
		}
	}
	if Linearize(-0.5) != -linearizeSlow(0.5) {
		t.Fatal("odd extension")
	}
}

func TestLabRoundTrip(t *testing.T) {
	for _, xyz := range []Vec3{{0.2, 0.3, 0.4}, {0.001, 0.0005, 0.002}, D65, {}} {
		near(t, "lab", LabToXYZ(XYZToLab(xyz, D65), D65), xyz, 1e-14)
	}
}

func TestGamutMap(t *testing.T) {
	// In-gamut colors pass through unchanged.
	near(t, "in gamut", GamutMapLinear(Vec3{0.2, 0.4, 0.6}), Vec3{0.2, 0.4, 0.6}, 1e-15)
	// Highly saturated OKLCH maps into gamut with similar lightness and hue.
	for h := 0.0; h < 360; h += 7.5 {
		for _, l := range []float64{0.2, 0.5, 0.8, 0.95} {
			in := Vec3{l, 0.4, h}
			out := GamutMapOKLCH(in)
			if !InGamutLinear(out) {
				t.Fatalf("not in gamut: %v -> %v", in, out)
			}
			got := ToPolar(LinearSRGBToOKLab(out), AchromaticOKLCH)
			if math.Abs(got[0]-l) > 0.03 {
				t.Fatalf("lightness drift %v -> %v", in, got)
			}
			// The result must be within one JND (0.02 in OKLab) of the chroma-reduced
			// color with the original lightness and hue.
			reduced := FromPolar(Vec3{l, got[1], h})
			if d := Vec3Distance(LinearSRGBToOKLab(out), reduced); d > 0.0201 {
				t.Fatalf("hue/lightness drift %v -> %v (dE %g)", in, got, d)
			}
		}
	}
	near(t, "white", GamutMapOKLCH(Vec3{1.2, 0.3, 40}), Vec3{1, 1, 1}, 0)
	near(t, "black", GamutMapOKLCH(Vec3{-0.1, 0.3, 40}), Vec3{}, 0)
	near(t, "nan", GamutMapLinear(Vec3{math.NaN(), 0, 0}), Vec3{}, 0)
}

func TestNormalizeHue(t *testing.T) {
	cases := map[float64]float64{0: 0, 360: 0, 720: 0, -90: 270, -400: 320, 1080.5: 0.5, -1e-18: 0}
	for in, want := range cases {
		if got := NormalizeHue(in); math.Abs(got-want) > 1e-9 {
			t.Fatalf("NormalizeHue(%v) = %v, want %v", in, got, want)
		}
	}
	if NormalizeHue(math.NaN()) != 0 || Clamp01(math.NaN()) != 0 {
		t.Fatal("NaN handling")
	}
}
