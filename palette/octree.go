package palette

import (
	"math"
	"slices"

	"github.com/0mega24/golor/v2"
)

// octNode is an octree bucket: a packed per-level RGB key plus count-weighted channel sums.
type octNode struct {
	key        uint32
	r, g, b, a float64
	count      int
}

func (n *octNode) mean() golor.Color {
	c := float64(n.count)
	return golor.RGBAf(n.r/c, n.g/c, n.b/c, n.a/c)
}

func (n *octNode) add(o *octNode) {
	n.r += o.r
	n.g += o.g
	n.b += o.b
	n.a += o.a
	n.count += o.count
}

// extractOctree coarsens channel-bit buckets, then merges the finer level to reach n.
func extractOctree(swatches []swatch, n int) []golor.Color {
	// The finest level consists of the distinct colors.
	previous := make([]octNode, len(swatches))
	for i, s := range swatches {
		w := float64(s.count)
		previous[i] = octNode{
			key: s.key, count: s.count,
			r: s.color.R * w, g: s.color.G * w, b: s.color.B * w, a: s.color.A * w,
		}
	}
	parent := func(k uint32) uint32 { // Map the full RGBA key to the 7-bit level key.
		return (k>>25)<<16 | ((k>>17)&0x7f)<<8 | (k>>9)&0x7f
	}
	for bits := 7; bits >= 0; bits-- {
		level := groupNodes(previous, parent)
		switch {
		case len(level) == n:
			return nodeColors(level)
		case len(level) < n:
			return nodeColors(reduceNodes(previous, n))
		case bits == 0:
			return nodeColors(level)
		}
		previous = level
		parent = func(k uint32) uint32 { // Map the level key to the next coarser level key.
			return ((k>>16)>>1)<<16 | (((k>>8)&0xff)>>1)<<8 | (k&0xff)>>1
		}
	}
	return nil
}

// groupNodes merges nodes that share a parent key, returning the parents in key order.
func groupNodes(nodes []octNode, parent func(uint32) uint32) []octNode {
	pairs := make([]uint64, len(nodes))
	for i := range nodes {
		pairs[i] = uint64(parent(nodes[i].key))<<32 | uint64(i)
	}
	radixSortHigh(pairs)
	distinct := 1
	for i := 1; i < len(pairs); i++ {
		if pairs[i]>>32 != pairs[i-1]>>32 {
			distinct++
		}
	}
	out := make([]octNode, 0, distinct)
	for i := 0; i < len(pairs); {
		key := uint32(pairs[i] >> 32)
		node := octNode{key: key}
		for ; i < len(pairs) && pairs[i]>>32 == uint64(key); i++ {
			node.add(&nodes[uint32(pairs[i])])
		}
		out = append(out, node)
	}
	return out
}

// radixSortHigh stably sorts pairs by their high 32 bits.
func radixSortHigh(pairs []uint64) {
	if len(pairs) < 2 {
		return
	}
	buf := make([]uint64, len(pairs))
	src, dst := pairs, buf
	for shift := uint(32); shift < 64; shift += 8 {
		var counts [256]int
		for _, p := range src {
			counts[(p>>shift)&0xff]++
		}
		if counts[(src[0]>>shift)&0xff] == len(src) {
			continue
		}
		pos := 0
		for i, c := range &counts {
			counts[i] = pos
			pos += c
		}
		for _, p := range src {
			b := (p >> shift) & 0xff
			dst[counts[b]] = p
			counts[b]++
		}
		src, dst = dst, src
	}
	if &src[0] != &pairs[0] {
		copy(pairs, src)
	}
}

// reduceNodes merges nearest RGB means with cached neighbors (O(B) per merge).
// Lowest-index tie breaks keep results deterministic.
func reduceNodes(nodes []octNode, n int) []octNode {
	b := len(nodes)
	if b <= n {
		return nodes
	}
	nodes = append([]octNode(nil), nodes...)
	mr, mg, mb := make([]float64, b), make([]float64, b), make([]float64, b)
	setMean := func(i int) {
		c := nodes[i].mean()
		mr[i], mg[i], mb[i] = c.R, c.G, c.B
	}
	dist := func(i, j int) float64 {
		dr, dg, db := mr[i]-mr[j], mg[i]-mg[j], mb[i]-mb[j]
		return dr*dr + dg*dg + db*db
	}
	live := make([]int, b) // Indexes of unmerged buckets in ascending order.
	nn := make([]int, b)
	nd := make([]float64, b)
	for i := range nodes {
		setMean(i)
		live[i] = i
		nn[i] = -1
		nd[i] = math.Inf(1)
	}
	for i := 0; i < b; i++ {
		for j := i + 1; j < b; j++ {
			d := dist(i, j)
			if d < nd[i] {
				nd[i], nn[i] = d, j
			}
			if d < nd[j] {
				nd[j], nn[j] = d, i
			}
		}
	}
	nearest := func(i int) {
		nn[i], nd[i] = -1, math.Inf(1)
		for _, j := range live {
			if j == i {
				continue
			}
			if d := dist(i, j); d < nd[i] {
				nd[i], nn[i] = d, j
			}
		}
	}
	a := 0
	for _, i := range live {
		if nd[i] < nd[a] {
			a = i
		}
	}
	for len(live) > n {
		lo, hi := min(a, nn[a]), max(a, nn[a])
		nodes[lo].add(&nodes[hi])
		setMean(lo)
		pos, _ := slices.BinarySearch(live, hi)
		live = slices.Delete(live, pos, pos+1)
		nearest(lo)
		a = -1
		for _, i := range live {
			if i != lo {
				if nn[i] == lo || nn[i] == hi {
					nearest(i)
				} else if d := dist(i, lo); d < nd[i] || (d == nd[i] && lo < nn[i]) {
					nd[i], nn[i] = d, lo
				}
			}
			if a < 0 || nd[i] < nd[a] {
				a = i
			}
		}
	}
	out := make([]octNode, 0, n)
	for _, i := range live {
		out = append(out, nodes[i])
	}
	return out
}

func nodeColors(nodes []octNode) []golor.Color {
	swatches := make([]swatch, len(nodes))
	for i := range nodes {
		swatches[i] = newSwatch(nodes[i].mean(), nodes[i].count)
	}
	slices.SortFunc(swatches, compareSwatches)
	return swatchColors(swatches)
}
