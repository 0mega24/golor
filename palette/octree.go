package palette

import "github.com/0mega24/golor/v2"

func extractOctree(swatches []swatch, n int) []golor.Color {
	var previous []swatch
	for bits := 7; bits >= 0; bits-- {
		buckets := octreeBuckets(swatches, bits)
		switch {
		case len(buckets) == n:
			sortSwatches(buckets)
			return swatchColors(buckets)
		case len(buckets) < n && previous != nil:
			reduced := reduceBuckets(previous, n)
			sortSwatches(reduced)
			return swatchColors(reduced)
		case bits == 0:
			sortSwatches(buckets)
			return swatchColors(buckets)
		}
		previous = buckets
	}
	return nil
}

func reduceBuckets(buckets []swatch, n int) []swatch {
	reduced := append([]swatch(nil), buckets...)
	for len(reduced) > n {
		a, b := closestBuckets(reduced)
		if a > b {
			a, b = b, a
		}
		merged := mergeSwatches(reduced[a], reduced[b])
		reduced = append(reduced[:b], reduced[b+1:]...)
		reduced[a] = merged
	}
	return reduced
}

func closestBuckets(buckets []swatch) (int, int) {
	bestA, bestB := 0, 1
	bestDistance := rgbDistanceSquared(buckets[0].color, buckets[1].color)
	for i := 0; i < len(buckets); i++ {
		for j := i + 1; j < len(buckets); j++ {
			distance := rgbDistanceSquared(buckets[i].color, buckets[j].color)
			if distance < bestDistance {
				bestA, bestB = i, j
				bestDistance = distance
			}
		}
	}
	return bestA, bestB
}

func mergeSwatches(a, b swatch) swatch {
	total := a.count + b.count
	return swatch{
		color: golor.RGBAf(
			(a.color.R*float64(a.count)+b.color.R*float64(b.count))/float64(total),
			(a.color.G*float64(a.count)+b.color.G*float64(b.count))/float64(total),
			(a.color.B*float64(a.count)+b.color.B*float64(b.count))/float64(total),
			(a.color.A*float64(a.count)+b.color.A*float64(b.count))/float64(total),
		),
		count: total,
	}
}

func octreeBuckets(swatches []swatch, bits int) []swatch {
	type accum struct {
		r, g, b, a float64
		count      int
	}
	buckets := map[uint32]accum{}
	for _, s := range swatches {
		key := octreeKey(s.color, bits)
		a := buckets[key]
		weight := float64(s.count)
		a.r += s.color.R * weight
		a.g += s.color.G * weight
		a.b += s.color.B * weight
		a.a += s.color.A * weight
		a.count += s.count
		buckets[key] = a
	}
	result := make([]swatch, 0, len(buckets))
	for _, a := range buckets {
		result = append(result, swatch{
			color: golor.RGBAf(a.r/float64(a.count), a.g/float64(a.count), a.b/float64(a.count), a.a/float64(a.count)),
			count: a.count,
		})
	}
	return result
}

func octreeKey(c golor.Color, bits int) uint32 {
	if bits <= 0 {
		return 0
	}
	shift := uint(8 - bits)
	return uint32(c.R8()>>shift)<<16 | uint32(c.G8()>>shift)<<8 | uint32(c.B8()>>shift)
}
