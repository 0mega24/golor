package adjust

import (
	"math"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// Warm shifts the hue toward orange (~30°) by amount (0=no shift, 1=full shift to orange).
func Warm(c golor.Color, amount float64) golor.Color {
	return shiftToward(c, 30, amount)
}

// Cool shifts the hue toward blue (~210°) by amount (0=no shift, 1=full shift to blue).
func Cool(c golor.Color, amount float64) golor.Color {
	return shiftToward(c, 210, amount)
}

// shiftToward moves the hue of c toward targetDeg by the given fraction of the remaining distance.
func shiftToward(c golor.Color, targetDeg, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	diff := targetDeg - hsl.H
	// Take shortest path around the hue circle.
	for diff > 180 {
		diff -= 360
	}
	for diff < -180 {
		diff += 360
	}
	hsl.H = math.Mod(hsl.H+diff*clamp01(amount)+360*10, 360)
	return convert.FromHSL(hsl)
}
