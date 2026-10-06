package golor

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/0mega24/golor/v2/internal/space"
)

// cssFunction splits "name(body)" without allocating. The name must immediately precede
// the parenthesis, as in CSS.
func cssFunction(s string) (name, body string, err error) {
	input := strings.TrimSpace(s)
	open := strings.IndexByte(input, '(')
	if open < 1 || !strings.HasSuffix(input, ")") {
		return "", "", fmt.Errorf("golor: invalid CSS color %q", s)
	}
	name = input[:open]
	body = strings.TrimSpace(input[open+1 : len(input)-1])
	if body == "" || strings.ContainsAny(name, " \t\n\r\f") || strings.ContainsAny(body, "()") {
		return "", "", fmt.Errorf("golor: invalid CSS color %q", s)
	}
	return name, body, nil
}

// cssArgs holds the tokens of a color function body.
type cssArgs struct {
	comps    [4]string
	n        int
	alpha    string
	hasAlpha bool
	legacy   bool // Comma-separated syntax.
}

func isCSSSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f'
}

func trimCSSSpace(s string) string {
	for len(s) > 0 && isCSSSpace(s[0]) {
		s = s[1:]
	}
	for len(s) > 0 && isCSSSpace(s[len(s)-1]) {
		s = s[:len(s)-1]
	}
	return s
}

func hasCSSSpace(s string) bool {
	for i := 0; i < len(s); i++ {
		if isCSSSpace(s[i]) {
			return true
		}
	}
	return false
}

// splitCSSArgs tokenizes a function body in either legacy comma syntax
// ("a, b, c[, alpha]") or modern space syntax ("a b c [/ alpha]").
func splitCSSArgs(body string) (cssArgs, bool) {
	var a cssArgs
	if strings.IndexByte(body, ',') >= 0 {
		if strings.IndexByte(body, '/') >= 0 {
			return a, false
		}
		a.legacy = true
		var items [4]string
		n := 0
		for {
			comma := strings.IndexByte(body, ',')
			item := body
			if comma >= 0 {
				item = body[:comma]
			}
			item = trimCSSSpace(item)
			if item == "" || hasCSSSpace(item) || n == len(items) {
				return a, false
			}
			items[n] = item
			n++
			if comma < 0 {
				break
			}
			body = body[comma+1:]
		}
		switch n {
		case 3:
		case 4:
			a.alpha, a.hasAlpha = items[3], true
		default:
			return a, false
		}
		copy(a.comps[:], items[:3])
		a.n = 3
		return a, true
	}
	if slash := strings.IndexByte(body, '/'); slash >= 0 {
		alpha := trimCSSSpace(body[slash+1:])
		if alpha == "" || hasCSSSpace(alpha) || strings.IndexByte(alpha, '/') >= 0 {
			return a, false
		}
		a.alpha, a.hasAlpha = alpha, true
		body = body[:slash]
	}
	i := 0
	for i < len(body) {
		for i < len(body) && isCSSSpace(body[i]) {
			i++
		}
		if i == len(body) {
			break
		}
		start := i
		for i < len(body) && !isCSSSpace(body[i]) {
			i++
		}
		if a.n == len(a.comps) {
			return a, false
		}
		a.comps[a.n] = body[start:i]
		a.n++
	}
	return a, true
}

type compKind uint8

const (
	compNumber compKind = iota
	compPercent
	compAngle // Value in degrees.
	compNone
)

type cssComp struct {
	kind compKind
	v    float64
}

// isCSSNumber reports whether s matches the CSS <number> grammar
// ([+-]? digits? (. digits)? ([eE][+-]? digits)?), which is stricter than ParseFloat.
func isCSSNumber(s string) bool {
	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	digits := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
		digits++
	}
	if i < len(s) && s[i] == '.' {
		i++
		frac := 0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
			frac++
		}
		if frac == 0 {
			return false
		}
		digits += frac
	}
	if digits == 0 {
		return false
	}
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			i++
		}
		exp := 0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
			exp++
		}
		if exp == 0 {
			return false
		}
	}
	return i == len(s)
}

func parseCSSNumber(s string) (float64, bool) {
	// Fast path for short unsigned integers, the most common CSS components.
	if n := len(s); n > 0 && n <= 6 {
		v := 0
		for i := 0; i < n; i++ {
			d := s[i] - '0'
			if d > 9 {
				v = -1
				break
			}
			v = v*10 + int(d)
		}
		if v >= 0 {
			return float64(v), true
		}
	}
	if !isCSSNumber(s) {
		return 0, false
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsInf(v, 0) || math.IsNaN(v) {
		return 0, false
	}
	return v, true
}

func hasSuffixFold(s, suffix string) bool {
	return len(s) > len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix)
}

