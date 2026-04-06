package harmony

import "github.com/0mega24/golor"

// Analogous returns three colors: c shifted left by angle, c itself, and c shifted right by angle.
// angle is the spread in degrees (e.g. 30).
func Analogous(c golor.Color, angle float64) [3]golor.Color {
	return [3]golor.Color{rotateHue(c, -angle), c, rotateHue(c, angle)}
}
