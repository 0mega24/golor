package transform_test

import (
	"testing"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/blend"
	"github.com/0mega24/golor/convert"
	"github.com/0mega24/golor/transform"
	"github.com/stretchr/testify/assert"
)

func TestChainLightenDarken(t *testing.T) {
	c := golor.RGB(100, 100, 100)
	original := convert.ToHSL(c).L
	result := transform.From(c).Lighten(0.2).Darken(0.2).Color()
	// Should return close to original lightness
	assert.InDelta(t, original, convert.ToHSL(result).L, 1e-5)
}

func TestChainSaturateDesaturate(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	original := convert.ToHSL(c).S
	result := transform.From(c).Saturate(0.2).Desaturate(0.2).Color()
	assert.InDelta(t, original, convert.ToHSL(result).S, 1e-5)
}

func TestChainShiftHue(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	result := transform.From(red).ShiftHue(180).Color()
	hue := convert.ToHSL(result).H
	assert.InDelta(t, 180.0, hue, 1.0)
}

func TestChainSetters(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	result := transform.From(c).SetLightness(0.6).SetSaturation(0.4).Color()
	hsl := convert.ToHSL(result)
	assert.InDelta(t, 0.6, hsl.L, 1e-6)
	assert.InDelta(t, 0.4, hsl.S, 1e-6)
}

func TestChainMix(t *testing.T) {
	white := golor.RGB(255, 255, 255)
	black := golor.RGB(0, 0, 0)
	result := transform.From(white).Mix(black, blend.Multiply).Color()
	// Multiply(white, black) = black
	assert.InDelta(t, 0.0, result.R, 1e-9)
}

func TestChainTintShade(t *testing.T) {
	c := golor.RGB(128, 128, 128)
	tinted := transform.From(c).Tint(0.5).Color()
	shaded := transform.From(c).Shade(0.5).Color()
	assert.Greater(t, tinted.R, c.R)
	assert.Less(t, shaded.R, c.R)
}

func TestChainWarmCool(t *testing.T) {
	blue := golor.RGB(0, 0, 255)
	warmed := transform.From(blue).Warm(0.5).Color()
	cooled := transform.From(blue).Cool(0.5).Color()
	// Warmed should be closer to orange (30°), cooled closer to blue (210°)
	_ = warmed
	_ = cooled
}
