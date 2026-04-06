package contrast

import "github.com/0mega24/golor"

// ContrastRatio returns the WCAG 2.1 contrast ratio between two colors in [1, 21].
func ContrastRatio(c1, c2 golor.Color) float64 {
	l1 := Luminance(c1)
	l2 := Luminance(c2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// MeetsContrast reports whether c1 and c2 meet the given minimum contrast ratio.
func MeetsContrast(c1, c2 golor.Color, minRatio float64) bool {
	return ContrastRatio(c1, c2) >= minRatio
}
