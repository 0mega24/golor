package palette

import (
	"math"
	"math/rand"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/internal/space"
)

const (
	// binThreshold enables 5-bit-per-channel bins with count-weighted means.
	binThreshold = 1 << 14
	// moveEpsilon stops centroid updates below 0.001 OKLab distance.
	moveEpsilon = 1e-3
)

// points is the k-means input in struct-of-arrays form: OKLab coordinates, pixel
// counts as weights, and summed alpha for the output color.
type points struct {
	l, a, b  []float64
	weight   []float64
	alphaSum []float64
	count    []int
	first    int // Index of the most frequent point; ties use the lowest key.
}

func newPoints(swatches []swatch) points {
	if len(swatches) > binThreshold {
		return binnedPoints(swatches)
	}
	p := makePoints(len(swatches))
	best := -1
	for i, s := range swatches {
		p.set(i, s.color.R, s.color.G, s.color.B, s.count, s.color.A*float64(s.count))
		if best < 0 || s.count > swatches[best].count {
			best = i // Swatches are in key order, so the first maximum has the lowest key.
		}
	}
	p.first = best
	return p
}

func makePoints(n int) points {
	return points{
		l: make([]float64, n), a: make([]float64, n), b: make([]float64, n),
		weight: make([]float64, n), alphaSum: make([]float64, n), count: make([]int, n),
	}
}

func (p *points) set(i int, r, g, b float64, count int, alphaSum float64) {
	lab := space.LinearSRGBToOKLab(space.LinearizeVec(space.Vec3{r, g, b}))
	p.l[i], p.a[i], p.b[i] = lab[0], lab[1], lab[2]
	p.weight[i] = float64(count)
	p.alphaSum[i] = alphaSum
	p.count[i] = count
}

func binnedPoints(swatches []swatch) points {
	type bin struct {
		r, g, b, a float64
		count      int
	}
	bins := make([]bin, 1<<15)
	for _, s := range swatches {
		k := s.key
		idx := (k>>27)<<10 | ((k>>19)&0x1f)<<5 | (k>>11)&0x1f
		w := float64(s.count)
		bn := &bins[idx]
		bn.r += s.color.R * w
		bn.g += s.color.G * w
		bn.b += s.color.B * w
		bn.a += s.color.A * w
		bn.count += s.count
	}
	used := 0
	for i := range bins {
		if bins[i].count > 0 {
			used++
		}
	}
	p := makePoints(used)
	p.first = -1
	i := 0
	for _, bn := range bins {
		if bn.count == 0 {
			continue
		}
		c := float64(bn.count)
		p.set(i, bn.r/c, bn.g/c, bn.b/c, bn.count, bn.a)
		if p.first < 0 || bn.count > p.count[p.first] {
			p.first = i
		}
		i++
	}
	return p
}

// chunkSize is the fixed unit of parallel work. Each chunk accumulates its own centroid
// deltas, which are reduced in chunk order, so results are identical for any GOMAXPROCS.
const chunkSize = 2048

// kmeans holds the Lloyd/Hamerly state in struct-of-arrays form.
type kmeans struct {
	p                 *points
	cl, ca, cb        []float64 // Centroids.
	assign            []int
	upper, lower      []float64 // Hamerly bounds: distances to the assigned centroid and the second closest.
	sums              []float64 // Per cluster: l, a, b, weight.
	half, moved       []float64
	maxMove, maxMove2 float64
	maxIdx            int
}

// chunkResult is one chunk's contribution: centroid sum deltas and reassignment count.
type chunkResult struct {
	delta   []float64
	changed int
}

func (km *kmeans) fullScan(i int) (best int, d1, d2 float64) {
	d1, d2 = math.Inf(1), math.Inf(1)
	pl, pa, pb := km.p.l[i], km.p.a[i], km.p.b[i]
	for j := range km.cl {
		dl, da, db := pl-km.cl[j], pa-km.ca[j], pb-km.cb[j]
		d := dl*dl + da*da + db*db
		if d < d1 {
			d2, d1, best = d1, d, j
		} else if d < d2 {
			d2 = d
		}
	}
	return best, math.Sqrt(d1), math.Sqrt(d2)
}

