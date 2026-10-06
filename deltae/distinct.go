package deltae

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// EnsureDistinct returns a copy separated by minDeltaE (ΔE76 over white), preserving alpha.
// Earlier colors take precedence; later colors move in gamut-mapped LCH.
// Separation is best effort when the requested spacing is unreachable.
func EnsureDistinct(colors []golor.Color, minDeltaE float64) []golor.Color {
	result := make([]golor.Color, len(colors))
	copy(result, colors)
	if len(result) < 2 || !(minDeltaE > 0) {
		return result
	}
	labs := make([]convert.LAB, len(result))
	for i, c := range result {
		labs[i] = convert.ToLAB(compositeWhite(c))
	}
	stuck := make([]bool, len(result))
	for round := 0; round < 32*len(result); round++ {
		i, j, d := closestPair(labs, stuck, minDeltaE)
		if j < 0 || d >= minDeltaE {
			return result
		}
		moved, lab, ok := bestMove(result[j], j, labs, minDeltaE)
		if !ok {
			stuck[j] = true // No move improves j; try the remaining pairs.
			if i >= 0 && !stuck[i] {
				// Let the earlier color move as a last resort.
				if moved, lab, ok = bestMove(result[i], i, labs, minDeltaE); ok {
					result[i], labs[i] = moved, lab
					stuck[j] = false
				}
			}
			continue
		}
		result[j], labs[j] = moved, lab
	}
	return result
}

// closestPair returns the closest pair (i < j) whose later member is not stuck.
func closestPair(labs []convert.LAB, stuck []bool, minDeltaE float64) (bi, bj int, best float64) {
	bi, bj, best = -1, -1, math.Inf(1)
	for i := range labs {
		for j := i + 1; j < len(labs); j++ {
			if stuck[j] && stuck[i] {
				continue
			}
			if d := labDistance(labs[i], labs[j]); d < best {
				bi, bj, best = i, j, d
			}
			if best == 0 && minDeltaE > 0 {
				return bi, bj, best
			}
		}
	}
	return bi, bj, best
}

// bestMove tries increasing LCH moves; it returns a passing move or the best improvement.
func bestMove(c golor.Color, idx int, labs []convert.LAB, minDeltaE float64) (golor.Color, convert.LAB, bool) {
	current := minDistance(labs[idx], idx, labs)
	base := convert.ToLCH(c)
	bestColor, bestLab, bestDist := c, labs[idx], current
	for _, mag := range []float64{0.5, 1, 1.5, 2, 3, 4, 6} {
		step := mag * minDeltaE
		candidates := []convert.LCH{
			{L: base.L + step, C: base.C, H: base.H},
			{L: base.L - step, C: base.C, H: base.H},
			{L: base.L, C: base.C + step, H: base.H},
			{L: base.L, C: math.Max(0, base.C-step), H: base.H},
		}
		if base.C <= 1 {
			// Greys have no hue to keep, so try chroma in six directions.
			for h := 0.0; h < 360; h += 60 {
				candidates = append(candidates, convert.LCH{L: base.L, C: step, H: h})
			}
		} else {
			// Hue rotation whose chord length is step.
			dh := 2 * math.Asin(math.Min(1, step/(2*base.C))) * 180 / math.Pi
			candidates = append(candidates,
				convert.LCH{L: base.L, C: base.C, H: base.H + dh},
				convert.LCH{L: base.L, C: base.C, H: base.H - dh})
		}
		for _, cand := range candidates {
			cand.L = math.Max(0, math.Min(100, cand.L))
			col := convert.FromLCHMapped(cand)
			col.A = c.A
			lab := convert.ToLAB(compositeWhite(col))
			d := minDistance(lab, idx, labs)
			if d > bestDist+1e-9 {
				bestColor, bestLab, bestDist = col, lab, d
			}
		}
		if bestDist >= minDeltaE {
			break
		}
	}
	return bestColor, bestLab, bestDist > current+1e-9
}

func minDistance(lab convert.LAB, skip int, labs []convert.LAB) float64 {
	m := math.Inf(1)
	for k, other := range labs {
		if k != skip {
			m = math.Min(m, labDistance(lab, other))
		}
	}
	return m
}

func labDistance(a, b convert.LAB) float64 {
	dl, da, db := a.L-b.L, a.A-b.A, a.B-b.B
	return math.Sqrt(dl*dl + da*da + db*db)
}
