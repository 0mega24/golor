// Package deltae provides perceptual color difference metrics (ΔE76, CIEDE2000, and
// OKLab distance).
package deltae

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/internal/space"
)

// DeltaE76 returns the CIE 1976 color difference (Euclidean distance in L*a*b* space).
// Non-opaque colors are composited over white before comparison.
func DeltaE76(c1, c2 golor.Color) float64 {
	l1 := convert.ToLAB(compositeWhite(c1))
	l2 := convert.ToLAB(compositeWhite(c2))
	dL := l1.L - l2.L
	da := l1.A - l2.A
	db := l1.B - l2.B
	return math.Sqrt(dL*dL + da*da + db*db)
}

// DeltaEOK returns the Euclidean distance in OKLab, an inexpensive perceptual difference
// where about 0.02 is just noticeable. Non-opaque colors are composited over white.
func DeltaEOK(c1, c2 golor.Color) float64 { //nolint:revive // Named to match DeltaE76 and DeltaE2000.
	return convert.DeltaEOK(compositeWhite(c1), compositeWhite(c2))
}

func compositeWhite(c golor.Color) golor.Color {
	a := space.Clamp01(c.A)
	return golor.RGBf(
		c.R*a+1*(1-a),
		c.G*a+1*(1-a),
		c.B*a+1*(1-a),
	)
}
