package palette_test

import (
	"image"
	stdcolor "image/color"
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/palette"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// iconImage is a transparent 100x100 image with a red square and a blue top row.
func iconImage() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	for y := 40; y < 60; y++ {
		for x := 40; x < 60; x++ {
			img.SetNRGBA(x, y, stdcolor.NRGBA{220, 40, 40, 255})
		}
	}
	for x := 0; x < 100; x++ {
		img.SetNRGBA(x, 0, stdcolor.NRGBA{40, 40, 220, 255})
	}
	return img
}

func TestTransparentPixelsIgnored(t *testing.T) {
	for _, a := range algorithms {
		got := palette.Extract(iconImage(), 3, palette.WithAlgorithm(a.algo))
		require.Len(t, got, 2, a.name)
		assert.Equal(t, golor.RGB(220, 40, 40), got[0], a.name)
		assert.Equal(t, golor.RGB(40, 40, 220), got[1], a.name)
	}
	// WithMinAlpha(0) counts transparent pixels, as v2.0.0 did.
	got := palette.Extract(iconImage(), 3, palette.WithMinAlpha(0))
	require.Len(t, got, 3)
	assert.Equal(t, golor.RGBA(0, 0, 0, 0), got[0])
	// A fully transparent image has no colors.
	assert.Nil(t, palette.Extract(image.NewNRGBA(image.Rect(0, 0, 4, 4)), 3))
}

func TestMinAlphaThreshold(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	img.SetNRGBA(0, 0, stdcolor.NRGBA{R: 255, A: 255})
	img.SetNRGBA(1, 0, stdcolor.NRGBA{G: 255, A: 100})
	img.SetNRGBA(2, 0, stdcolor.NRGBA{B: 255, A: 10})
	assert.Len(t, palette.Extract(img, 5), 3)
	assert.Len(t, palette.Extract(img, 5, palette.WithMinAlpha(50)), 2)
	assert.Len(t, palette.Extract(img, 5, palette.WithMinAlpha(255)), 1)
}

// quadrants returns an image of four solid colors with ±8 per-channel noise.
func quadrants(size int) (*image.NRGBA, []golor.Color) {
	base := []golor.Color{golor.RGB(200, 30, 40), golor.RGB(30, 160, 60), golor.RGB(40, 60, 200), golor.RGB(230, 210, 60)}
	r := rand.New(rand.NewSource(5))
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	jitter := func(v uint8) uint8 { return uint8(int(v) + r.Intn(17) - 8) }
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c := base[(y*2/size)*2+x*2/size].NRGBA()
			img.SetNRGBA(x, y, stdcolor.NRGBA{jitter(c.R), jitter(c.G), jitter(c.B), 255})
		}
	}
	return img, base
}

func TestGenericImageType(t *testing.T) {
	img, _ := quadrants(32)
	gray := image.NewGray(image.Rect(0, 0, 8, 8))
	for i := range gray.Pix {
		gray.Pix[i] = uint8(i % 3 * 100)
	}
	got := palette.Extract(gray, 3)
	require.Len(t, got, 3)
	assert.Equal(t, palette.Extract(img, 4), palette.Extract(img.SubImage(img.Bounds()), 4))
}
