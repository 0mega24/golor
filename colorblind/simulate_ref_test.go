package colorblind_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/colorblind"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/deltae"
	"github.com/stretchr/testify/assert"
)

// TestSimulateReferenceValues checks values computed from the precomputed libDaltonLens matrices.
func TestSimulateReferenceValues(t *testing.T) {
	cases := []struct {
		in   string
		d    colorblind.Deficiency
		want string
	}{
		// Dichromats along the blue-yellow axis are unaffected by protan/deutan.
		{"#0000ff", colorblind.Deuteranopia, "#0000ff"},
		{"#0000ff", colorblind.Protanopia, "#0000ff"},
		{"#ffffff", colorblind.Deuteranopia, "#ffffff"},
		{"#ffffff", colorblind.Tritanopia, "#ffffff"},
		{"#000000", colorblind.Tritanopia, "#000000"},
	}
	for _, tc := range cases {
		in, _ := golor.Hex(tc.in)
		assert.Equal(t, tc.want, colorblind.Simulate(in, tc.d).String(), "%s under %d", tc.in, tc.d)
	}
	// Red and green collapse onto the same yellowish axis for deuteranopes.
	r := colorblind.Simulate(golor.RGB(255, 0, 0), colorblind.Deuteranopia)
	g := colorblind.Simulate(golor.RGB(0, 255, 0), colorblind.Deuteranopia)
	assert.InDelta(t, r.R, r.G, 1e-9)
	assert.InDelta(t, g.R, g.G, 1e-9)
	assert.Less(t, r.B, 0.05)
	// Tritanopes confuse blue and green: both lose their blue/yellow distinction.
	b := colorblind.Simulate(golor.RGB(0, 0, 255), colorblind.Tritanopia)
	assert.Greater(t, b.G, 0.3)
}

func TestSimulateUnknownDeficiencyIsIdentity(t *testing.T) {
	c := golor.RGBA(200, 100, 50, 77)
	assert.Equal(t, c, colorblind.Simulate(c, colorblind.Deficiency(7)))
	assert.Equal(t, c, colorblind.Simulate(c, colorblind.Deficiency(-1)))
}

// TestAccessiblePaletteMatchesNaive checks that caching simulated Lab values preserves the
// decisions of the implementation that re-simulates every prior color.
func TestAccessiblePaletteMatchesNaive(t *testing.T) {
	r := rand.New(rand.NewSource(9))
	for trial := 0; trial < 200; trial++ {
		n := 2 + r.Intn(10)
		in := make([]golor.Color, n)
		for i := range in {
			// Cluster colors so that many require hue shifts.
			in[i] = golor.RGBA(uint8(100+r.Intn(40)), uint8(100+r.Intn(40)), uint8(r.Intn(256)), uint8(128+r.Intn(128)))
		}
		for _, d := range []colorblind.Deficiency{colorblind.Deuteranopia, colorblind.Protanopia, colorblind.Tritanopia} {
			assert.Equal(t, naiveAccessiblePalette(in, d), colorblind.AccessiblePalette(in, d))
		}
	}
}

func naiveAccessiblePalette(colors []golor.Color, d colorblind.Deficiency) []golor.Color {
	result := append([]golor.Color(nil), colors...)
	minPrior := func(i int) float64 {
		m := math.Inf(1)
		for j := 0; j < i; j++ {
			m = math.Min(m, deltae.DeltaE76(colorblind.Simulate(result[j], d), colorblind.Simulate(result[i], d)))
		}
		return m
	}
	steps := []float64{15, 30, 45, 60, 90, 120, 150, 180, 210, 240, 270, 300, 330}
	for i := 1; i < len(result); i++ {
		if minPrior(i) >= 10 {
			continue
		}
		original := result[i]
		best, bestD := original, minPrior(i)
		hsl := convert.ToHSL(original)
		chosen := false
		for a := 0; a < colorblind.MaxAccessiblePaletteAttempts; a++ {
			shift := float64((a-len(steps)+1)*37) + 15
			if a < len(steps) {
				shift = steps[a]
			}
			h := hsl
			h.H = math.Mod(hsl.H+shift, 360)
			cand := convert.FromHSL(h)
			cand.A = original.A
			result[i] = cand
			dist := minPrior(i)
			if dist > bestD {
				best, bestD = cand, dist
			}
			if dist >= 10 {
				chosen = true
				break
			}
		}
		if !chosen {
			result[i] = best
		}
	}
	return result
}
