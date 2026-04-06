package harmony

import "github.com/0mega24/golor"

// Tetradic returns four colors equally spaced around the hue wheel (90° apart).
func Tetradic(c golor.Color) [4]golor.Color {
	return [4]golor.Color{c, rotateHue(c, 90), rotateHue(c, 180), rotateHue(c, 270)}
}
