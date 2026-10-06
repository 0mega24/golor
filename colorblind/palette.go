package colorblind

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

const (
	minAccessibleDeltaE = 10.0
	// MaxAccessiblePaletteAttempts caps hue-shift attempts per color.
	MaxAccessiblePaletteAttempts = 24
)

// AccessiblePalette hue-shifts colors closer than ΔE76 10 under d, composited over white.
// Earlier colors take precedence; alpha is preserved. Separation is best effort.
func AccessiblePalette(colors []golor.Color, d Deficiency) []golor.Color {
	result := make([]golor.Color, len(colors))
	copy(result, colors)
	// Cache Lab values of simulated colors already placed in the palette. Computing them once per
	// placed color keeps the search O(n² · attempts) instead of O(n³ · attempts).
	placed := make([]convert.LAB, 0, len(colors))
	for i := range result {
		if i > 0 && minDistance(placed, result[i], d) < minAccessibleDeltaE {
			result[i] = bestAccessibleCandidate(placed, result[i], d)
		}
		placed = append(placed, simulatedLab(result[i], d))
	}
	return result
}

func bestAccessibleCandidate(placed []convert.LAB, original golor.Color, d Deficiency) golor.Color {
	best := original
	bestDistance := minDistance(placed, original, d)
	hsl := convert.ToHSL(original)
	for attempt := 0; attempt < MaxAccessiblePaletteAttempts; attempt++ {
		candidateHSL := hsl
		candidateHSL.H = cbHueShift(hsl.H, adaptiveShift(attempt))
		candidate := convert.FromHSL(candidateHSL)
		candidate.A = original.A
		distance := minDistance(placed, candidate, d)
		if distance > bestDistance {
			best = candidate
			bestDistance = distance
		}
		if distance >= minAccessibleDeltaE {
			return candidate
		}
	}
	return best
}

// simulatedLab returns the L*a*b* of c as seen under d, composited over white
// (the same preparation deltae.DeltaE76 applies).
func simulatedLab(c golor.Color, d Deficiency) convert.LAB {
	s := Simulate(c, d)
	a := s.A
	return convert.ToLAB(golor.RGBf(s.R*a+1-a, s.G*a+1-a, s.B*a+1-a))
}

func minDistance(placed []convert.LAB, c golor.Color, d Deficiency) float64 {
	current := simulatedLab(c, d)
	m := math.Inf(1)
	for _, p := range placed {
		dl, da, db := p.L-current.L, p.A-current.A, p.B-current.B
		if dist := math.Sqrt(dl*dl + da*da + db*db); dist < m {
			m = dist
		}
	}
	return m
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
