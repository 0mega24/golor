package palette

import (
	"sort"

	"github.com/0mega24/golor/v2"
)

type bucket struct {
	swatches []swatch
	count    int
}

func extractMedianCut(swatches []swatch, n int) []golor.Color {
	buckets := []bucket{{swatches: append([]swatch(nil), swatches...), count: totalCount(swatches)}}
	for len(buckets) < n {
		index := bucketToSplit(buckets)
		if index < 0 {
			break
		}
		left, right := splitBucket(buckets[index])
		if len(left.swatches) == 0 || len(right.swatches) == 0 {
			break
		}
		buckets = append(buckets[:index], buckets[index+1:]...)
		buckets = append(buckets, left, right)
	}
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].count > buckets[j].count
	})
	colors := make([]golor.Color, len(buckets))
	for i, bucket := range buckets {
		colors[i] = weightedAverage(bucket.swatches)
	}
	return colors
}

func bucketToSplit(buckets []bucket) int {
	best := -1
	bestScore := -1.0
	for i, bucket := range buckets {
		if len(bucket.swatches) < 2 {
			continue
		}
		score := bucketRange(bucket.swatches) * float64(bucket.count)
		if score > bestScore {
			bestScore = score
			best = i
		}
	}
	return best
}

func splitBucket(b bucket) (bucket, bucket) {
	swatches := append([]swatch(nil), b.swatches...)
	channel := widestChannel(swatches)
	sort.Slice(swatches, func(i, j int) bool {
		return channelValue(swatches[i].color, channel) < channelValue(swatches[j].color, channel)
	})
	half := b.count / 2
	running := 0
	split := 1
	for i, s := range swatches {
		running += s.count
		if running >= half {
			split = i + 1
			break
		}
	}
	if split >= len(swatches) {
		split = len(swatches) - 1
	}
	left := swatches[:split]
	right := swatches[split:]
	return bucket{swatches: left, count: totalCount(left)}, bucket{swatches: right, count: totalCount(right)}
}

func bucketRange(swatches []swatch) float64 {
	channel := widestChannel(swatches)
	min, max := 1.0, 0.0
	for _, s := range swatches {
		v := channelValue(s.color, channel)
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return max - min
}

func widestChannel(swatches []swatch) int {
	mins := [3]float64{1, 1, 1}
	maxs := [3]float64{0, 0, 0}
	for _, s := range swatches {
		values := [3]float64{s.color.R, s.color.G, s.color.B}
		for i, v := range values {
			if v < mins[i] {
				mins[i] = v
			}
			if v > maxs[i] {
				maxs[i] = v
			}
		}
	}
	best := 0
	bestRange := maxs[0] - mins[0]
	for i := 1; i < 3; i++ {
		if r := maxs[i] - mins[i]; r > bestRange {
			bestRange = r
			best = i
		}
	}
	return best
}

func channelValue(c golor.Color, channel int) float64 {
	switch channel {
	case 0:
		return c.R
	case 1:
		return c.G
	default:
		return c.B
	}
}
