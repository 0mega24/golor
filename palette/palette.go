// Package palette extracts dominant colors with MedianCut, Octree, or seeded OKLab KMeans.
// Transparent pixels are skipped by default; see WithMinAlpha.
package palette

import (
	"cmp"
	"image"
	"slices"

	"github.com/0mega24/golor/v2"
)

// Algorithm identifies the palette extraction algorithm to use.
type Algorithm int

const (
	// MedianCut recursively splits color buckets on their widest RGB dimension.
	MedianCut Algorithm = iota
	// KMeans clusters colors in OKLab with k-means++ initialization and deterministic seeded randomness.
	KMeans
	// Octree quantizes colors through RGB bit buckets, similar to classic octree reduction.
	Octree
)

const (
	// DefaultSeed is used by KMeans when no seed option is provided.
	DefaultSeed int64 = 1
	// DefaultMaxIterations bounds KMeans convergence work.
	DefaultMaxIterations = 32
	// DefaultMinAlpha is the default minimum 8-bit alpha for a pixel to be counted.
	// It skips only fully transparent pixels.
	DefaultMinAlpha uint8 = 1
)

// DefaultAlgorithm is the algorithm used by Extract when no algorithm option is provided.
const DefaultAlgorithm = MedianCut

// Option configures palette extraction.
type Option func(*options)

type options struct {
	algorithm     Algorithm
	seed          int64
	maxIterations int
	minAlpha      uint8
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

// WithMinAlpha skips alpha below a; the default skips transparent pixels, while 0 counts all.
func WithMinAlpha(a uint8) Option {
	return func(o *options) {
		o.minAlpha = a
	}
}

// Extract returns up to n dominant colors, most frequent first.
// Nil/empty images, no qualifying pixels, or n <= 0 return nil.
func Extract(img image.Image, n int, opts ...Option) []golor.Color {
	cfg := options{
		algorithm:     DefaultAlgorithm,
		seed:          DefaultSeed,
		maxIterations: DefaultMaxIterations,
		minAlpha:      DefaultMinAlpha,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.maxIterations <= 0 {
		cfg.maxIterations = DefaultMaxIterations
	}
	if n <= 0 {
		return nil
	}
	swatches := imageSwatches(img, cfg.minAlpha)
	if len(swatches) == 0 {
		return nil
	}
	if n >= len(swatches) {
		slices.SortFunc(swatches, compareSwatches)
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

// swatch is a distinct color and its pixel count. key is the packed 8-bit RGBA value
// of color, cached to avoid recomputing it during sorting.
type swatch struct {
	color golor.Color
	key   uint32
	count int
}

func newSwatch(c golor.Color, count int) swatch {
	return swatch{color: c, key: colorKey(c), count: count}
}

// imageSwatches counts distinct 8-bit colors in img and returns them in ascending key order.
func imageSwatches(img image.Image, minAlpha uint8) []swatch {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		return nil
	}
	keys := pixelKeys(img, bounds, minAlpha)
	if len(keys) == 0 {
		return nil
	}
	radixSort(keys)
	distinct := 1
	for i := 1; i < len(keys); i++ {
		if keys[i] != keys[i-1] {
			distinct++
		}
	}
	swatches := make([]swatch, 0, distinct)
	for i := 0; i < len(keys); {
		j := i + 1
		for j < len(keys) && keys[j] == keys[i] {
			j++
		}
		swatches = append(swatches, swatch{color: keyColor(keys[i]), key: keys[i], count: j - i})
		i = j
	}
	return swatches
}

// pixelKeys packs pixels above minAlpha; direct NRGBA/RGBA reads match FromStdColor exactly.
func pixelKeys(img image.Image, b image.Rectangle, minAlpha uint8) []uint32 {
	keys := make([]uint32, 0, b.Dx()*b.Dy())
	switch src := img.(type) {
	case *image.NRGBA:
		t := tables()
		for y := b.Min.Y; y < b.Max.Y; y++ {
			row := src.Pix[src.PixOffset(b.Min.X, y):src.PixOffset(b.Max.X, y)]
			for i := 0; i+3 < len(row); i += 4 {
				a := row[i+3]
				if t.alpha[a] < minAlpha {
					continue
				}
				lut := &t.straight[a]
				keys = append(keys, uint32(lut[row[i]])<<24|uint32(lut[row[i+1]])<<16|uint32(lut[row[i+2]])<<8|uint32(t.alpha[a]))
			}
		}
	case *image.RGBA:
		t := tables()
		for y := b.Min.Y; y < b.Max.Y; y++ {
			row := src.Pix[src.PixOffset(b.Min.X, y):src.PixOffset(b.Max.X, y)]
			for i := 0; i+3 < len(row); i += 4 {
				a := row[i+3]
				if t.alpha[a] < minAlpha {
					continue
				}
				lut := &t.premul[a]
				keys = append(keys, uint32(lut[row[i]])<<24|uint32(lut[row[i+1]])<<16|uint32(lut[row[i+2]])<<8|uint32(t.alpha[a]))
			}
		}
	default:
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				c := golor.FromStdColor(img.At(x, y))
				if c.A8() < minAlpha {
					continue
				}
				keys = append(keys, colorKey(c))
			}
		}
	}
	return keys
}

// radixSort sorts keys ascending with an LSD byte radix sort, skipping passes where
// every key shares the same byte (for example, the alpha byte of opaque images).
func radixSort(keys []uint32) {
	if len(keys) < 256 {
		slices.Sort(keys)
		return
	}
	buf := make([]uint32, len(keys))
	src, dst := keys, buf
	for shift := uint(0); shift < 32; shift += 8 {
		var counts [256]int
		for _, k := range src {
			counts[(k>>shift)&0xff]++
		}
		if counts[(src[0]>>shift)&0xff] == len(src) {
			continue
		}
		pos := 0
		for i, c := range &counts {
			counts[i] = pos
			pos += c
		}
		for _, k := range src {
			b := (k >> shift) & 0xff
			dst[counts[b]] = k
			counts[b]++
		}
		src, dst = dst, src
	}
	if &src[0] != &keys[0] {
		copy(keys, src)
	}
}

// compareSwatches orders swatches most frequent first, then by ascending key.
func compareSwatches(a, b swatch) int {
	if a.count != b.count {
		return cmp.Compare(b.count, a.count)
	}
	return cmp.Compare(a.key, b.key)
}

func colorKey(c golor.Color) uint32 {
	return uint32(c.R8())<<24 | uint32(c.G8())<<16 | uint32(c.B8())<<8 | uint32(c.A8())
}

func keyColor(key uint32) golor.Color {
	return golor.RGBA(uint8(key>>24), uint8(key>>16), uint8(key>>8), uint8(key))
}

func swatchColors(swatches []swatch) []golor.Color {
	colors := make([]golor.Color, len(swatches))
	for i, s := range swatches {
		colors[i] = s.color
	}
	return colors
}
