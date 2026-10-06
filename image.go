package golor

import (
	"image"
	stdcolor "image/color"
)

// ImagePixels extracts row-major colors, using direct NRGBA/RGBA reads that match FromStdColor.
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
	switch src := img.(type) {
	case *image.NRGBA:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			row := src.Pix[src.PixOffset(bounds.Min.X, y):]
			for i := 0; i < width*4; i += 4 {
				if row[i+3] == 0 {
					pixels = append(pixels, Color{})
					continue
				}
				pixels = append(pixels, RGBA(row[i], row[i+1], row[i+2], row[i+3]))
			}
		}
	case *image.RGBA:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			row := src.Pix[src.PixOffset(bounds.Min.X, y):]
			for i := 0; i < width*4; i += 4 {
				pixels = append(pixels, fromPremultiplied8(row[i], row[i+1], row[i+2], row[i+3]))
			}
		}
	default:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				pixels = append(pixels, FromStdColor(img.At(x, y)))
			}
		}
	}
	return pixels
}

// fromPremultiplied8 matches FromStdColor(color.RGBA{r, g, b, a}) exactly: scaling both
// numerator and denominator by 0x101 does not change a correctly rounded quotient.
func fromPremultiplied8(r, g, b, a uint8) Color {
	if a == 0 {
		return Color{}
	}
	fa := float64(a)
	return RGBAf(float64(r)/fa, float64(g)/fa, float64(b)/fa, fa/255)
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
	width := bounds.Dx()
	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := dst.Pix[dst.PixOffset(bounds.Min.X, y):]
		for x := 0; x < width; x++ {
			if i >= len(pixels) {
				return
			}
			putNRGBA(row[x*4:], pixels[i].NRGBA())
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
	width := bounds.Dx()
	nrgba := c.NRGBA()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := dst.Pix[dst.PixOffset(bounds.Min.X, y):]
		for x := 0; x < width; x++ {
			putNRGBA(row[x*4:], nrgba)
		}
	}
}

func putNRGBA(p []uint8, c stdcolor.NRGBA) {
	p[0], p[1], p[2], p[3] = c.R, c.G, c.B, c.A
}

// NRGBA returns c as an unpremultiplied standard library color.NRGBA value.
func (c Color) NRGBA() stdcolor.NRGBA {
	return stdcolor.NRGBA{R: c.R8(), G: c.G8(), B: c.B8(), A: c.A8()}
}
