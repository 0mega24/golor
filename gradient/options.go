package gradient

import (
	"fmt"
	"math"
	"sort"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/internal/space"
)

// Space selects the coordinates used to interpolate colors. The zero value is sRGB.
type Space uint8

const (
	// SpaceRGB interpolates encoded sRGB channels.
	SpaceRGB Space = iota
	// SpaceLinearRGB interpolates linear-light sRGB channels.
	SpaceLinearRGB
	// SpaceHSL interpolates hue, saturation, and lightness.
	SpaceHSL
	// SpaceHSV interpolates hue, saturation, and value.
	SpaceHSV
	// SpaceLAB interpolates CIE L*a*b* coordinates.
	SpaceLAB
	// SpaceLCH interpolates CIE lightness, chroma, and hue.
	SpaceLCH
	// SpaceOKLab interpolates OKLab coordinates.
	SpaceOKLab
	// SpaceOKLCH interpolates OKLab lightness, chroma, and hue.
	SpaceOKLCH
)

// HueRoute selects the angular path in HSL, HSV, LCH, and OKLCH.
type HueRoute uint8

const (
	// HueShorter takes the shorter arc (the default).
	HueShorter HueRoute = iota
	// HueLonger takes the longer arc; equal hues make a full increasing turn.
	HueLonger
	// HueIncreasing follows increasing angles through the hue wheel.
	HueIncreasing
	// HueDecreasing follows decreasing angles through the hue wheel.
	HueDecreasing
)

// Curve maps segment position to progress. Results are clamped to [0,1];
// non-finite results use the input. Exact stops are preserved.
type Curve func(float64) float64

// Linear returns its input unchanged.
func Linear(t float64) float64 { return t }

// EaseIn accelerates quadratically.
func EaseIn(t float64) float64 { return t * t }

// EaseOut decelerates quadratically.
func EaseOut(t float64) float64 { return t * (2 - t) }

// EaseInOut eases both ends using smoothstep.
func EaseInOut(t float64) float64 { return t * t * (3 - 2*t) }

// Logarithmic returns log1p(strength*t)/log1p(strength); positive strength speeds the start.
// Strength in (-1,0) slows the start; zero, <= -1, or non-finite strength uses Linear.
func Logarithmic(strength float64) Curve {
	if !finite(strength) || strength <= -1 || strength == 0 {
		return Linear
	}
	denominator := math.Log1p(strength)
	return func(t float64) float64 { return math.Log1p(strength*t) / denominator }
}

// ChromaCurve adjusts alpha-aware LCH/OKLCH chroma between stops, before gamut mapping.
// t is eased progress; negative results clamp to zero, non-finite results retain chroma.
type ChromaCurve func(t, chroma float64) float64

// Options controls Generate and Between; defaults are sRGB, linear progress, and shorter hue.
// Chroma controls require LCH/OKLCH; MinChroma follows Chroma and may be reduced by gamut mapping.
type Options struct {
	Space     Space
	Hue       HueRoute
	Easing    Curve
	Chroma    ChromaCurve
	MinChroma float64
}

// Stop sets a color at a nondecreasing position in [0,1]; the last duplicate wins.
// Colors outside the stops' span are held at the nearest stop.
type Stop struct {
	Color    golor.Color
	Position float64
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func (o Options) validate() error {
	if o.Space > SpaceOKLCH {
		return fmt.Errorf("gradient: invalid color space %d", o.Space)
	}
	if o.Hue > HueDecreasing {
		return fmt.Errorf("gradient: invalid hue route %d", o.Hue)
	}
	if !finite(o.MinChroma) || o.MinChroma < 0 {
		return fmt.Errorf("gradient: minimum chroma must be finite and nonnegative")
	}
	if (o.Chroma != nil || o.MinChroma != 0) && o.Space != SpaceLCH && o.Space != SpaceOKLCH {
		return fmt.Errorf("gradient: chroma controls require LCH or OKLCH")
	}
	return nil
}

func (o Options) newMixer(a, b golor.Color) mixer {
	switch o.Space {
	case SpaceRGB:
		return rgbMixer(a, b)
	case SpaceLinearRGB:
		return linearRGBMixer(a, b)
	case SpaceHSL:
		return hslMixerWithOptions(a, b, o)
	case SpaceHSV:
		return hsvMixerWithOptions(a, b, o)
	case SpaceLAB:
		return labMixer(a, b)
	case SpaceLCH:
		return lchMixerWithOptions(a, b, o)
	case SpaceOKLab:
		return oklabMixer(a, b)
	case SpaceOKLCH:
		return oklchMixerWithOptions(a, b, o)
	}
	return nil // Generate validates the space before building mixers.
}

// Between samples n colors using Generate with endpoints at positions 0 and 1.
func Between(a, b golor.Color, n int, options Options) ([]golor.Color, error) {
	return Generate([]Stop{{Color: a, Position: 0}, {Color: b, Position: 1}}, n, options)
}

// Generate samples [0,1], easing each segment while preserving stops and the input slice.
// No stops or n <= 0 returns nil; n == 1 samples 0. Colors are clamped to sRGB/alpha.
// Options and positions are validated even when no samples are requested.
func Generate(stops []Stop, n int, options Options) ([]golor.Color, error) {
	if err := options.validate(); err != nil {
		return nil, err
	}
	for i, s := range stops {
		if !finite(s.Position) || s.Position < 0 || s.Position > 1 {
			return nil, fmt.Errorf("gradient: stop %d position must be finite and in [0,1]", i)
		}
		if i > 0 && s.Position < stops[i-1].Position {
			return nil, fmt.Errorf("gradient: stop positions must be nondecreasing")
		}
	}
	if n <= 0 || len(stops) == 0 {
		return nil, nil
	}
	colors := make([]golor.Color, len(stops))
	for i, s := range stops {
		colors[i] = golor.RGBAf(s.Color.R, s.Color.G, s.Color.B, s.Color.A)
	}
	segs := make([]mixer, len(stops)-1)
	for i := range segs {
		segs[i] = options.newMixer(colors[i], colors[i+1])
	}
	out := make([]golor.Color, n)
	for i := range out {
		position := 0.0
		if n > 1 {
			position = float64(i) / float64(n-1)
		}
		// First strictly greater stop: the previous stop is the rightmost match.
		next := sort.Search(len(stops), func(j int) bool { return stops[j].Position > position })
		switch {
		case next == 0:
			out[i] = colors[0]
		case next == len(stops):
			out[i] = colors[len(colors)-1]
		case position == stops[next-1].Position:
			out[i] = colors[next-1]
		default:
			t := (position - stops[next-1].Position) / (stops[next].Position - stops[next-1].Position)
			if options.Easing != nil {
				if eased := options.Easing(t); finite(eased) {
					t = space.Clamp01(eased)
				}
			}
			out[i] = segs[next-1](t)
		}
	}
	return out, nil
}
