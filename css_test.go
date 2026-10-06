package golor_test

import (
	"image"
	stdcolor "image/color"
	"math"
	"strings"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHexForms(t *testing.T) {
	cases := map[string]string{
		"fff": "#ffffff", "#abc": "#aabbcc", "abcd": "#aabbccdd", "#ABCDEF": "#abcdef",
		"\t#00ff0080\n": "#00ff0080", "#000f": "#000000",
	}
	for in, want := range cases {
		c, err := golor.Hex(in)
		require.NoError(t, err, in)
		assert.Equal(t, want, c.String(), in)
	}
	for _, in := range []string{"", "#", "#f", "#ff", "#fffff", "#fffffff", "#fffffffff", "##fff", "#+ff", "#0xffff", "# fff"} {
		_, err := golor.Hex(in)
		assert.Error(t, err, in)
	}
	_, err := golor.Hex("#zzz")
	assert.EqualError(t, err, `golor: invalid hex color "#zzz"`)
}

func TestNaNAndOutOfRange(t *testing.T) {
	c := golor.RGBf(math.NaN(), 0.5, math.Inf(1))
	assert.Equal(t, 0.0, c.R)
	assert.Equal(t, 1.0, c.B)
	assert.Equal(t, 0.0, golor.RGBAf(0, 0, 0, math.NaN()).A)

	lit := golor.Color{R: math.NaN(), G: 2, B: -1, A: math.NaN()}
	assert.Equal(t, "#00ff0000", lit.String())
	assert.Equal(t, uint8(0), lit.R8())
	r, g, b, a := lit.RGBA()
	assert.Equal(t, [4]uint32{0, 0, 0, 0}, [4]uint32{r, g, b, a})
}

func TestStringFormat(t *testing.T) {
	assert.Equal(t, "#0a0b0c", golor.RGB(10, 11, 12).String())
	assert.Equal(t, "#0a0b0cfe", golor.RGBA(10, 11, 12, 254).String())
	assert.Equal(t, "#ffffff00", golor.RGBA(255, 255, 255, 0).String())
}

// genericImage hides the concrete type so ImagePixels takes its img.At fallback.
type genericImage struct{ image.Image }

func TestImagePixelsFastPathsMatchFromStdColor(t *testing.T) {
	nrgba := image.NewNRGBA(image.Rect(-3, 5, 253, 261))
	rgba := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for a := 0; a < 256; a++ {
		for v := 0; v < 256; v++ {
			x, y := v, a
			nrgba.SetNRGBA(x-3, y+5, stdcolor.NRGBA{uint8(v), uint8(255 - v), uint8(v / 3), uint8(a)})
			// Include invalid premultiplied values (channel > alpha).
			rgba.SetRGBA(x, y, stdcolor.RGBA{uint8(v), uint8(v / 2), uint8(255 - v), uint8(a)})
		}
	}
	for _, img := range []image.Image{
		nrgba, rgba,
		nrgba.SubImage(image.Rect(10, 20, 50, 30)),
		rgba.SubImage(image.Rect(100, 3, 101, 200)),
	} {
		fast := golor.ImagePixels(img)
		slow := golor.ImagePixels(genericImage{img})
		require.Equal(t, len(slow), len(fast))
		for i := range fast {
			if fast[i] != slow[i] {
				t.Fatalf("%T pixel %d: fast %+v, generic %+v", img, i, fast[i], slow[i])
			}
		}
	}
}

func TestFromStdColorStraightAlphaIsLossless(t *testing.T) {
	for a := 1; a < 256; a++ {
		for v := 0; v < 256; v++ {
			src := stdcolor.NRGBA{uint8(v), uint8(255 - v), uint8(v / 2), uint8(a)}
			if got := golor.FromStdColor(src).NRGBA(); got != src {
				t.Fatalf("%v -> %v", src, got)
			}
		}
	}
	c := golor.RGBAf(0.123456789, 0.5, 1, 0.001)
	assert.Equal(t, c, golor.FromStdColor(c))
	assert.Equal(t, golor.Color{}, golor.FromStdColor(stdcolor.NRGBA{R: 9, A: 0}))
	assert.Equal(t, golor.RGBA(255, 0, 0, 128), golor.FromStdColor(stdcolor.NRGBA64{R: 0xffff, A: 0x8080}))
}

func TestPaintAndFill(t *testing.T) {
	img := image.NewNRGBA(image.Rect(2, 2, 5, 4))
	golor.Fill(img, golor.RGBA(1, 2, 3, 4))
	assert.Equal(t, stdcolor.NRGBA{1, 2, 3, 4}, img.NRGBAAt(4, 3))
	golor.PaintPixels(img, []golor.Color{golor.RGB(9, 9, 9), golor.RGB(8, 8, 8), golor.RGB(7, 7, 7), golor.RGB(6, 6, 6)})
	assert.Equal(t, stdcolor.NRGBA{9, 9, 9, 255}, img.NRGBAAt(2, 2))
	assert.Equal(t, stdcolor.NRGBA{6, 6, 6, 255}, img.NRGBAAt(2, 3))
	assert.Equal(t, stdcolor.NRGBA{1, 2, 3, 4}, img.NRGBAAt(3, 3))
	sub := img.SubImage(image.Rect(3, 3, 5, 4)).(*image.NRGBA)
	golor.Fill(sub, golor.RGB(0, 0, 0))
	assert.Equal(t, stdcolor.NRGBA{0, 0, 0, 255}, img.NRGBAAt(4, 3))
	assert.Equal(t, stdcolor.NRGBA{6, 6, 6, 255}, img.NRGBAAt(2, 3))
}

func FuzzHex(f *testing.F) {
	for _, s := range []string{"#fff", "ffff", "#ff6b35", "#ff6b3580", "  #ABC ", "#", "#zz"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		c, err := golor.Hex(s)
		if err != nil {
			return
		}
		h := strings.TrimPrefix(strings.TrimSpace(s), "#")
		if n := len(h); n != 3 && n != 4 && n != 6 && n != 8 {
			t.Fatalf("Hex(%q) accepted length %d", s, n)
		}
		back, err := golor.Hex(c.String())
		if err != nil || back != c {
			t.Fatalf("Hex(%q) = %s, re-parsed %s (%v)", s, c, back, err)
		}
	})
}