func (km *kmeans) add(delta []float64, i, j int, sign float64) {
	w := km.p.weight[i] * sign
	delta[4*j] += km.p.l[i] * w
	delta[4*j+1] += km.p.a[i] * w
	delta[4*j+2] += km.p.b[i] * w
	delta[4*j+3] += w
}

// assignChunk assigns points [lo, hi). On the first pass every point gets a full scan;
// afterwards Hamerly's bounds skip points that cannot have changed cluster.
func (km *kmeans) assignChunk(lo, hi int, first bool, res *chunkResult) {
	for j := range res.delta {
		res.delta[j] = 0
	}
	res.changed = 0
	for i := lo; i < hi; i++ {
		if first {
			best, d1, d2 := km.fullScan(i)
			km.assign[i], km.upper[i], km.lower[i] = best, d1, d2
			km.add(res.delta, i, best, 1)
			continue
		}
		a := km.assign[i]
		km.upper[i] += km.moved[a]
		if a == km.maxIdx {
			km.lower[i] -= km.maxMove2
		} else {
			km.lower[i] -= km.maxMove
		}
		bound := math.Max(km.half[a], km.lower[i])
		if km.upper[i] <= bound {
			continue
		}
		dl, da, db := km.p.l[i]-km.cl[a], km.p.a[i]-km.ca[a], km.p.b[i]-km.cb[a]
		km.upper[i] = math.Sqrt(dl*dl + da*da + db*db)
		if km.upper[i] <= bound {
			continue
		}
		best, d1, d2 := km.fullScan(i)
		km.upper[i], km.lower[i] = d1, d2
		if best != a {
			km.add(res.delta, i, a, -1)
			km.add(res.delta, i, best, 1)
			km.assign[i] = best
			res.changed++
		}
	}
}

// assignAll runs assignChunk over every chunk (in parallel when there is enough work),
// then adds the chunk deltas to the cluster sums in chunk order.
func (km *kmeans) assignAll(results []chunkResult, first bool) int {
	m := len(km.assign)
	chunks := len(results)
	workers := min(runtime.GOMAXPROCS(0), chunks)
	if m*len(km.cl) < 1<<16 {
		workers = 1
	}
	run := func(c int) { km.assignChunk(c*chunkSize, min(m, (c+1)*chunkSize), first, &results[c]) }
	if workers <= 1 {
		for c := range results {
			run(c)
		}
	} else {
		var next atomic.Int64
		var wg sync.WaitGroup
		for range workers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for c := int(next.Add(1) - 1); c < chunks; c = int(next.Add(1) - 1) {
					run(c)
				}
			}()
		}
		wg.Wait()
	}
	changed := 0
	for c := range results {
		for j, d := range results[c].delta {
			km.sums[j] += d
		}
		changed += results[c].changed
	}
	return changed
}

// updateCentroids moves centroids to their cluster means and refreshes the per-centroid
// movement and half-distance tables used by the bounds.
func (km *kmeans) updateCentroids() {
	k := len(km.cl)
	km.maxMove, km.maxMove2, km.maxIdx = 0, 0, -1
	for j := 0; j < k; j++ {
		km.moved[j] = 0
		w := km.sums[4*j+3]
		if w <= 0 {
			continue
		}
		nl, na, nb := km.sums[4*j]/w, km.sums[4*j+1]/w, km.sums[4*j+2]/w
		dl, da, db := nl-km.cl[j], na-km.ca[j], nb-km.cb[j]
		km.moved[j] = math.Sqrt(dl*dl + da*da + db*db)
		km.cl[j], km.ca[j], km.cb[j] = nl, na, nb
		if km.moved[j] > km.maxMove {
			km.maxMove2, km.maxMove, km.maxIdx = km.maxMove, km.moved[j], j
		} else if km.moved[j] > km.maxMove2 {
			km.maxMove2 = km.moved[j]
		}
	}
	for j := 0; j < k; j++ {
		nearest := math.Inf(1)
		for o := 0; o < k; o++ {
			if o == j {
				continue
			}
			dl, da, db := km.cl[j]-km.cl[o], km.ca[j]-km.ca[o], km.cb[j]-km.cb[o]
			nearest = math.Min(nearest, dl*dl+da*da+db*db)
		}
		km.half[j] = math.Sqrt(nearest) / 2
	}
}

