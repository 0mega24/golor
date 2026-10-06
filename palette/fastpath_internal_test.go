package palette

import (
	"image"
	stdcolor "image/color"
	"slices"
	"testing"
)

// opaqueWrapper hides the concrete image type so pixelKeys takes the generic img.At path.
type opaqueWrapper struct{ image.Image }

func TestFastPathsMatchGenericPath(t *testing.T) {
	nrgba := image.NewNRGBA(image.Rect(3, 5, 3+256, 5+256)) // A non-zero origin exercises PixOffset.
	rgba := image.NewRGBA(nrgba.Bounds())
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			v, a := uint8(x), uint8(y)
			nrgba.SetNRGBA(3+x, 5+y, stdcolor.NRGBA{R: v, G: 255 - v, B: v / 3, A: a})
			// Includes invalid premultiplied values (channel > alpha).
			rgba.SetRGBA(3+x, 5+y, stdcolor.RGBA{R: v, G: v / 2, B: 255 - v, A: a})
		}
	}
	sub := nrgba.SubImage(image.Rect(40, 60, 200, 100)).(*image.NRGBA)
	for _, img := range []image.Image{nrgba, rgba, sub} {
		for _, minAlpha := range []uint8{0, 1, 128} {
			fast := pixelKeys(img, img.Bounds(), minAlpha)
			slow := pixelKeys(opaqueWrapper{img}, img.Bounds(), minAlpha)
			if !slices.Equal(fast, slow) {
				t.Fatalf("%T minAlpha=%d: fast path differs from generic path", img, minAlpha)
			}
		}
	}
}

func TestRadixSort(t *testing.T) {
	keys := []uint32{}
	for i := 0; i < 5000; i++ {
		keys = append(keys, uint32(i*2654435761)|0xff)
	}
	want := slices.Clone(keys)
	slices.Sort(want)
	radixSort(keys)
	if !slices.Equal(keys, want) {
		t.Fatal("radixSort mismatch")
	}
}
