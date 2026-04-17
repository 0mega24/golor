package gradient_test

import (
	"testing"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/gradient"
	"github.com/stretchr/testify/assert"
)

var (
	black = golor.RGB(0, 0, 0)
	white = golor.RGB(255, 255, 255)
	red   = golor.RGB(255, 0, 0)
	blue  = golor.RGB(0, 0, 255)
)

func TestRGBEndpoints(t *testing.T) {
	result := gradient.RGB(black, white, 5)
	assert.Equal(t, 5, len(result))
	assert.InDelta(t, black.R, result[0].R, 1e-9)
	assert.InDelta(t, white.R, result[4].R, 1e-9)
}

func TestHSLEndpoints(t *testing.T) {
	result := gradient.HSL(red, blue, 5)
	assert.Equal(t, 5, len(result))
	assert.InDelta(t, red.R, result[0].R, 1e-5)
	assert.InDelta(t, blue.B, result[4].B, 1e-5)
}

func TestLCHEndpoints(t *testing.T) {
	result := gradient.LCH(black, white, 3)
	assert.Equal(t, 3, len(result))
	assert.InDelta(t, black.R, result[0].R, 1e-5)
	assert.InDelta(t, white.R, result[2].R, 1e-5)
}

func TestLABEndpoints(t *testing.T) {
	result := gradient.LAB(red, blue, 3)
	assert.Equal(t, 3, len(result))
	assert.InDelta(t, red.R, result[0].R, 1e-5)
	assert.InDelta(t, blue.B, result[2].B, 1e-5)
}

func TestGradientEdgeCases(t *testing.T) {
	assert.Nil(t, gradient.RGB(black, white, 0))
	single := gradient.RGB(red, blue, 1)
	assert.Equal(t, 1, len(single))
	assert.InDelta(t, red.R, single[0].R, 1e-9)
}

func TestMultiStopRGB(t *testing.T) {
	stops := []golor.Color{black, red, white}
	result := gradient.MultiStopRGB(stops, 5)
	assert.Equal(t, 5, len(result))
	assert.InDelta(t, black.R, result[0].R, 1e-9)
	assert.InDelta(t, white.R, result[4].R, 1e-9)
}

func TestMultiStopLCH(t *testing.T) {
	stops := []golor.Color{black, white}
	result := gradient.MultiStopLCH(stops, 3)
	assert.Equal(t, 3, len(result))
	assert.InDelta(t, black.R, result[0].R, 1e-5)
	assert.InDelta(t, white.R, result[2].R, 1e-5)
}
