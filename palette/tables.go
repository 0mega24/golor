package palette

import (
	stdcolor "image/color"
	"sync"

	"github.com/0mega24/golor/v2"
)

// keyTables matches FromStdColor's 8-bit results, including lossy low-alpha
// unpremultiplication, so direct pixel reads agree with img.At.
type keyTables struct {
	straight [256][256]uint8 // Indexed by [alpha][channel] for image.NRGBA.
	premul   [256][256]uint8 // Indexed by [alpha][channel] for image.RGBA.
	alpha    [256]uint8
}

var tables = sync.OnceValue(func() *keyTables {
	t := &keyTables{}
	for a := 0; a < 256; a++ {
		t.alpha[a] = golor.FromStdColor(stdcolor.NRGBA{A: uint8(a)}).A8()
		for v := 0; v < 256; v++ {
			c := uint8(v)
			t.straight[a][v] = golor.FromStdColor(stdcolor.NRGBA{R: c, A: uint8(a)}).R8()
			t.premul[a][v] = golor.FromStdColor(stdcolor.RGBA{R: c, A: uint8(a)}).R8()
		}
	}
	return t
})
