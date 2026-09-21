// Package palette extracts dominant colors from images.
//
// MedianCut and Octree are deterministic, single-pass-ish quantizers suitable
// for large images. KMeans uses deterministic k-means++ initialization with a
// caller-controlled seed and a bounded iteration count.
package palette

import (
	"image"
	"math"
	"math/rand"
	"sort"

	"github.com/0mega24/golor/v2"
)

// Algorithm identifies the palette extraction algorithm to use.
type Algorithm int

const (
	// MedianCut recursively splits color buckets on their widest RGB dimension.
	MedianCut Algorithm = iota
	// KMeans clusters colors with k-means++ initialization and deterministic seeded randomness.
	KMeans
	// Octree quantizes colors through RGB bit buckets, similar to classic octree reduction.
	Octree
)

const (
	// DefaultSeed is used by KMeans when no seed option is provided.
	DefaultSeed int64 = 1
	// DefaultMaxIterations bounds KMeans convergence work.
	DefaultMaxIterations = 32
)

// DefaultAlgorithm is the algorithm used by Extract when no algorithm option is provided.
const DefaultAlgorithm = MedianCut

// Option configures palette extraction.
type Option func(*options)

type options struct {
	algorithm     Algorithm
	seed          int64
	maxIterations int
}

// WithAlgorithm selects the extraction algorithm.
func WithAlgorithm(a Algorithm) Option {
	return func(o *options) {
		o.algorithm = a
	}
}

// WithSeed sets the deterministic seed used by KMeans.
func WithSeed(seed int64) Option {
	return func(o *options) {
		o.seed = seed
	}
}

// WithMaxIterations sets the KMeans iteration cap.
func WithMaxIterations(n int) Option {
	return func(o *options) {
		o.maxIterations = n
	}
}

// Extract returns up to n dominant colors from img, ordered most-dominant first.
// It returns nil for nil images, empty images, or n <= 0.
func Extract(img image.Image, n int, opts ...Option) []golor.Color {
	cfg := options{
		algorithm:     DefaultAlgorithm,
		seed:          DefaultSeed,
		maxIterations: DefaultMaxIterations,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.maxIterations <= 0 {
		cfg.maxIterations = DefaultMaxIterations
	}

	swatches := imageSwatches(img)
	if n <= 0 || len(swatches) == 0 {
		return nil
	}
	sortSwatches(swatches)
	if n >= len(swatches) {
		return swatchColors(swatches)
	}

	switch cfg.algorithm {
	case KMeans:
		return extractKMeans(swatches, n, cfg.seed, cfg.maxIterations)
	case Octree:
		return extractOctree(swatches, n)
	case MedianCut:
		return extractMedianCut(swatches, n)
	default:
		return extractMedianCut(swatches, n)
	}
}

type swatch struct {
	color golor.Color
	count int
}

func imageSwatches(img image.Image) []swatch {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		return nil
	}
	counts := map[uint32]int{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := golor.FromStdColor(img.At(x, y))
			key := colorKey(c)
			counts[key]++
		}
	}
	swatches := make([]swatch, 0, len(counts))
	for key, count := range counts {
		swatches = append(swatches, swatch{color: keyColor(key), count: count})
	}
	return swatches
}

func colorKey(c golor.Color) uint32 {
	return uint32(c.R8())<<24 | uint32(c.G8())<<16 | uint32(c.B8())<<8 | uint32(c.A8())
}

func keyColor(key uint32) golor.Color {
	return golor.RGBA(uint8(key>>24), uint8(key>>16), uint8(key>>8), uint8(key))
}

func sortSwatches(swatches []swatch) {
	sort.Slice(swatches, func(i, j int) bool {
		if swatches[i].count == swatches[j].count {
			return colorKey(swatches[i].color) < colorKey(swatches[j].color)
		}
		return swatches[i].count > swatches[j].count
	})
}

func swatchColors(swatches []swatch) []golor.Color {
	colors := make([]golor.Color, len(swatches))
	for i, s := range swatches {
		colors[i] = s.color
	}
	return colors
}

func weightedAverage(swatches []swatch) golor.Color {
	var r, g, b, a float64
	var total int
	for _, s := range swatches {
		weight := float64(s.count)
		r += s.color.R * weight
		g += s.color.G * weight
		b += s.color.B * weight
		a += s.color.A * weight
		total += s.count
	}
	if total == 0 {
		return golor.RGBAf(0, 0, 0, 0)
	}
	return golor.RGBAf(r/float64(total), g/float64(total), b/float64(total), a/float64(total))
}

func totalCount(swatches []swatch) int {
	total := 0
	for _, s := range swatches {
		total += s.count
	}
	return total
}

func rgbDistanceSquared(a, b golor.Color) float64 {
	dr := a.R - b.R
	dg := a.G - b.G
	db := a.B - b.B
	return dr*dr + dg*dg + db*db
}

func nearestCentroid(c golor.Color, centroids []golor.Color) int {
	best := 0
	bestDistance := math.Inf(1)
	for i, centroid := range centroids {
		distance := rgbDistanceSquared(c, centroid)
		if distance < bestDistance {
			bestDistance = distance
			best = i
		}
	}
	return best
}

func seededRand(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed)) // #nosec G404 -- deterministic k-means seeding, not security randomness.
}
