package colorblind_test

import (
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/adjust"
	"github.com/0mega24/golor/v2/colorblind"
	"github.com/stretchr/testify/assert"
)

func TestSimulate(t *testing.T) {
	cases := []struct {
		name string
		d    colorblind.Deficiency
	}{
		{"deuteranopia", colorblind.Deuteranopia},
		{"protanopia", colorblind.Protanopia},
		{"tritanopia", colorblind.Tritanopia},
	}
	c := golor.RGB(100, 200, 50)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := colorblind.Simulate(c, tc.d)
			assert.GreaterOrEqual(t, result.R, 0.0)
			assert.LessOrEqual(t, result.R, 1.0)
			assert.GreaterOrEqual(t, result.G, 0.0)
			assert.LessOrEqual(t, result.G, 1.0)
			assert.GreaterOrEqual(t, result.B, 0.0)
			assert.LessOrEqual(t, result.B, 1.0)
		})
	}
}

func TestAccessiblePalette(t *testing.T) {
	colors := []golor.Color{
		golor.RGB(255, 0, 0),
		golor.RGB(0, 255, 0),
		golor.RGB(0, 0, 255),
	}
	result := colorblind.AccessiblePalette(colors, colorblind.Deuteranopia)
	assert.Equal(t, len(colors), len(result))
}

func TestColorblindPreservesAlpha(t *testing.T) {
	c := golor.RGBA(100, 200, 50, 77)
	assert.Equal(t, c.A8(), colorblind.Simulate(c, colorblind.Deuteranopia).A8())

	result := colorblind.AccessiblePalette([]golor.Color{c, c}, colorblind.Deuteranopia)
	assert.Equal(t, c.A8(), result[0].A8())
	assert.Equal(t, c.A8(), result[1].A8())
}

func TestAccessiblePalettePathologicalPaletteReturns(t *testing.T) {
	colors := make([]golor.Color, 10)
	base := golor.RGBA(180, 60, 50, 128)
	for i := range colors {
		colors[i] = adjust.ShiftHue(base, float64(i%5))
	}

	result := colorblind.AccessiblePalette(colors, colorblind.Deuteranopia)
	assert.Len(t, result, len(colors))
	for _, c := range result {
		assert.Equal(t, base.A8(), c.A8())
	}
}

func TestMaxAccessiblePaletteAttemptsIsBounded(t *testing.T) {
	assert.Greater(t, colorblind.MaxAccessiblePaletteAttempts, 0)
	assert.LessOrEqual(t, colorblind.MaxAccessiblePaletteAttempts, 64)
}
