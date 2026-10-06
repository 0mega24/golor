package harmony

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/adjust"
	"github.com/0mega24/golor/v2/convert"
)

// OKLCH schemes preserve alpha and rotate hue with gamut mapping; greys are unchanged.

// ComplementaryOKLCH returns c and its OKLCH complement (hue rotated 180°).
func ComplementaryOKLCH(c golor.Color) [2]golor.Color {
	return [2]golor.Color{c, adjust.ShiftHueOKLCH(c, 180)}
}

// TriadicOKLCH returns c and two colors 120° apart in OKLCH hue.
func TriadicOKLCH(c golor.Color) [3]golor.Color {
	return [3]golor.Color{c, adjust.ShiftHueOKLCH(c, 120), adjust.ShiftHueOKLCH(c, 240)}
}

// AnalogousOKLCH returns c rotated by -angle, c, and c rotated by +angle in OKLCH hue.
func AnalogousOKLCH(c golor.Color, angle float64) [3]golor.Color {
	return [3]golor.Color{adjust.ShiftHueOKLCH(c, -angle), c, adjust.ShiftHueOKLCH(c, angle)}
}

// SplitComplementaryOKLCH returns c and the two colors ±30° from its OKLCH complement.
func SplitComplementaryOKLCH(c golor.Color) [3]golor.Color {
	return [3]golor.Color{c, adjust.ShiftHueOKLCH(c, 150), adjust.ShiftHueOKLCH(c, 210)}
}

// TetradicOKLCH returns four colors 90° apart in OKLCH hue.
func TetradicOKLCH(c golor.Color) [4]golor.Color {
	return [4]golor.Color{c, adjust.ShiftHueOKLCH(c, 90), adjust.ShiftHueOKLCH(c, 180), adjust.ShiftHueOKLCH(c, 270)}
}

// ExpandOKLCH samples lightness ±0.35 around c, clamped to [0.05,0.98].
// Hue/chroma and alpha are preserved; results are gamut-mapped.
func ExpandOKLCH(c golor.Color, n int) []golor.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []golor.Color{c}
	}
	o := convert.ToOKLCH(c)
	lo, hi := max(0.05, o.L-0.35), min(0.98, o.L+0.35)
	out := make([]golor.Color, n)
	for i := range out {
		step := o
		step.L = lo + float64(i)/float64(n-1)*(hi-lo)
		out[i] = convert.FromOKLCHMapped(step)
		out[i].A = c.A
	}
	return out
}
