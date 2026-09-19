package contrast

import "github.com/0mega24/golor/v2"

// Ratio returns the WCAG 2.1 contrast ratio between two colors in [1, 21].
// Non-opaque colors are composited in order: c2 over white, then c1 over that background.
func Ratio(c1, c2 golor.Color) float64 {
	fg, bg := resolvePair(c1, c2)
	l1 := Luminance(fg)
	l2 := Luminance(bg)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// Meets reports whether c1 and c2 meet the given minimum contrast ratio after Ratio's alpha compositing.
func Meets(c1, c2 golor.Color, minRatio float64) bool {
	return Ratio(c1, c2) >= minRatio
}

func resolvePair(c1, c2 golor.Color) (golor.Color, golor.Color) {
	white := golor.RGB(255, 255, 255)
	bg := composite(c2, white)
	return composite(c1, bg), bg
}

func composite(fg, bg golor.Color) golor.Color {
	a := fg.A
	return golor.RGBf(
		fg.R*a+bg.R*(1-a),
		fg.G*a+bg.G*(1-a),
		fg.B*a+bg.B*(1-a),
	)
}
