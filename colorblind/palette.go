package colorblind

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
	"github.com/0mega24/golor/deltae"
)

// AccessiblePalette returns a version of colors that is distinguishable under d.
// Identifies pairs with DeltaE76 < 10 when simulated and shifts hue to restore separability.
func AccessiblePalette(colors []golor.Color, d Deficiency) []golor.Color {
	result := make([]golor.Color, len(colors))
	copy(result, colors)
	const minDeltaE = 10.0
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			si := Simulate(result[i], d)
			sj := Simulate(result[j], d)
			if deltae.DeltaE76(si, sj) < minDeltaE {
				hsl := convert.ToHSL(result[j])
				hsl.H = cbHueShift(hsl.H, 30)
				result[j] = convert.FromHSL(hsl)
			}
		}
	}
	return result
}

func cbHueShift(h, delta float64) float64 {
	h += delta
	for h >= 360 {
		h -= 360
	}
	return h
}
