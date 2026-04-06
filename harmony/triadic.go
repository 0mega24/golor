package harmony

import "github.com/0mega24/golor"

// Triadic returns c and two colors equally spaced around the hue wheel (120° apart).
func Triadic(c golor.Color) [3]golor.Color {
	return [3]golor.Color{c, rotateHue(c, 120), rotateHue(c, 240)}
}
