package contrast

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// The APCA constants below are from APCA-W3 0.0.98G-4g (https://github.com/Myndex/apca-w3).
const (
	apcaMainTRC   = 2.4
	apcaRco       = 0.2126729
	apcaGco       = 0.7151522
	apcaBco       = 0.0721750
	apcaNormBG    = 0.56
	apcaNormTXT   = 0.57
	apcaRevTXT    = 0.62
	apcaRevBG     = 0.65
	apcaBlkThrs   = 0.022
	apcaBlkClmp   = 1.414
	apcaScale     = 1.14
	apcaLoOffset  = 0.027
	apcaDeltaYMin = 0.0005
	apcaLoClip    = 0.1
	apcaLcPerSAPC = 100
)

// APCA returns signed APCA-W3 Lc: positive for dark text, negative for light text.
// Alpha compositing follows Ratio. Thresholds depend on font size and weight.
func APCA(text, bg golor.Color) float64 {
	fg, back := resolvePair(text, bg)
	yt, yb := apcaY(fg), apcaY(back)
	if math.Abs(yb-yt) < apcaDeltaYMin {
		return 0
	}
	var lc float64
	if yb > yt { // Dark text on a light background.
		sapc := (math.Pow(yb, apcaNormBG) - math.Pow(yt, apcaNormTXT)) * apcaScale
		if sapc >= apcaLoClip {
			lc = sapc - apcaLoOffset
		}
	} else { // Light text on a dark background.
		sapc := (math.Pow(yb, apcaRevBG) - math.Pow(yt, apcaRevTXT)) * apcaScale
		if sapc <= -apcaLoClip {
			lc = sapc + apcaLoOffset
		}
	}
	return lc * apcaLcPerSAPC
}

// apcaY uses a 2.4 power curve with a near-black soft clamp.
func apcaY(c golor.Color) float64 {
	y := apcaRco*math.Pow(space.Clamp01(c.R), apcaMainTRC) +
		apcaGco*math.Pow(space.Clamp01(c.G), apcaMainTRC) +
		apcaBco*math.Pow(space.Clamp01(c.B), apcaMainTRC)
	if y < apcaBlkThrs {
		y += math.Pow(apcaBlkThrs-y, apcaBlkClmp)
	}
	return y
}
