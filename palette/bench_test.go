package palette_test

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/palette"
)

var sinkColors []golor.Color

// photoLike creates a smooth image with many correlated colors, like a photo with gradients.
func photoLike(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetNRGBA(x, y, color.NRGBA{uint8(x * 255 / w), uint8(y * 255 / h), uint8((x + y) % 64 * 4), 255})
		}
	}
	return img
}

// noise creates an image with nearly unique pixel colors, the worst case for swatch counting.
func noise(w, h int) *image.NRGBA {
	r := rand.New(rand.NewSource(7))
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	r.Read(img.Pix)
	for i := 3; i < len(img.Pix); i += 4 {
		img.Pix[i] = 255
	}
	return img
}

var algorithms = []struct {
	name string
	algo palette.Algorithm
}{{"MedianCut", palette.MedianCut}, {"KMeans", palette.KMeans}, {"Octree", palette.Octree}}

func BenchmarkExtract(b *testing.B) {
	imgs := []struct {
		name string
		img  image.Image
	}{{"photo512", photoLike(512, 512)}, {"noise512", noise(512, 512)}}
	for _, im := range imgs {
		for _, a := range algorithms {
			for _, n := range []int{8, 64} {
				b.Run(fmt.Sprintf("%s/%s/n=%d", im.name, a.name, n), func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						sinkColors = palette.Extract(im.img, n, palette.WithAlgorithm(a.algo))
					}
				})
			}
		}
	}
}

// BenchmarkOctreeLevelBoundary covers n just below and just above an octree level size,
// which used to trigger an O(B^3) bucket merge.
func BenchmarkOctreeLevelBoundary(b *testing.B) {
	img := noise(256, 256)
	for _, n := range []int{500, 520} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sinkColors = palette.Extract(img, n, palette.WithAlgorithm(palette.Octree))
			}
		})
	}
}

func BenchmarkKMeansLarge(b *testing.B) {
	img := photoLike(2048, 2048)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkColors = palette.Extract(img, 16, palette.WithAlgorithm(palette.KMeans))
	}
}
