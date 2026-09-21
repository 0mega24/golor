package golor

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	_ "embed"
)

// CSS named-color data comes from W3C CSS Color Module Level 4.
//
//go:embed css_named_colors.json
var namedColorsJSON []byte

var namedColors = loadNamedColors()

// ParseCSS parses a CSS named color, hex color, rgb()/rgba(), or hsl()/hsla() color.
func ParseCSS(s string) (Color, error) {
	input := strings.TrimSpace(s)
	lower := strings.ToLower(input)
	switch {
	case strings.HasPrefix(input, "#"):
		return Hex(input)
	case strings.HasPrefix(lower, "rgb("), strings.HasPrefix(lower, "rgba("):
		return ParseRGB(input)
	case strings.HasPrefix(lower, "hsl("), strings.HasPrefix(lower, "hsla("):
		return ParseHSL(input)
	default:
		return Named(input)
	}
}

// ParseRGB parses CSS rgb()/rgba() syntax, including comma and space/slash forms.
func ParseRGB(s string) (Color, error) {
	name, body, err := cssFunction(s)
	if err != nil {
		return Color{}, err
	}
	if name != "rgb" && name != "rgba" {
		return Color{}, fmt.Errorf("golor: expected rgb()/rgba(), got %q", s)
	}
	parts, alpha, hasAlpha := splitCSSComponents(body)
	if name == "rgba" && !hasAlpha {
		if len(parts) != 4 {
			return Color{}, fmt.Errorf("golor: invalid RGB color %q", s)
		}
		alpha = parts[3]
		parts = parts[:3]
		hasAlpha = true
	}
	if len(parts) != 3 {
		return Color{}, fmt.Errorf("golor: invalid RGB color %q", s)
	}
	r, err := parseRGBChannel(parts[0])
	if err != nil {
		return Color{}, fmt.Errorf("golor: invalid RGB color %q", s)
	}
	g, err := parseRGBChannel(parts[1])
	if err != nil {
		return Color{}, fmt.Errorf("golor: invalid RGB color %q", s)
	}
	b, err := parseRGBChannel(parts[2])
	if err != nil {
		return Color{}, fmt.Errorf("golor: invalid RGB color %q", s)
	}
	a := 1.0
	if hasAlpha {
		a, err = parseAlpha(alpha)
		if err != nil {
			return Color{}, fmt.Errorf("golor: invalid RGB color %q", s)
		}
	}
	return RGBAf(r, g, b, a), nil
}

// ParseHSL parses CSS hsl()/hsla() syntax, including comma and space/slash forms.
func ParseHSL(s string) (Color, error) {
	name, body, err := cssFunction(s)
	if err != nil {
		return Color{}, err
	}
	if name != "hsl" && name != "hsla" {
		return Color{}, fmt.Errorf("golor: expected hsl()/hsla(), got %q", s)
	}
	parts, alpha, hasAlpha := splitCSSComponents(body)
	if name == "hsla" && !hasAlpha {
		if len(parts) != 4 {
			return Color{}, fmt.Errorf("golor: invalid HSL color %q", s)
		}
		alpha = parts[3]
		parts = parts[:3]
		hasAlpha = true
	}
	if len(parts) != 3 {
		return Color{}, fmt.Errorf("golor: invalid HSL color %q", s)
	}
	h, err := parseHue(parts[0])
	if err != nil {
		return Color{}, fmt.Errorf("golor: invalid HSL color %q", s)
	}
	sat, err := parsePercent(parts[1])
	if err != nil {
		return Color{}, fmt.Errorf("golor: invalid HSL color %q", s)
	}
	light, err := parsePercent(parts[2])
	if err != nil {
		return Color{}, fmt.Errorf("golor: invalid HSL color %q", s)
	}
	a := 1.0
	if hasAlpha {
		a, err = parseAlpha(alpha)
		if err != nil {
			return Color{}, fmt.Errorf("golor: invalid HSL color %q", s)
		}
	}
	return hslToColor(h, sat, light, a), nil
}

// Named returns the CSS Color Level 4 named color matching name.
func Named(name string) (Color, error) {
	c, ok := namedColors[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return Color{}, fmt.Errorf("golor: unknown named color %q", name)
	}
	return c, nil
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

// CSSRGBString returns c formatted as CSS rgb() or rgba() syntax.
func (c Color) CSSRGBString() string {
	if c.A8() == 255 {
		return fmt.Sprintf("rgb(%d %d %d)", c.R8(), c.G8(), c.B8())
	}
	return fmt.Sprintf("rgb(%d %d %d / %.3g)", c.R8(), c.G8(), c.B8(), clamp01(c.A))
}

// CSSHSLString returns c formatted as CSS hsl() syntax.
func (c Color) CSSHSLString() string {
	h, s, l := colorToHSL(c)
	if c.A8() == 255 {
		return fmt.Sprintf("hsl(%.3g %.3g%% %.3g%%)", h, s*100, l*100)
	}
	return fmt.Sprintf("hsl(%.3g %.3g%% %.3g%% / %.3g)", h, s*100, l*100, clamp01(c.A))
}

func cssFunction(s string) (name, body string, err error) {
	input := strings.TrimSpace(s)
	open := strings.IndexByte(input, '(')
	if open < 1 || !strings.HasSuffix(input, ")") {
		return "", "", fmt.Errorf("golor: invalid CSS color %q", s)
	}
	name = strings.ToLower(strings.TrimSpace(input[:open]))
	body = strings.TrimSpace(input[open+1 : len(input)-1])
	if body == "" {
		return "", "", fmt.Errorf("golor: invalid CSS color %q", s)
	}
	return name, body, nil
}

func splitCSSComponents(body string) ([]string, string, bool) {
	body = strings.TrimSpace(body)
	var alpha string
	hasAlpha := false
	if slash := strings.IndexByte(body, '/'); slash >= 0 {
		alpha = strings.TrimSpace(body[slash+1:])
		body = strings.TrimSpace(body[:slash])
		hasAlpha = true
	}
	body = strings.ReplaceAll(body, ",", " ")
	fields := strings.FieldsFunc(body, unicode.IsSpace)
	return fields, alpha, hasAlpha
}

func parseRGBChannel(s string) (float64, error) {
	if strings.HasSuffix(s, "%") {
		return parsePercent(s)
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v < 0 || v > 255 {
		return 0, fmt.Errorf("invalid RGB channel")
	}
	return v / 255, nil
}

func parseAlpha(s string) (float64, error) {
	if strings.HasSuffix(s, "%") {
		return parsePercent(s)
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v < 0 || v > 1 {
		return 0, fmt.Errorf("invalid alpha")
	}
	return v, nil
}

func parsePercent(s string) (float64, error) {
	if !strings.HasSuffix(s, "%") {
		return 0, fmt.Errorf("expected percentage")
	}
	v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
	if err != nil || v < 0 || v > 100 {
		return 0, fmt.Errorf("invalid percentage")
	}
	return v / 100, nil
}

func parseHue(s string) (float64, error) {
	s = strings.TrimSuffix(s, "deg")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v < 0 || v >= 360 {
		return 0, fmt.Errorf("invalid hue")
	}
	return v, nil
}

func hslToColor(h, s, l, a float64) Color {
	if s == 0 {
		return RGBAf(l, l, l, a)
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	return RGBAf(hueToRGB(p, q, h/360+1.0/3), hueToRGB(p, q, h/360), hueToRGB(p, q, h/360-1.0/3), a)
}

func colorToHSL(c Color) (h, s, l float64) {
	r, g, b := c.R, c.G, c.B
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

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 1.0/2:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	default:
		return p
	}
}
