package harmony

import "github.com/0mega24/golor"

// SplitComplementary returns c and two colors flanking its complement (±30° from 180°).
func SplitComplementary(c golor.Color) [3]golor.Color {
	return [3]golor.Color{c, rotateHue(c, 150), rotateHue(c, 210)}
}
