package blend_test

import (
	"testing"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/blend"
	"github.com/stretchr/testify/assert"
)

var (
	black = golor.RGB(0, 0, 0)
	white = golor.RGB(255, 255, 255)
	red   = golor.RGB(255, 0, 0)
	gray  = golor.RGBf(0.5, 0.5, 0.5)
)

func TestMix(t *testing.T) {
	result := blend.Mix(black, white, 0.5)
	assert.InDelta(t, 0.5, result.R, 1e-6)
	assert.InDelta(t, 0.5, result.G, 1e-6)
	assert.InDelta(t, 0.5, result.B, 1e-6)

	assert.Equal(t, black, blend.Mix(black, white, 0))
	assert.Equal(t, white, blend.Mix(black, white, 1))
}

func TestMultiply(t *testing.T) {
	// black * anything = black
	result := blend.Multiply(black, white)
	assert.InDelta(t, 0.0, result.R, 1e-9)
	// white * color = color
	result = blend.Multiply(white, red)
	assert.InDelta(t, 1.0, result.R, 1e-9)
	assert.InDelta(t, 0.0, result.G, 1e-9)
}

func TestScreen(t *testing.T) {
	// screen(white, anything) = white
	result := blend.Screen(white, gray)
	assert.InDelta(t, 1.0, result.R, 1e-9)
	// screen(black, black) = black
	result = blend.Screen(black, black)
	assert.InDelta(t, 0.0, result.R, 1e-9)
}

func TestOverlay(t *testing.T) {
	// overlay(gray, gray): both branches at 0.5 → should be ~0.5
	result := blend.Overlay(gray, gray)
	assert.InDelta(t, 0.5, result.R, 1e-6)
}

func TestHardLight(t *testing.T) {
	// HardLight is Overlay with base/layer swapped
	result := blend.HardLight(gray, gray)
	assert.InDelta(t, 0.5, result.R, 1e-6)
}

func TestDifference(t *testing.T) {
	// difference(white, black) = white
	result := blend.Difference(white, black)
	assert.InDelta(t, 1.0, result.R, 1e-9)
	// difference(c, c) = black
	result = blend.Difference(red, red)
	assert.InDelta(t, 0.0, result.R, 1e-9)
}

func TestSoftLight(t *testing.T) {
	// softlight(base, 0.5) = base (no effect at 50% gray layer)
	result := blend.SoftLight(red, gray)
	assert.InDelta(t, red.R, result.R, 0.01)
}

func TestLuminosity(t *testing.T) {
	// Luminosity: takes L from layer; H,S from base
	// Applying white layer should make result very bright
	result := blend.Luminosity(red, white)
	assert.Greater(t, result.R+result.G+result.B, 2.0)
}
