package palette

import (
	"math"

	"github.com/0mega24/golor/v2"
)

func extractKMeans(swatches []swatch, n int, seed int64, maxIterations int) []golor.Color {
	centroids := initialCentroids(swatches, n, seed)
	assignments := make([]int, len(swatches))
	for i := range assignments {
		assignments[i] = -1
	}
	for iter := 0; iter < maxIterations; iter++ {
		changed := false
		for i, s := range swatches {
			next := nearestCentroid(s.color, centroids)
			if assignments[i] != next {
				assignments[i] = next
				changed = true
			}
		}
		centroids = recomputeCentroids(swatches, assignments, centroids)
		if !changed {
			break
		}
	}
	return clustersToColors(swatches, assignments, centroids)
}

func initialCentroids(swatches []swatch, n int, seed int64) []golor.Color {
	rng := seededRand(seed)
	centroids := []golor.Color{swatches[0].color}
	for len(centroids) < n {
		distances := make([]float64, len(swatches))
		total := 0.0
		for i, s := range swatches {
			distance := rgbDistanceSquared(s.color, centroids[nearestCentroid(s.color, centroids)])
			weighted := distance * float64(s.count)
			distances[i] = weighted
			total += weighted
		}
		if total == 0 {
			break
		}
		target := rng.Float64() * total
		running := 0.0
		chosen := len(swatches) - 1
		for i, distance := range distances {
			running += distance
			if running >= target {
				chosen = i
				break
			}
		}
		centroids = append(centroids, swatches[chosen].color)
	}
	return centroids
}

func recomputeCentroids(swatches []swatch, assignments []int, previous []golor.Color) []golor.Color {
	sums := make([][4]float64, len(previous))
	counts := make([]int, len(previous))
	for i, s := range swatches {
		cluster := assignments[i]
		weight := float64(s.count)
		sums[cluster][0] += s.color.R * weight
		sums[cluster][1] += s.color.G * weight
		sums[cluster][2] += s.color.B * weight
		sums[cluster][3] += s.color.A * weight
		counts[cluster] += s.count
	}
	centroids := make([]golor.Color, len(previous))
	for i := range previous {
		if counts[i] == 0 {
			centroids[i] = previous[i]
			continue
		}
		centroids[i] = golor.RGBAf(
			sums[i][0]/float64(counts[i]),
			sums[i][1]/float64(counts[i]),
			sums[i][2]/float64(counts[i]),
			sums[i][3]/float64(counts[i]),
		)
	}
	return centroids
}

func clustersToColors(swatches []swatch, assignments []int, centroids []golor.Color) []golor.Color {
	counts := make([]int, len(centroids))
	for i, s := range swatches {
		counts[assignments[i]] += s.count
	}
	clusters := make([]swatch, 0, len(centroids))
	for i, c := range centroids {
		if counts[i] == 0 || math.IsNaN(c.R) {
			continue
		}
		clusters = append(clusters, swatch{color: c, count: counts[i]})
	}
	sortSwatches(clusters)
	return swatchColors(clusters)
}
