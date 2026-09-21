package palette_test

import (
	"image"
	stdcolor "image/color"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/palette"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractSolidRegions(t *testing.T) {
	img := solidRegionImage()
	cases := []struct {
		name      string
		algorithm palette.Algorithm
	}{
		{"median-cut", palette.MedianCut},
		{"kmeans", palette.KMeans},
		{"octree", palette.Octree},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := palette.Extract(img, 3, palette.WithAlgorithm(tc.algorithm), palette.WithSeed(7))
			require.Len(t, got, 3)
			assert.Equal(t, golor.RGB(255, 0, 0), got[0])
			assert.Equal(t, golor.RGB(0, 255, 0), got[1])
			assert.Equal(t, golor.RGB(0, 0, 255), got[2])
		})
	}
}

func TestExtractFewerDistinctThanRequested(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	img.SetNRGBA(0, 0, stdcolor.NRGBA{R: 255, A: 255})
	img.SetNRGBA(1, 0, stdcolor.NRGBA{R: 255, A: 255})
	img.SetNRGBA(2, 0, stdcolor.NRGBA{B: 255, A: 255})

	got := palette.Extract(img, 5)
	require.Len(t, got, 2)
	assert.Equal(t, golor.RGB(255, 0, 0), got[0])
	assert.Equal(t, golor.RGB(0, 0, 255), got[1])
}

func TestExtractDegenerateInputs(t *testing.T) {
	assert.Nil(t, palette.Extract(nil, 3))
	assert.Nil(t, palette.Extract(image.NewNRGBA(image.Rect(0, 0, 0, 0)), 3))
	assert.Nil(t, palette.Extract(solidRegionImage(), 0))

	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.SetNRGBA(0, 0, stdcolor.NRGBA{R: 10, G: 20, B: 30, A: 255})
	got := palette.Extract(img, 3, palette.WithAlgorithm(palette.Octree))
	require.Len(t, got, 1)
	assert.Equal(t, golor.RGB(10, 20, 30), got[0])
}

func TestKMeansDeterministic(t *testing.T) {
	img := solidRegionImage()
	a := palette.Extract(img, 2, palette.WithAlgorithm(palette.KMeans), palette.WithSeed(99), palette.WithMaxIterations(10))
	b := palette.Extract(img, 2, palette.WithAlgorithm(palette.KMeans), palette.WithSeed(99), palette.WithMaxIterations(10))
	assert.Equal(t, a, b)
	require.Len(t, a, 2)
}

func TestKMeansMaxIterationsFallback(t *testing.T) {
	got := palette.Extract(solidRegionImage(), 2, palette.WithAlgorithm(palette.KMeans), palette.WithMaxIterations(0))
	require.Len(t, got, 2)
}

func TestAlgorithmsReducePalette(t *testing.T) {
	img := solidRegionImage()
	for _, algorithm := range []palette.Algorithm{palette.MedianCut, palette.Octree} {
		got := palette.Extract(img, 2, palette.WithAlgorithm(algorithm))
		require.Len(t, got, 2)
		for _, c := range got {
			assert.GreaterOrEqual(t, c.R, 0.0)
			assert.LessOrEqual(t, c.R, 1.0)
			assert.GreaterOrEqual(t, c.G, 0.0)
			assert.LessOrEqual(t, c.G, 1.0)
			assert.GreaterOrEqual(t, c.B, 0.0)
			assert.LessOrEqual(t, c.B, 1.0)
		}
	}
}

func solidRegionImage() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 6, 1))
	for x := 0; x < 3; x++ {
		img.SetNRGBA(x, 0, stdcolor.NRGBA{R: 255, A: 255})
	}
	for x := 3; x < 5; x++ {
		img.SetNRGBA(x, 0, stdcolor.NRGBA{G: 255, A: 255})
	}
	img.SetNRGBA(5, 0, stdcolor.NRGBA{B: 255, A: 255})
	return img
}
