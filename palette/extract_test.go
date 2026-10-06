package palette_test

import (
	"image"
	stdcolor "image/color"
	"math"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
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

func TestAlgorithmsRecoverQuadrantColors(t *testing.T) {
	img, base := quadrants(128)
	for _, a := range algorithms {
		got := palette.Extract(img, 4, palette.WithAlgorithm(a.algo))
		require.Len(t, got, 4, a.name)
		for _, want := range base {
			best := math.Inf(1)
			for _, c := range got {
				best = math.Min(best, convert.DeltaEOK(c, want))
			}
			assert.Less(t, best, 0.02, "%s: no extracted color near %s (got %v)", a.name, want, got)
		}
	}
}

func TestKMeansIndependentOfGOMAXPROCS(t *testing.T) {
	img := noise(256, 256) // Bins into enough points to run the assignment step in parallel.
	prev := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(prev)
	want := palette.Extract(img, 16, palette.WithAlgorithm(palette.KMeans), palette.WithSeed(3))
	for _, procs := range []int{2, 8} {
		runtime.GOMAXPROCS(procs)
		got := palette.Extract(img, 16, palette.WithAlgorithm(palette.KMeans), palette.WithSeed(3))
		assert.Equal(t, want, got, "GOMAXPROCS=%d", procs)
	}
	runtime.GOMAXPROCS(prev)
	other := palette.Extract(img, 16, palette.WithAlgorithm(palette.KMeans), palette.WithSeed(4))
	assert.NotEqual(t, want, other, "seed should matter")
}

func TestOctreeLevelBoundary(t *testing.T) {
	img := noise(256, 256)
	start := time.Now()
	got := palette.Extract(img, 520, palette.WithAlgorithm(palette.Octree))
	assert.Len(t, got, 520)
	// O(B^3) merging takes ~15s; the current implementation takes about 0.1s. The bound allows for slow CI.
	assert.Less(t, time.Since(start), 10*time.Second)
	again := palette.Extract(img, 520, palette.WithAlgorithm(palette.Octree))
	assert.Equal(t, got, again)
}

func TestExtractReturnsExactlyNWhenPossible(t *testing.T) {
	img := noise(64, 64)
	for _, a := range algorithms {
		for _, n := range []int{1, 7, 64, 300} {
			assert.Len(t, palette.Extract(img, n, palette.WithAlgorithm(a.algo)), n, "%s n=%d", a.name, n)
		}
	}
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
