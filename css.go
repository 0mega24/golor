package golor

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	_ "embed"

	"github.com/0mega24/golor/v2/internal/space"
)

// CSS named-color data comes from W3C CSS Color Module Level 4.
//
//go:embed css_named_colors.json
var namedColorsJSON []byte

var namedColors = loadNamedColors()

// namedIndex holds the named colors (except transparent) sorted by name, with their
// OKLab coordinates precomputed for NearestNamed.
var namedIndex = buildNamedIndex(namedColors)

type namedEntry struct {
	name  string
	color Color
	lab   space.Vec3
}

// ParseCSS parses context-free CSS Color 4 syntax, clamping components and mapping to sRGB.
// It accepts none and angle units; currentcolor and system colors return errors.
// Supported forms and spaces are listed in the README.
func ParseCSS(s string) (Color, error) {
	input := strings.TrimSpace(s)
	switch {
	case input == "":
		return Color{}, fmt.Errorf("golor: empty CSS color")
	case input[0] == '#':
		return Hex(input)
	case strings.IndexByte(input, '(') < 0:
		return Named(input)
	}
	name, body, err := cssFunction(s)
	if err != nil {
		return Color{}, err
	}
	fn, label := lookupCSSFunction(name)
	if fn == nil {
		return Color{}, fmt.Errorf("golor: unsupported CSS color function %q", s)
	}
	c, ok := fn(body)
	if !ok {
		return Color{}, fmt.Errorf("golor: invalid %s color %q", label, s)
	}
	return c, nil
}

// ParseRGB parses CSS rgb()/rgba() syntax, including legacy comma and modern space/slash
// forms. rgb and rgba are aliases. Out-of-range channels and alpha are clamped.
func ParseRGB(s string) (Color, error) {
	return parseNamedFunction(s, "rgb", "RGB", parseRGBBody)
}

// ParseHSL accepts legacy/modern hsl()/hsla(), wrapped angle units, and clamped S/L.
func ParseHSL(s string) (Color, error) {
	return parseNamedFunction(s, "hsl", "HSL", parseHSLBody)
}

func parseNamedFunction(s, want, label string, fn func(string) (Color, bool)) (Color, error) {
	name, body, err := cssFunction(s)
	if err != nil {
		return Color{}, err
	}
	if !strings.EqualFold(name, want) && !strings.EqualFold(name, want+"a") {
		return Color{}, fmt.Errorf("golor: expected %s()/%sa(), got %q", want, want, s)
	}
	c, ok := fn(body)
	if !ok {
		return Color{}, fmt.Errorf("golor: invalid %s color %q", label, s)
	}
	return c, nil
}

func lookupCSSFunction(name string) (fn func(string) (Color, bool), label string) {
	switch {
	case strings.EqualFold(name, "rgb"), strings.EqualFold(name, "rgba"):
		return parseRGBBody, "RGB"
	case strings.EqualFold(name, "hsl"), strings.EqualFold(name, "hsla"):
		return parseHSLBody, "HSL"
	case strings.EqualFold(name, "hwb"):
		return parseHWBBody, "HWB"
	case strings.EqualFold(name, "lab"):
		return parseLabBody, "Lab"
	case strings.EqualFold(name, "lch"):
		return parseLCHBody, "LCH"
	case strings.EqualFold(name, "oklab"):
		return parseOKLabBody, "OKLab"
	case strings.EqualFold(name, "oklch"):
		return parseOKLCHBody, "OKLCH"
	case strings.EqualFold(name, "color"):
		return parseColorBody, "color()"
	}
	return nil, ""
}

// Named returns the CSS Color Level 4 named color matching name (case-insensitive).
func Named(name string) (Color, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	c, ok := namedColors[key]
	if !ok {
		if key == "currentcolor" {
			return Color{}, fmt.Errorf("golor: currentcolor depends on the element's context and cannot be resolved")
		}
		return Color{}, fmt.Errorf("golor: unknown named color %q", name)
	}
	return c, nil
}

// NearestNamed returns the CSS named color closest to c by perceptual OKLab distance.
// Alpha is ignored and transparent is never returned. Ties are resolved by name.
func NearestNamed(c Color) (name string, color Color) {
	lab := oklabOf(c)
	best := -1
	bestDistance := math.Inf(1)
	for i := range namedIndex {
		e := &namedIndex[i]
		d0, d1, d2 := lab[0]-e.lab[0], lab[1]-e.lab[1], lab[2]-e.lab[2]
		if d := d0*d0 + d1*d1 + d2*d2; d < bestDistance {
			best, bestDistance = i, d
		}
	}
	return namedIndex[best].name, namedIndex[best].color
}