func parseCSSComp(tok string) (cssComp, bool) {
	if strings.EqualFold(tok, "none") {
		return cssComp{kind: compNone}, true
	}
	if strings.HasSuffix(tok, "%") {
		v, ok := parseCSSNumber(tok[:len(tok)-1])
		return cssComp{compPercent, v}, ok
	}
	units := [...]struct {
		suffix string
		scale  float64
	}{{"grad", 0.9}, {"turn", 360}, {"deg", 1}, {"rad", 180 / math.Pi}}
	for _, u := range units {
		if hasSuffixFold(tok, u.suffix) {
			v, ok := parseCSSNumber(tok[:len(tok)-len(u.suffix)])
			return cssComp{compAngle, v * u.scale}, ok && !math.IsInf(v*u.scale, 0)
		}
	}
	v, ok := parseCSSNumber(tok)
	return cssComp{compNumber, v}, ok
}

// parseHue accepts a number (degrees), an angle, or none (modern syntax only).
func parseHue(tok string, legacy bool) (float64, bool) {
	c, ok := parseCSSComp(tok)
	if !ok {
		return 0, false
	}
	switch c.kind {
	case compNumber, compAngle:
		return space.NormalizeHue(c.v), true
	case compNone:
		return 0, !legacy
	case compPercent:
		return 0, false
	}
	return 0, false
}

// parseScaled accepts a number, a percentage (where 100% equals pct100), or none
// (modern syntax only).
func parseScaled(tok string, pct100 float64, legacy bool) (cssComp, float64, bool) {
	c, ok := parseCSSComp(tok)
	if !ok {
		return c, 0, false
	}
	switch c.kind {
	case compNumber:
		return c, c.v, true
	case compPercent:
		return c, c.v / 100 * pct100, true
	case compNone:
		return c, 0, !legacy
	case compAngle:
		return c, 0, false
	}
	return c, 0, false
}

func parseCSSAlpha(a cssArgs) (float64, bool) {
	if !a.hasAlpha {
		return 1, true
	}
	_, v, ok := parseScaled(a.alpha, 1, a.legacy)
	return space.Clamp01(v), ok
}

func parseRGBBody(body string) (Color, bool) {
	a, ok := splitCSSArgs(body)
	if !ok || a.n != 3 {
		return Color{}, false
	}
	var ch [3]float64
	var kinds [3]compKind
	for i := range ch {
		c, v, ok := parseScaled(a.comps[i], 255, a.legacy)
		if !ok {
			return Color{}, false
		}
		ch[i], kinds[i] = v/255, c.kind
	}
	if a.legacy && (kinds[0] != kinds[1] || kinds[1] != kinds[2]) {
		return Color{}, false // Legacy rgb() cannot mix numbers and percentages.
	}
	alpha, ok := parseCSSAlpha(a)
	if !ok {
		return Color{}, false
	}
	return RGBAf(ch[0], ch[1], ch[2], alpha), true
}

func parseHSLBody(body string) (Color, bool) {
	a, ok := splitCSSArgs(body)
	if !ok || a.n != 3 {
		return Color{}, false
	}
	h, ok := parseHue(a.comps[0], a.legacy)
	if !ok {
		return Color{}, false
	}
	var sl [2]float64
	for i := range sl {
		c, v, ok := parseScaled(a.comps[i+1], 100, a.legacy)
		if !ok || (a.legacy && c.kind != compPercent) {
			return Color{}, false
		}
		sl[i] = space.Clamp01(v / 100)
	}
	alpha, ok := parseCSSAlpha(a)
	if !ok {
		return Color{}, false
	}
	return hslToColor(h, sl[0], sl[1], alpha), true
}

func parseHWBBody(body string) (Color, bool) {
	a, ok := splitCSSArgs(body)
	if !ok || a.legacy || a.n != 3 {
		return Color{}, false
	}
	h, ok := parseHue(a.comps[0], false)
	if !ok {
		return Color{}, false
	}
	var wb [2]float64
	for i := range wb {
		_, v, ok := parseScaled(a.comps[i+1], 100, false)
		if !ok {
			return Color{}, false
		}
		wb[i] = space.Clamp01(v / 100)
	}
	alpha, ok := parseCSSAlpha(a)
	if !ok {
		return Color{}, false
	}
	w, b := wb[0], wb[1]
	if w+b >= 1 {
		g := w / (w + b)
		return RGBAf(g, g, g, alpha), true
	}
	base := hslToColor(h, 1, 0.5, 1)
	scale := 1 - w - b
	return RGBAf(base.R*scale+w, base.G*scale+w, base.B*scale+w, alpha), true
}

// parseThree parses three modern-syntax components with the given percentage
// references.
func parseThree(a cssArgs, offset int, pct100 [3]float64) (v [3]float64, ok bool) {
	for i := range v {
		_, v[i], ok = parseScaled(a.comps[offset+i], pct100[i], false)
		if !ok {
			return v, false
		}
	}
	return v, true
}

