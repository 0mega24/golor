package harmony_test

import (
	"math"
	"testing"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
	"github.com/0mega24/golor/harmony"
	"github.com/stretchr/testify/assert"
)

func hue(c golor.Color) float64 { return convert.ToHSL(c).H }

func hueDist(h1, h2 float64) float64 {
	d := math.Abs(h1 - h2)
	if d > 180 {
		d = 360 - d
	}
	return d
}

func TestComplementary(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	pair := harmony.Complementary(red)
	assert.InDelta(t, 180.0, hueDist(hue(pair[0]), hue(pair[1])), 1.0)
}

func TestTriadic(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	trio := harmony.Triadic(red)
	assert.InDelta(t, 120.0, hueDist(hue(trio[0]), hue(trio[1])), 1.0)
	assert.InDelta(t, 120.0, hueDist(hue(trio[1]), hue(trio[2])), 1.0)
}

func TestAnalogous(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	trio := harmony.Analogous(red, 30)
	assert.InDelta(t, 30.0, hueDist(hue(trio[0]), hue(trio[1])), 1.0)
	assert.InDelta(t, 30.0, hueDist(hue(trio[1]), hue(trio[2])), 1.0)
}

func TestTetradic(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	quad := harmony.Tetradic(red)
	assert.InDelta(t, 90.0, hueDist(hue(quad[0]), hue(quad[1])), 1.0)
}

func TestSplitComplementary(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	trio := harmony.SplitComplementary(red)
	assert.InDelta(t, 150.0, hueDist(hue(trio[0]), hue(trio[1])), 1.0)
	assert.InDelta(t, 150.0, hueDist(hue(trio[0]), hue(trio[2])), 1.0)
}

func TestExpand(t *testing.T) {
	c := golor.RGB(100, 150, 200)
	result := harmony.Expand(c, 5)
	assert.Equal(t, 5, len(result))
	// Should go from darker to lighter
	l0 := convert.ToHSL(result[0]).L
	l4 := convert.ToHSL(result[4]).L
	assert.Less(t, l0, l4)
}

func TestExpandEdgeCases(t *testing.T) {
	c := golor.RGB(128, 128, 128)
	assert.Nil(t, harmony.Expand(c, 0))
	single := harmony.Expand(c, 1)
	assert.Equal(t, 1, len(single))
	assert.Equal(t, c, single[0])
}
