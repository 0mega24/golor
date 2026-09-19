package golor

import (
	"image"
	stdcolor "image/color"
)

// ImagePixels extracts every pixel from img as Colors in row-major order.
func ImagePixels(img image.Image) []Color {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width <= 0 || height <= 0 {
		return nil
	}
	pixels := make([]Color, 0, width*height)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixels = append(pixels, FromStdColor(img.At(x, y)))
		}
	}
	return pixels
}

// NewNRGBA creates an *image.NRGBA from pixels in row-major order.
// Extra pixels are ignored; missing pixels remain transparent.
func NewNRGBA(width, height int, pixels []Color) *image.NRGBA {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	PaintPixels(dst, pixels)
	return dst
}

// PaintPixels writes pixels into dst in row-major order.
// Extra pixels are ignored; missing pixels leave the existing destination unchanged.
func PaintPixels(dst *image.NRGBA, pixels []Color) {
	if dst == nil {
		return
	}
	bounds := dst.Bounds()
	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if i >= len(pixels) {
				return
			}
			dst.SetNRGBA(x, y, pixels[i].NRGBA())
			i++
		}
	}
}

// Fill paints every pixel in dst with c.
func Fill(dst *image.NRGBA, c Color) {
	if dst == nil {
		return
	}
	bounds := dst.Bounds()
	nrgba := c.NRGBA()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.SetNRGBA(x, y, nrgba)
		}
	}
}

// NRGBA returns c as an unpremultiplied standard library color.NRGBA value.
func (c Color) NRGBA() stdcolor.NRGBA {
	return stdcolor.NRGBA{R: c.R8(), G: c.G8(), B: c.B8(), A: c.A8()}
}
