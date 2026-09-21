package colorblind

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/deltae"
)

const (
	minAccessibleDeltaE = 10.0
	// MaxAccessiblePaletteAttempts bounds the adaptive hue-shift search per color.
	// If no attempted shift fully separates a color from prior colors, AccessiblePalette returns
	// the best candidate found for that color.
	MaxAccessiblePaletteAttempts = 24
)

// AccessiblePalette returns a best-effort version of colors that is distinguishable under d.
// It preserves alpha and adaptively searches bounded hue shifts for colors whose simulated
// DeltaE76 distance is below 10. Full pairwise separability is not guaranteed for all inputs
// such as palettes clustered around one hue; after MaxAccessiblePaletteAttempts per color,
// the best candidate found is returned rather than hanging or returning an error.
func AccessiblePalette(colors []golor.Color, d Deficiency) []golor.Color {
	result := make([]golor.Color, len(colors))
	copy(result, colors)
	for i := 0; i < len(result); i++ {
		if i == 0 || minDistanceToPrior(result, i, d) >= minAccessibleDeltaE {
			continue
		}
		result[i] = bestAccessibleCandidate(result, i, d)
	}
	return result
}

func bestAccessibleCandidate(colors []golor.Color, index int, d Deficiency) golor.Color {
	original := colors[index]
	best := original
	bestDistance := minDistanceToPrior(colors, index, d)
	hsl := convert.ToHSL(original)
	for attempt := 0; attempt < MaxAccessiblePaletteAttempts; attempt++ {
		candidateHSL := hsl
		candidateHSL.H = cbHueShift(hsl.H, adaptiveShift(attempt))
		candidate := convert.FromHSL(candidateHSL)
		candidate.A = original.A
		colors[index] = candidate
		distance := minDistanceToPrior(colors, index, d)
		if distance > bestDistance {
			best = candidate
			bestDistance = distance
		}
		if distance >= minAccessibleDeltaE {
			colors[index] = original
			return candidate
		}
	}
	colors[index] = original
	return best
}

func minDistanceToPrior(colors []golor.Color, index int, d Deficiency) float64 {
	if index <= 0 {
		return math.Inf(1)
	}
	current := Simulate(colors[index], d)
	minDistance := math.Inf(1)
	for i := 0; i < index; i++ {
		distance := deltae.DeltaE76(Simulate(colors[i], d), current)
		if distance < minDistance {
			minDistance = distance
		}
	}
	return minDistance
}

func adaptiveShift(attempt int) float64 {
	steps := []float64{15, 30, 45, 60, 90, 120, 150, 180, 210, 240, 270, 300, 330}
	if attempt < len(steps) {
		return steps[attempt]
	}
	return float64((attempt-len(steps)+1)*37) + 15
}

func cbHueShift(h, delta float64) float64 {
	h = math.Mod(h+delta, 360)
	if h < 0 {
		h += 360
	}
	return h
}
