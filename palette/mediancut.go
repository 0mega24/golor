package palette

import (
	"cmp"
	"slices"

	"github.com/0mega24/golor/v2"
)

// mcItem is a compact swatch for median cut: the packed 8-bit RGBA key and its count.
type mcItem struct {
	key   uint32
	count int
}

// bucket is a median-cut box. channel and spread (the widest RGB channel and its 8-bit
// range) are computed once when the bucket is created.
type bucket struct {
	items   []mcItem
	count   int
	channel uint
	spread  int
}

func newBucket(items []mcItem) bucket {
	b := bucket{items: items}
	if len(items) == 0 {
		return b
	}
	mins := [3]uint32{255, 255, 255}
	var maxs [3]uint32
	for _, it := range items {
		b.count += it.count
		r, g, bl := it.key>>24, (it.key>>16)&0xff, (it.key>>8)&0xff
		mins[0], maxs[0] = min(mins[0], r), max(maxs[0], r)
		mins[1], maxs[1] = min(mins[1], g), max(maxs[1], g)
		mins[2], maxs[2] = min(mins[2], bl), max(maxs[2], bl)
	}
	b.spread = int(maxs[0]) - int(mins[0])
	for c := 1; c < 3; c++ {
		if r := int(maxs[c]) - int(mins[c]); r > b.spread {
			b.spread, b.channel = r, uint(c)
		}
	}
	return b
}

// channel8 extracts the 8-bit R (0), G (1), or B (2) channel from a packed key.
func channel8(key uint32, c uint) uint32 {
	return (key >> (24 - 8*c)) & 0xff
}

func extractMedianCut(swatches []swatch, n int) []golor.Color {
	items := make([]mcItem, len(swatches))
	for i, s := range swatches {
		items[i] = mcItem{key: s.key, count: s.count}
	}
	scratch := make([]mcItem, len(items))
	buckets := []bucket{newBucket(items)}
	for len(buckets) < n {
		index := bucketToSplit(buckets)
		if index < 0 {
			break
		}
		left, right := splitBucket(buckets[index], scratch)
		if len(left.items) == 0 || len(right.items) == 0 {
			break
		}
		buckets = append(buckets[:index], buckets[index+1:]...)
		buckets = append(buckets, left, right)
	}
	slices.SortStableFunc(buckets, func(a, b bucket) int { return cmp.Compare(b.count, a.count) })
	colors := make([]golor.Color, len(buckets))
	for i, bucket := range buckets {
		colors[i] = bucketAverage(bucket.items)
	}
	return colors
}

func bucketAverage(items []mcItem) golor.Color {
	var r, g, b, a, total int
	for _, it := range items {
		r += int(it.key>>24) * it.count
		g += int((it.key>>16)&0xff) * it.count
		b += int((it.key>>8)&0xff) * it.count
		a += int(it.key&0xff) * it.count
		total += it.count
	}
	if total == 0 {
		return golor.RGBAf(0, 0, 0, 0)
	}
	t := float64(total) * 255
	return golor.RGBAf(float64(r)/t, float64(g)/t, float64(b)/t, float64(a)/t)
}

func bucketToSplit(buckets []bucket) int {
	best := -1
	bestScore := -1.0
	for i, bucket := range buckets {
		if len(bucket.items) < 2 {
			continue
		}
		score := float64(bucket.spread) / 255 * float64(bucket.count)
		if score > bestScore {
			bestScore = score
			best = i
		}
	}
	return best
}

// splitBucket orders the bucket's items along its widest channel with a stable counting
// sort (O(len), deterministic) and splits at the count-weighted median.
func splitBucket(b bucket, scratch []mcItem) (bucket, bucket) {
	items := b.items // Owned by b, which is discarded after the split.
	ch := b.channel
	var offsets [257]int
	for _, it := range items {
		offsets[channel8(it.key, ch)+1]++
	}
	for i := 1; i < len(offsets); i++ {
		offsets[i] += offsets[i-1]
	}
	tmp := scratch[:len(items)]
	for _, it := range items {
		v := channel8(it.key, ch)
		tmp[offsets[v]] = it
		offsets[v]++
	}
	copy(items, tmp)

	half := b.count / 2
	running := 0
	split := 1
	for i, it := range items {
		running += it.count
		if running >= half {
			split = i + 1
			break
		}
	}
	if split >= len(items) {
		split = len(items) - 1
	}
	return newBucket(items[:split:split]), newBucket(items[split:])
}