func oklabOf(c Color) space.Vec3 {
	return space.LinearSRGBToOKLab(space.LinearizeVec(space.Vec3{space.Clamp01(c.R), space.Clamp01(c.G), space.Clamp01(c.B)}))
}

func loadNamedColors() map[string]Color {
	var raw map[string]string
	if err := json.Unmarshal(namedColorsJSON, &raw); err != nil {
		panic(fmt.Sprintf("golor: invalid embedded CSS named colors: %v", err))
	}
	colors := make(map[string]Color, len(raw))
	for name, hex := range raw {
		c, err := Hex(hex)
		if err != nil {
			panic(fmt.Sprintf("golor: invalid embedded CSS named color %q: %v", name, err))
		}
		colors[name] = c
	}
	return colors
}

func buildNamedIndex(colors map[string]Color) []namedEntry {
	entries := make([]namedEntry, 0, len(colors))
	for name, c := range colors {
		if name == "transparent" {
			continue
		}
		entries = append(entries, namedEntry{name: name, color: c, lab: oklabOf(c)})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })
	return entries
}

// CSSRGBString returns c formatted as CSS rgb() syntax, e.g. "rgb(255 107 53)" or
// "rgb(255 107 53 / 0.502)". ParseCSS returns the same 8-bit color.
func (c Color) CSSRGBString() string {
	body := strconv.Itoa(int(c.R8())) + " " + strconv.Itoa(int(c.G8())) + " " + strconv.Itoa(int(c.B8()))
	return "rgb(" + body + cssAlphaSuffix(c) + ")"
}

// CSSHSLString formats hsl() with enough precision to round-trip c's 8-bit channels.
func (c Color) CSSHSLString() string {
	h, s, l := colorToHSL(c.quantized())
	h = math.Round(h*100) / 100
	if h >= 360 {
		h -= 360
	}
	return "hsl(" + formatCSSNumber(h, 2) + " " + formatCSSNumber(s*100, 2) + "% " +
		formatCSSNumber(l*100, 2) + "%" + cssAlphaSuffix(c) + ")"
}

// CSSOKLCHString formats oklch() to round-trip c's 8-bit channels; greys use C=H=0.
func (c Color) CSSOKLCHString() string {
	lch := space.ToPolar(oklabOf(c.quantized()), space.AchromaticOKLCH)
	chroma := math.Round(lch[1]*1e5) / 1e5
	h := math.Round(lch[2]*1000) / 1000
	if h >= 360 || chroma == 0 {
		h = 0
	}
	return "oklch(" + formatCSSNumber(lch[0]*100, 3) + "% " + formatCSSNumber(chroma, 5) + " " +
		formatCSSNumber(h, 3) + cssAlphaSuffix(c) + ")"
}

// quantized returns c's color channels rounded to the 8-bit grid.
func (c Color) quantized() Color {
	return RGB(c.R8(), c.G8(), c.B8())
}

func cssAlphaSuffix(c Color) string {
	a := c.A8()
	if a == 255 {
		return ""
	}
	// Three decimals are enough to recover any 8-bit alpha (1/255 > 0.002).
	return " / " + formatCSSNumber(float64(a)/255, 3)
}

// formatCSSNumber formats v with at most prec decimals and no trailing zeros.
func formatCSSNumber(v float64, prec int) string {
	s := strconv.FormatFloat(v, 'f', prec, 64)
	if strings.IndexByte(s, '.') >= 0 {
		s = strings.TrimRight(s, "0")
		s = strings.TrimSuffix(s, ".")
	}
	if s == "-0" {
		return "0"
	}
	return s
}

func colorToHSL(c Color) (h, s, l float64) {
	r, g, b := space.Clamp01(c.R), space.Clamp01(c.G), space.Clamp01(c.B)
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l = (max + min) / 2
	delta := max - min
	if delta == 0 {
		return 0, 0, l
	}
	if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2 - max - min)
	}
	switch max {
	case r:
		h = (g - b) / delta
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/delta + 2
	default:
		h = (r-g)/delta + 4
	}
	return h * 60, s, l
}
