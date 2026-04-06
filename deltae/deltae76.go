package deltae

import (
	"math"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// DeltaE76 returns the CIE 1976 color difference (Euclidean distance in L*a*b* space).
func DeltaE76(c1, c2 golor.Color) float64 {
	l1 := convert.ToLAB(c1)
	l2 := convert.ToLAB(c2)
	dL := l1.L - l2.L
	da := l1.A - l2.A
	db := l1.B - l2.B
	return math.Sqrt(dL*dL + da*da + db*db)
}
