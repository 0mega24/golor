package golor

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

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

var namedColors = map[string]Color{
	"aliceblue":            RGB(240, 248, 255),
	"antiquewhite":         RGB(250, 235, 215),
	"aqua":                 RGB(0, 255, 255),
	"aquamarine":           RGB(127, 255, 212),
	"azure":                RGB(240, 255, 255),
	"beige":                RGB(245, 245, 220),
	"bisque":               RGB(255, 228, 196),
	"black":                RGB(0, 0, 0),
	"blanchedalmond":       RGB(255, 235, 205),
	"blue":                 RGB(0, 0, 255),
	"blueviolet":           RGB(138, 43, 226),
	"brown":                RGB(165, 42, 42),
	"burlywood":            RGB(222, 184, 135),
	"cadetblue":            RGB(95, 158, 160),
	"chartreuse":           RGB(127, 255, 0),
	"chocolate":            RGB(210, 105, 30),
	"coral":                RGB(255, 127, 80),
	"cornflowerblue":       RGB(100, 149, 237),
	"cornsilk":             RGB(255, 248, 220),
	"crimson":              RGB(220, 20, 60),
	"cyan":                 RGB(0, 255, 255),
	"darkblue":             RGB(0, 0, 139),
	"darkcyan":             RGB(0, 139, 139),
	"darkgoldenrod":        RGB(184, 134, 11),
	"darkgray":             RGB(169, 169, 169),
	"darkgreen":            RGB(0, 100, 0),
	"darkgrey":             RGB(169, 169, 169),
	"darkkhaki":            RGB(189, 183, 107),
	"darkmagenta":          RGB(139, 0, 139),
	"darkolivegreen":       RGB(85, 107, 47),
	"darkorange":           RGB(255, 140, 0),
	"darkorchid":           RGB(153, 50, 204),
	"darkred":              RGB(139, 0, 0),
	"darksalmon":           RGB(233, 150, 122),
	"darkseagreen":         RGB(143, 188, 143),
	"darkslateblue":        RGB(72, 61, 139),
	"darkslategray":        RGB(47, 79, 79),
	"darkslategrey":        RGB(47, 79, 79),
	"darkturquoise":        RGB(0, 206, 209),
	"darkviolet":           RGB(148, 0, 211),
	"deeppink":             RGB(255, 20, 147),
	"deepskyblue":          RGB(0, 191, 255),
	"dimgray":              RGB(105, 105, 105),
	"dimgrey":              RGB(105, 105, 105),
	"dodgerblue":           RGB(30, 144, 255),
	"firebrick":            RGB(178, 34, 34),
	"floralwhite":          RGB(255, 250, 240),
	"forestgreen":          RGB(34, 139, 34),
	"fuchsia":              RGB(255, 0, 255),
	"gainsboro":            RGB(220, 220, 220),
	"ghostwhite":           RGB(248, 248, 255),
	"gold":                 RGB(255, 215, 0),
	"goldenrod":            RGB(218, 165, 32),
	"gray":                 RGB(128, 128, 128),
	"green":                RGB(0, 128, 0),
	"greenyellow":          RGB(173, 255, 47),
	"grey":                 RGB(128, 128, 128),
	"honeydew":             RGB(240, 255, 240),
	"hotpink":              RGB(255, 105, 180),
	"indianred":            RGB(205, 92, 92),
	"indigo":               RGB(75, 0, 130),
	"ivory":                RGB(255, 255, 240),
	"khaki":                RGB(240, 230, 140),
	"lavender":             RGB(230, 230, 250),
	"lavenderblush":        RGB(255, 240, 245),
	"lawngreen":            RGB(124, 252, 0),
	"lemonchiffon":         RGB(255, 250, 205),
	"lightblue":            RGB(173, 216, 230),
	"lightcoral":           RGB(240, 128, 128),
	"lightcyan":            RGB(224, 255, 255),
	"lightgoldenrodyellow": RGB(250, 250, 210),
	"lightgray":            RGB(211, 211, 211),
	"lightgreen":           RGB(144, 238, 144),
	"lightgrey":            RGB(211, 211, 211),
	"lightpink":            RGB(255, 182, 193),
	"lightsalmon":          RGB(255, 160, 122),
	"lightseagreen":        RGB(32, 178, 170),
	"lightskyblue":         RGB(135, 206, 250),
	"lightslategray":       RGB(119, 136, 153),
	"lightslategrey":       RGB(119, 136, 153),
	"lightsteelblue":       RGB(176, 196, 222),
	"lightyellow":          RGB(255, 255, 224),
	"lime":                 RGB(0, 255, 0),
	"limegreen":            RGB(50, 205, 50),
	"linen":                RGB(250, 240, 230),
	"magenta":              RGB(255, 0, 255),
	"maroon":               RGB(128, 0, 0),
	"mediumaquamarine":     RGB(102, 205, 170),
	"mediumblue":           RGB(0, 0, 205),
	"mediumorchid":         RGB(186, 85, 211),
	"mediumpurple":         RGB(147, 112, 219),
	"mediumseagreen":       RGB(60, 179, 113),
	"mediumslateblue":      RGB(123, 104, 238),
	"mediumspringgreen":    RGB(0, 250, 154),
	"mediumturquoise":      RGB(72, 209, 204),
	"mediumvioletred":      RGB(199, 21, 133),
	"midnightblue":         RGB(25, 25, 112),
	"mintcream":            RGB(245, 255, 250),
	"mistyrose":            RGB(255, 228, 225),
	"moccasin":             RGB(255, 228, 181),
	"navajowhite":          RGB(255, 222, 173),
	"navy":                 RGB(0, 0, 128),
	"oldlace":              RGB(253, 245, 230),
	"olive":                RGB(128, 128, 0),
	"olivedrab":            RGB(107, 142, 35),
	"orange":               RGB(255, 165, 0),
	"orangered":            RGB(255, 69, 0),
	"orchid":               RGB(218, 112, 214),
	"palegoldenrod":        RGB(238, 232, 170),
	"palegreen":            RGB(152, 251, 152),
	"paleturquoise":        RGB(175, 238, 238),
	"palevioletred":        RGB(219, 112, 147),
	"papayawhip":           RGB(255, 239, 213),
	"peachpuff":            RGB(255, 218, 185),
	"peru":                 RGB(205, 133, 63),
	"pink":                 RGB(255, 192, 203),
	"plum":                 RGB(221, 160, 221),
	"powderblue":           RGB(176, 224, 230),
	"purple":               RGB(128, 0, 128),
	"rebeccapurple":        RGB(102, 51, 153),
	"red":                  RGB(255, 0, 0),
	"rosybrown":            RGB(188, 143, 143),
	"royalblue":            RGB(65, 105, 225),
	"saddlebrown":          RGB(139, 69, 19),
	"salmon":               RGB(250, 128, 114),
	"sandybrown":           RGB(244, 164, 96),
	"seagreen":             RGB(46, 139, 87),
	"seashell":             RGB(255, 245, 238),
	"sienna":               RGB(160, 82, 45),
	"silver":               RGB(192, 192, 192),
	"skyblue":              RGB(135, 206, 235),
	"slateblue":            RGB(106, 90, 205),
	"slategray":            RGB(112, 128, 144),
	"slategrey":            RGB(112, 128, 144),
	"snow":                 RGB(255, 250, 250),
	"springgreen":          RGB(0, 255, 127),
	"steelblue":            RGB(70, 130, 180),
	"tan":                  RGB(210, 180, 140),
	"teal":                 RGB(0, 128, 128),
	"thistle":              RGB(216, 191, 216),
	"tomato":               RGB(255, 99, 71),
	"transparent":          RGBA(0, 0, 0, 0),
	"turquoise":            RGB(64, 224, 208),
	"violet":               RGB(238, 130, 238),
	"wheat":                RGB(245, 222, 179),
	"white":                RGB(255, 255, 255),
	"whitesmoke":           RGB(245, 245, 245),
	"yellow":               RGB(255, 255, 0),
	"yellowgreen":          RGB(154, 205, 50),
}
