package contrast

import "github.com/0mega24/golor/v2"

// EnforceContrast adjusts fg until it meets minRatio contrast against bg after alpha compositing.
// It lightens fg if bg is dark, darkens fg if bg is light. Returns the adjusted color with fg's alpha preserved.
// Uses binary search (bisection) blending toward white or black.
func EnforceContrast(fg, bg golor.Color, minRatio float64) golor.Color {
	if Meets(fg, bg, minRatio) {
		return fg
	}
	bgLum := Luminance(bg)
	var target golor.Color
	if bgLum < 0.5 {
		target = golor.RGB(255, 255, 255)
	} else {
		target = golor.RGB(0, 0, 0)
	}
	lo, hi := 0.0, 1.0
	for i := 0; i < 64; i++ {
		mid := (lo + hi) / 2
		candidate := golor.Color{
			R: fg.R + mid*(target.R-fg.R),
			G: fg.G + mid*(target.G-fg.G),
			B: fg.B + mid*(target.B-fg.B),
			A: fg.A,
		}
		if Meets(candidate, bg, minRatio) {
			hi = mid
		} else {
			lo = mid
		}
	}
	mid := (lo + hi) / 2
	return golor.RGBAf(
		fg.R+mid*(target.R-fg.R),
		fg.G+mid*(target.G-fg.G),
		fg.B+mid*(target.B-fg.B),
		fg.A,
	)
}