// extractKMeans uses OKLab k-means++, Hamerly bounds, and deterministic parallel chunks.
func extractKMeans(swatches []swatch, n int, seed int64, maxIterations int) []golor.Color {
	p := newPoints(swatches)
	m := len(p.l)
	cl, ca, cb := initialCentroids(&p, min(n, m), seed)
	k := len(cl)
	km := &kmeans{
		p: &p, cl: cl, ca: ca, cb: cb,
		assign: make([]int, m), upper: make([]float64, m), lower: make([]float64, m),
		sums: make([]float64, 4*k), half: make([]float64, k), moved: make([]float64, k),
	}
	results := make([]chunkResult, (m+chunkSize-1)/chunkSize)
	for c := range results {
		results[c].delta = make([]float64, 4*k)
	}
	km.assignAll(results, true)
	for iter := 1; iter < maxIterations; iter++ {
		km.updateCentroids()
		if km.maxMove < moveEpsilon {
			break
		}
		if km.assignAll(results, false) == 0 {
			break
		}
	}
	return clusterColors(&p, km.assign, k)
}

// initialCentroids starts with the most frequent color, then uses weighted D² sampling.
func initialCentroids(p *points, n int, seed int64) (cl, ca, cb []float64) {
	rng := rand.New(rand.NewSource(seed)) // #nosec G404 -- deterministic k-means seeding, not security randomness.
	m := len(p.l)
	minDist := make([]float64, m)
	for i := range minDist {
		minDist[i] = math.Inf(1)
	}
	chosen := p.first
	for {
		cl, ca, cb = append(cl, p.l[chosen]), append(ca, p.a[chosen]), append(cb, p.b[chosen])
		if len(cl) == n {
			return cl, ca, cb
		}
		total := 0.0
		for i := 0; i < m; i++ {
			dl, da, db := p.l[i]-p.l[chosen], p.a[i]-p.a[chosen], p.b[i]-p.b[chosen]
			if d := dl*dl + da*da + db*db; d < minDist[i] {
				minDist[i] = d
			}
			total += minDist[i] * p.weight[i]
		}
		if total == 0 {
			return cl, ca, cb
		}
		target := rng.Float64() * total
		running := 0.0
		chosen = -1
		for i := 0; i < m; i++ {
			if minDist[i] == 0 {
				continue // Already a centroid or a duplicate of one.
			}
			chosen = i
			running += minDist[i] * p.weight[i]
			if running >= target {
				break
			}
		}
	}
}

// clusterColors computes each non-empty cluster's count-weighted OKLab mean, converts the
// mean to sRGB, and orders clusters by decreasing frequency.
func clusterColors(p *points, assign []int, k int) []golor.Color {
	type acc struct{ l, a, b, w, alpha float64 }
	accs := make([]acc, k)
	counts := make([]int, k)
	for i, j := range assign {
		w := p.weight[i]
		accs[j].l += p.l[i] * w
		accs[j].a += p.a[i] * w
		accs[j].b += p.b[i] * w
		accs[j].w += w
		accs[j].alpha += p.alphaSum[i]
		counts[j] += p.count[i]
	}
	clusters := make([]swatch, 0, k)
	for j, c := range accs {
		if counts[j] == 0 {
			continue
		}
		col := convert.FromOKLab(convert.OKLab{L: c.l / c.w, A: c.a / c.w, B: c.b / c.w})
		col.A = space.Clamp01(c.alpha / c.w)
		if math.IsNaN(col.R) {
			continue
		}
		clusters = append(clusters, newSwatch(col, counts[j]))
	}
	slices.SortFunc(clusters, compareSwatches)
	return swatchColors(clusters)
}
