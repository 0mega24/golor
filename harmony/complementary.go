// Package harmony provides color harmony generators based on hue relationships.
package harmony

import "github.com/0mega24/golor"

// Complementary returns c and its complement (hue rotated 180°).
func Complementary(c golor.Color) [2]golor.Color {
	return [2]golor.Color{c, rotateHue(c, 180)}
}
