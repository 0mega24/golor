package colorblind_test

import (
	"testing"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/colorblind"
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