func parseLabBody(body string) (Color, bool) {
	a, ok := splitCSSArgs(body)
	if !ok || a.legacy || a.n != 3 {
		return Color{}, false
	}
	v, ok := parseThree(a, 0, [3]float64{100, 125, 125})
	if !ok {
		return Color{}, false
	}
	v[0] = math.Max(0, math.Min(100, v[0]))
	return finishCSSLab(v, a)
}

func parseLCHBody(body string) (Color, bool) {
	a, ok := splitCSSArgs(body)
	if !ok || a.legacy || a.n != 3 {
		return Color{}, false
	}
	var v [3]float64
	if _, v[0], ok = parseScaled(a.comps[0], 100, false); !ok {
		return Color{}, false
	}
	if _, v[1], ok = parseScaled(a.comps[1], 150, false); !ok {
		return Color{}, false
	}
	if v[2], ok = parseHue(a.comps[2], false); !ok {
		return Color{}, false
	}
	v[0] = math.Max(0, math.Min(100, v[0]))
	v[1] = math.Max(0, v[1])
	return finishCSSLab(space.FromPolar(v), a)
}

// finishCSSLab converts CSS lab() coordinates, which are relative to D50, to sRGB.
func finishCSSLab(lab space.Vec3, a cssArgs) (Color, bool) {
	xyz := space.D50ToD65.Mul(space.LabToXYZ(lab, space.D50))
	return finishLinear(space.XYZToSRGB.Mul(xyz), a)
}

func parseOKLabBody(body string) (Color, bool) {
	a, ok := splitCSSArgs(body)
	if !ok || a.legacy || a.n != 3 {
		return Color{}, false
	}
	v, ok := parseThree(a, 0, [3]float64{1, 0.4, 0.4})
	if !ok {
		return Color{}, false
	}
	v[0] = space.Clamp01(v[0])
	return finishLinear(space.OKLabToLinearSRGB(v), a)
}

func parseOKLCHBody(body string) (Color, bool) {
	a, ok := splitCSSArgs(body)
	if !ok || a.legacy || a.n != 3 {
		return Color{}, false
	}
	var v [3]float64
	if _, v[0], ok = parseScaled(a.comps[0], 1, false); !ok {
		return Color{}, false
	}
	if _, v[1], ok = parseScaled(a.comps[1], 0.4, false); !ok {
		return Color{}, false
	}
	if v[2], ok = parseHue(a.comps[2], false); !ok {
		return Color{}, false
	}
	alpha, ok := parseCSSAlpha(a)
	if !ok {
		return Color{}, false
	}
	m := space.GamutMapOKLCH(space.Vec3{space.Clamp01(v[0]), math.Max(0, v[1]), v[2]})
	return RGBAf(space.Delinearize(m[0]), space.Delinearize(m[1]), space.Delinearize(m[2]), alpha), true
}

func parseColorBody(body string) (Color, bool) {
	a, ok := splitCSSArgs(body)
	if !ok || a.legacy || a.n != 4 {
		return Color{}, false
	}
	v, ok := parseThree(a, 1, [3]float64{1, 1, 1})
	if !ok {
		return Color{}, false
	}
	var lin space.Vec3
	switch id := a.comps[0]; {
	case strings.EqualFold(id, "srgb"):
		if v[0] >= 0 && v[0] <= 1 && v[1] >= 0 && v[1] <= 1 && v[2] >= 0 && v[2] <= 1 {
			alpha, ok := parseCSSAlpha(a)
			return RGBAf(v[0], v[1], v[2], alpha), ok
		}
		lin = space.LinearizeVec(v)
	case strings.EqualFold(id, "srgb-linear"):
		lin = v
	case strings.EqualFold(id, "display-p3"):
		lin = space.XYZToSRGB.Mul(space.P3ToXYZ.Mul(space.LinearizeVec(v)))
	case strings.EqualFold(id, "xyz"), strings.EqualFold(id, "xyz-d65"):
		lin = space.XYZToSRGB.Mul(v)
	case strings.EqualFold(id, "xyz-d50"):
		lin = space.XYZToSRGB.Mul(space.D50ToD65.Mul(v))
	default:
		return Color{}, false
	}
	return finishLinear(lin, a)
}

// finishLinear gamut-maps linear sRGB and attaches the parsed alpha.
func finishLinear(lin space.Vec3, a cssArgs) (Color, bool) {
	alpha, ok := parseCSSAlpha(a)
	if !ok {
		return Color{}, false
	}
	m := space.GamutMapLinear(lin)
	return RGBAf(space.Delinearize(m[0]), space.Delinearize(m[1]), space.Delinearize(m[2]), alpha), true
}

// hslToColor expects h in [0, 360) and s, l in [0, 1].
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

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t >= 1 {
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
