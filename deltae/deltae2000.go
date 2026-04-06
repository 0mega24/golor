package deltae

import (
	"math"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// DeltaE2000 returns the CIE 2000 color difference (CIEDE2000 formula, Sharma et al. 2005).
func DeltaE2000(c1, c2 golor.Color) float64 {
	return ciede2000(convert.ToLAB(c1), convert.ToLAB(c2))
}

func ciede2000(lab1, lab2 convert.LAB) float64 {
	kL, kC, kH := 1.0, 1.0, 1.0

	c1 := math.Sqrt(lab1.A*lab1.A + lab1.B*lab1.B)
	c2 := math.Sqrt(lab2.A*lab2.A + lab2.B*lab2.B)
	cAvg := (c1 + c2) / 2
	cAvg7 := math.Pow(cAvg, 7)
	g := 0.5 * (1 - math.Sqrt(cAvg7/(cAvg7+math.Pow(25, 7))))

	a1p := lab1.A * (1 + g)
	a2p := lab2.A * (1 + g)

	c1p := math.Sqrt(a1p*a1p + lab1.B*lab1.B)
	c2p := math.Sqrt(a2p*a2p + lab2.B*lab2.B)

	h1p := math.Atan2(lab1.B, a1p) * 180 / math.Pi
	if h1p < 0 {
		h1p += 360
	}
	h2p := math.Atan2(lab2.B, a2p) * 180 / math.Pi
	if h2p < 0 {
		h2p += 360
	}

	dLp := lab2.L - lab1.L
	dCp := c2p - c1p

	var dhp float64
	if c1p*c2p == 0 {
		dhp = 0
	} else if math.Abs(h2p-h1p) <= 180 {
		dhp = h2p - h1p
	} else if h2p-h1p > 180 {
		dhp = h2p - h1p - 360
	} else {
		dhp = h2p - h1p + 360
	}
	dHp := 2 * math.Sqrt(c1p*c2p) * math.Sin(dhp*math.Pi/360)

	lAvgp := (lab1.L + lab2.L) / 2
	cAvgp := (c1p + c2p) / 2

	var hAvgp float64
	if c1p*c2p == 0 {
		hAvgp = h1p + h2p
	} else if math.Abs(h1p-h2p) <= 180 {
		hAvgp = (h1p + h2p) / 2
	} else if h1p+h2p < 360 {
		hAvgp = (h1p + h2p + 360) / 2
	} else {
		hAvgp = (h1p + h2p - 360) / 2
	}

	t := 1 -
		0.17*math.Cos((hAvgp-30)*math.Pi/180) +
		0.24*math.Cos(2*hAvgp*math.Pi/180) +
		0.32*math.Cos((3*hAvgp+6)*math.Pi/180) -
		0.20*math.Cos((4*hAvgp-63)*math.Pi/180)

	lAvgpSq := (lAvgp - 50) * (lAvgp - 50)
	sL := 1 + 0.015*lAvgpSq/math.Sqrt(20+lAvgpSq)
	sC := 1 + 0.045*cAvgp
	sH := 1 + 0.015*cAvgp*t

	cAvgp7 := math.Pow(cAvgp, 7)
	rc := 2 * math.Sqrt(cAvgp7/(cAvgp7+math.Pow(25, 7)))
	dTheta := 30 * math.Exp(-math.Pow((hAvgp-275)/25, 2))
	rT := -math.Sin(2*dTheta*math.Pi/180) * rc

	return math.Sqrt(
		math.Pow(dLp/(kL*sL), 2)+
			math.Pow(dCp/(kC*sC), 2)+
			math.Pow(dHp/(kH*sH), 2)+
			rT*(dCp/(kC*sC))*(dHp/(kH*sH)),
	)
}
