// Package gradient interpolates in RGB, linear-light RGB, HSL/HSV, LAB/LCH, and OKLab/OKLCH.
// Defaults use premultiplied alpha, shorter hue arcs, borrowed achromatic hues, and gamut mapping.
package gradient

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// mixer returns the color at position t in [0, 1] between two fixed endpoints.
type mixer func(t float64) golor.Color

// sample returns n evenly spaced colors from m, including both endpoints.
// n == 1 returns only the start color.
func sample(n int, start golor.Color, m mixer) []golor.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []golor.Color{start}
	}
	out := make([]golor.Color, n)
	for i := range out {
		out[i] = m(float64(i) / float64(n-1))
	}
	return out
}

// multiStop samples n colors through stops, building each segment's mixer once.
func multiStop(stops []golor.Color, n int, newMixer func(a, b golor.Color) mixer) []golor.Color {
	if n <= 0 || len(stops) == 0 {
		return nil
	}
	if len(stops) == 1 || n == 1 {
		return []golor.Color{stops[0]}
	}
	segs := make([]mixer, len(stops)-1)
	for i := range segs {
		segs[i] = newMixer(stops[i], stops[i+1])
	}
	out := make([]golor.Color, n)
	for i := range out {
		scaled := float64(i) / float64(n-1) * float64(len(segs))
		seg := min(int(scaled), len(segs)-1)
		out[i] = segs[seg](scaled - float64(seg))
	}
	return out
}

func alphas(a, b golor.Color) (float64, float64) {
	return space.Clamp01(a.A), space.Clamp01(b.A)
}

// mixPremul interpolates component x premultiplied by alpha and un-premultiplies the
// result. When the interpolated alpha is 0, it uses the straight value.
func mixPremul(x1, x2, a1, a2, alpha, t float64) float64 {
	if a1 == a2 || alpha <= 0 {
		return x1 + t*(x2-x1)
	}
	return (x1*a1 + t*(x2*a2-x1*a1)) / alpha
}

// rect builds a mixer for a rectangular space (RGB, Lab, OKLab).
func rect(a, b golor.Color, to func(golor.Color) space.Vec3, from func(space.Vec3) golor.Color) mixer {
	va, vb := to(a), to(b)
	aa, ab := alphas(a, b)
	return func(t float64) golor.Color {
		alpha := aa + t*(ab-aa)
		var v space.Vec3
		for i := range v {
			v[i] = mixPremul(va[i], vb[i], aa, ab, alpha, t)
		}
		out := from(v)
		out.A = alpha
		return out
	}
}

// polarColor represents a polar color as (lightness, chroma, hue°).
type polarColor struct {
	v          space.Vec3
	achromatic bool
}

// polarWithOptions builds a polar mixer; hue is not premultiplied.
func polarWithOptions(a, b golor.Color, to func(golor.Color) polarColor, from func(space.Vec3) golor.Color, options Options) mixer {
	pa, pb := to(a), to(b)
	switch {
	case pa.achromatic && !pb.achromatic:
		pa.v[2] = pb.v[2]
	case pb.achromatic && !pa.achromatic:
		pb.v[2] = pa.v[2]
	}
	aa, ab := alphas(a, b)
	return func(t float64) golor.Color {
		alpha := aa + t*(ab-aa)
		chroma := mixPremul(pa.v[1], pb.v[1], aa, ab, alpha, t)
		if t > 0 && t < 1 {
			if options.Chroma != nil {
				if adjusted := options.Chroma(t, chroma); finite(adjusted) {
					chroma = math.Max(0, adjusted)
				}
			}
			chroma = math.Max(chroma, options.MinChroma)
		}
		out := from(space.Vec3{
			mixPremul(pa.v[0], pb.v[0], aa, ab, alpha, t),
			chroma,
			lerpHueRoute(pa.v[2], pb.v[2], t, options.Hue),
		})
		out.A = alpha
		return out
	}
}

func lerpHueRoute(h1, h2, t float64, route HueRoute) float64 {
	h1, h2 = space.NormalizeHue(h1), space.NormalizeHue(h2)
	diff := h2 - h1
	switch route {
	case HueLonger:
		switch {
		case diff >= 0 && diff < 180:
			diff -= 360
		case diff > -180 && diff < 0:
			diff += 360
		}
		if diff == -360 {
			diff = 360
		}
	case HueIncreasing:
		if diff < 0 {
			diff += 360
		}
	case HueDecreasing:
		if diff > 0 {
			diff -= 360
		}
	case HueShorter:
		switch {
		case diff > 180:
			diff -= 360
		case diff < -180:
			diff += 360
		}
	}
	return space.NormalizeHue(h1 + t*diff)
}
