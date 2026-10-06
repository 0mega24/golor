package adjust_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/adjust"
	"github.com/0mega24/golor/v2/convert"
	"github.com/stretchr/testify/assert"
)

func TestLighten(t *testing.T) {
	c := golor.RGB(100, 100, 100)
	before := convert.ToHSL(c).L
	after := convert.ToHSL(adjust.Lighten(c, 0.2)).L
	assert.InDelta(t, before+0.2, after, 1e-6)
}

func TestDarken(t *testing.T) {
	c := golor.RGB(200, 200, 200)
	before := convert.ToHSL(c).L
	after := convert.ToHSL(adjust.Darken(c, 0.2)).L
	assert.InDelta(t, before-0.2, after, 1e-6)
}

func TestSetLightness(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	result := adjust.SetLightness(c, 0.7)
	assert.InDelta(t, 0.7, convert.ToHSL(result).L, 1e-6)
}

func TestLightenClamps(t *testing.T) {
	white := golor.RGB(255, 255, 255)
	result := adjust.Lighten(white, 0.5)
	assert.InDelta(t, 1.0, convert.ToHSL(result).L, 1e-6)
}

func TestSaturate(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	before := convert.ToHSL(c).S
	result := adjust.Saturate(c, 0.2)
	assert.InDelta(t, min01(before+0.2), convert.ToHSL(result).S, 1e-6)
}

func TestDesaturate(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	before := convert.ToHSL(c).S
	result := adjust.Desaturate(c, 0.1)
	assert.InDelta(t, max0(before-0.1), convert.ToHSL(result).S, 1e-6)
}

func TestSetSaturation(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	result := adjust.SetSaturation(c, 0.5)
	assert.InDelta(t, 0.5, convert.ToHSL(result).S, 1e-6)
}

func TestShiftHue(t *testing.T) {
	c := golor.RGB(255, 0, 0) // hue ~0°
	result := adjust.ShiftHue(c, 180)
	hue := convert.ToHSL(result).H
	assert.InDelta(t, 180.0, hue, 1.0)
}

func TestSetHue(t *testing.T) {
	c := golor.RGB(255, 0, 0)
	result := adjust.SetHue(c, 120)
	assert.InDelta(t, 120.0, convert.ToHSL(result).H, 1.0)
}

func TestSetValue(t *testing.T) {
	c := golor.RGB(200, 100, 50)
	result := adjust.SetValue(c, 0.3)
	assert.InDelta(t, 0.3, convert.ToHSV(result).V, 1e-6)
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
	before := convert.ToHSL(blue).H
	result := adjust.Warm(blue, 0.5)
	after := convert.ToHSL(result).H
	// Hue should have moved toward 30°
	target := 30.0
	distBefore := hueDist(before, target)
	distAfter := hueDist(after, target)
	assert.Less(t, distAfter, distBefore, "warm should move hue closer to orange")
}

func TestCool(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	before := convert.ToHSL(red).H
	result := adjust.Cool(red, 0.5)
	after := convert.ToHSL(result).H
	target := 210.0
	distBefore := hueDist(before, target)
	distAfter := hueDist(after, target)
	assert.Less(t, distAfter, distBefore, "cool should move hue closer to blue")
}

func TestAdjustPreservesAlpha(t *testing.T) {
	c := golor.RGBA(100, 150, 200, 77)
	results := []golor.Color{
		adjust.Lighten(c, 0.1),
		adjust.Darken(c, 0.1),
		adjust.Saturate(c, 0.1),
		adjust.Desaturate(c, 0.1),
		adjust.ShiftHue(c, 30),
		adjust.SetHue(c, 120),
		adjust.SetLightness(c, 0.5),
		adjust.SetSaturation(c, 0.5),
		adjust.SetValue(c, 0.5),
		adjust.Tint(c, 0.5),
		adjust.Shade(c, 0.5),
		adjust.Warm(c, 0.5),
		adjust.Cool(c, 0.5),
	}
	for _, result := range results {
		assert.Equal(t, c.A8(), result.A8())
	}
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
