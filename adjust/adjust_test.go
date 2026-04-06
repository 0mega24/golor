package adjust_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/adjust"
	"github.com/0mega24/golor/convert"
	"github.com/stretchr/testify/assert"
)

func hslOf(c golor.Color) convert.HSL { return convert.ToHSL(c) }
func hsvOf(c golor.Color) convert.HSV { return convert.ToHSV(c) }

func TestLighten(t *testing.T) {
	c := golor.RGB(100, 100, 100)
	before := hslOf(c).L
	after := hslOf(adjust.Lighten(c, 0.2)).L
	assert.InDelta(t, before+0.2, after, 1e-6)
}

func TestDarken(t *testing.T) {
	c := golor.RGB(200, 200, 200)
	before := hslOf(c).L
	after := hslOf(adjust.Darken(c, 0.2)).L
	assert.InDelta(t, before-0.2, after, 1e-6)
}

func TestSetLightness(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	result := adjust.SetLightness(c, 0.7)
	assert.InDelta(t, 0.7, hslOf(result).L, 1e-6)
}

func TestLightenClamps(t *testing.T) {
	white := golor.RGB(255, 255, 255)
	result := adjust.Lighten(white, 0.5)
	assert.InDelta(t, 1.0, hslOf(result).L, 1e-6)
}

func TestSaturate(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	before := hslOf(c).S
	result := adjust.Saturate(c, 0.2)
	assert.InDelta(t, min01(before+0.2), hslOf(result).S, 1e-6)
}

func TestDesaturate(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	before := hslOf(c).S
	result := adjust.Desaturate(c, 0.1)
	assert.InDelta(t, max0(before-0.1), hslOf(result).S, 1e-6)
}

func TestSetSaturation(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	result := adjust.SetSaturation(c, 0.5)
	assert.InDelta(t, 0.5, hslOf(result).S, 1e-6)
}

func TestShiftHue(t *testing.T) {
	c := golor.RGB(255, 0, 0) // hue ~0°
	result := adjust.ShiftHue(c, 180)
	hue := hslOf(result).H
	assert.InDelta(t, 180.0, hue, 1.0)
}

func TestSetHue(t *testing.T) {
	c := golor.RGB(255, 0, 0)
	result := adjust.SetHue(c, 120)
	assert.InDelta(t, 120.0, hslOf(result).H, 1.0)
}

func TestSetValue(t *testing.T) {
	c := golor.RGB(200, 100, 50)
	result := adjust.SetValue(c, 0.3)
	assert.InDelta(t, 0.3, hsvOf(result).V, 1e-6)
}

func TestTint(t *testing.T) {
	black := golor.RGB(0, 0, 0)
	result := adjust.Tint(black, 1.0)
	assert.InDelta(t, 1.0, result.R, 1e-6)
	assert.InDelta(t, 1.0, result.G, 1e-6)
	assert.InDelta(t, 1.0, result.B, 1e-6)
}

func TestShade(t *testing.T) {
	white := golor.RGB(255, 255, 255)
	result := adjust.Shade(white, 1.0)
	assert.InDelta(t, 0.0, result.R, 1e-6)
	assert.InDelta(t, 0.0, result.G, 1e-6)
	assert.InDelta(t, 0.0, result.B, 1e-6)
}

func TestWarm(t *testing.T) {
	// A blue color; warming should move hue toward orange
	blue := golor.RGB(0, 0, 255)
	before := hslOf(blue).H
	result := adjust.Warm(blue, 0.5)
	after := hslOf(result).H
	// Hue should have moved toward 30°
	target := 30.0
	distBefore := hueDist(before, target)
	distAfter := hueDist(after, target)
	assert.Less(t, distAfter, distBefore, "warm should move hue closer to orange")
}

func TestCool(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	before := hslOf(red).H
	result := adjust.Cool(red, 0.5)
	after := hslOf(result).H
	target := 210.0
	distBefore := hueDist(before, target)
	distAfter := hueDist(after, target)
	assert.Less(t, distAfter, distBefore, "cool should move hue closer to blue")
}

func hueDist(h1, h2 float64) float64 {
	d := math.Abs(h1 - h2)
	if d > 180 {
		d = 360 - d
	}
	return d
}

func min01(v float64) float64 {
	if v > 1 {
		return 1
	}
	return v
}

func max0(v float64) float64 {
	if v < 0 {
		return 0
	}
	return v
}
