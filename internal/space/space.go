// Package space shares color math as float64 triples to avoid import cycles.
package space

import "math"

// Vec3 is a color triple in the space documented by the function.
type Vec3 [3]float64

// Mat3 is a row-major 3x3 matrix.
type Mat3 [3][3]float64

// Mul returns m·v.
func (m *Mat3) Mul(v Vec3) Vec3 {
	return Vec3{
		m[0][0]*v[0] + m[0][1]*v[1] + m[0][2]*v[2],
		m[1][0]*v[0] + m[1][1]*v[1] + m[1][2]*v[2],
		m[2][0]*v[0] + m[2][1]*v[1] + m[2][2]*v[2],
	}
}

func mulMat(a, b Mat3) Mat3 {
	var out Mat3
	for i := range 3 {
		for j := range 3 {
			out[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j] + a[i][2]*b[2][j]
		}
	}
	return out
}

func invert(m Mat3) Mat3 {
	a, b, c := m[0][0], m[0][1], m[0][2]
	d, e, f := m[1][0], m[1][1], m[1][2]
	g, h, i := m[2][0], m[2][1], m[2][2]
	det := a*(e*i-f*h) - b*(d*i-f*g) + c*(d*h-e*g)
	return Mat3{
		{(e*i - f*h) / det, (c*h - b*i) / det, (b*f - c*e) / det},
		{(f*g - d*i) / det, (a*i - c*g) / det, (c*d - a*f) / det},
		{(d*h - e*g) / det, (b*g - a*h) / det, (a*e - b*d) / det},
	}
}

// D65 and D50 are CSS Color 4 white points (chromaticities 0.3127/0.3290 and 0.3457/0.3585).
var (
	D65 = Vec3{0.3127 / 0.3290, 1, (1 - 0.3127 - 0.3290) / 0.3290}
	D50 = Vec3{0.3457 / 0.3585, 1, (1 - 0.3457 - 0.3585) / 0.3585}
)

// SRGBToXYZ and XYZToSRGB convert linear sRGB <-> XYZ (D65) using the exact rational CSS Color 4 matrix so that
// white maps exactly onto D65 and greys have zero a*/b*.
var (
	SRGBToXYZ = Mat3{
		{506752.0 / 1228815, 87881.0 / 245763, 12673.0 / 70218},
		{87098.0 / 409605, 175762.0 / 245763, 12673.0 / 175545},
		{7918.0 / 409605, 87881.0 / 737289, 1001167.0 / 1053270},
	}
	XYZToSRGB = invert(SRGBToXYZ)
)

// P3ToXYZ and XYZToP3 convert linear Display P3 <-> XYZ (D65) using CSS Color 4 matrices.
var (
	P3ToXYZ = Mat3{
		{608311.0 / 1250200, 189793.0 / 714400, 198249.0 / 1000160},
		{35783.0 / 156275, 247089.0 / 357200, 198249.0 / 2500400},
		{0, 32229.0 / 714400, 5220557.0 / 5000800},
	}
	XYZToP3 = invert(P3ToXYZ)
)

// D50ToD65 and D65ToD50 apply Bradford chromatic adaptation between D50 and D65.
var (
	D50ToD65 = bradford(D50, D65)
	D65ToD50 = invert(D50ToD65)
)

func bradford(src, dst Vec3) Mat3 {
	m := Mat3{
		{0.8951, 0.2664, -0.1614},
		{-0.7502, 1.7135, 0.0367},
		{0.0389, -0.0685, 1.0296},
	}
	s := m.Mul(src)
	d := m.Mul(dst)
	scale := Mat3{{d[0] / s[0], 0, 0}, {0, d[1] / s[1], 0}, {0, 0, d[2] / s[2]}}
	return mulMat(invert(m), mulMat(scale, m))
}

// oklabM1, oklabM2, and their inverses are Björn Ottosson's linear sRGB-based OKLab matrices.
var (
	oklabM1 = Mat3{
		{0.4122214708, 0.5363325363, 0.0514459929},
		{0.2119034982, 0.6806995451, 0.1073969566},
		{0.0883024619, 0.2817188376, 0.6299787005},
	}
	oklabM2 = Mat3{
		{0.2104542553, 0.7936177850, -0.0040720468},
		{1.9779984951, -2.4285922050, 0.4505937099},
		{0.0259040371, 0.7827717662, -0.8086757660},
	}
	oklabM1Inv = invert(oklabM1)
	oklabM2Inv = invert(oklabM2)
)

var linearLUT = func() [256]float64 {
	var t [256]float64
	for i := range t {
		t[i] = linearizeSlow(float64(i) / 255)
	}
	return t
}()

// Linearize decodes sRGB with CSS Color 4 odd symmetry; 8-bit inputs use a lookup table.
func Linearize(v float64) float64 {
	if v >= 0 && v <= 1 {
		f := v * 255
		if r := math.Round(f); math.Abs(f-r) < 1e-9 {
			return linearLUT[int(r)]
		}
	}
	return linearizeSlow(v)
}

func linearizeSlow(v float64) float64 {
	a := math.Abs(v)
	if a <= 0.04045 {
		return v / 12.92
	}
	return math.Copysign(math.Exp(2.4*math.Log((a+0.055)/1.055)), v)
}

// Delinearize encodes sRGB with odd symmetry. The exact 0.04045/12.92 breakpoint
// keeps it consistent with Linearize for lossless round trips.
func Delinearize(v float64) float64 {
	if v > 0.04045/12.92 && v <= 1 {
		return 1.055*math.Exp(math.Log(v)/2.4) - 0.055
	}
	a := math.Abs(v)
	if a <= 0.04045/12.92 {
		return 12.92 * v
	}
	return math.Copysign(1.055*math.Pow(a, 1/2.4)-0.055, v)
}

// LinearizeVec linearizes each channel.
func LinearizeVec(v Vec3) Vec3 {
	return Vec3{Linearize(v[0]), Linearize(v[1]), Linearize(v[2])}
}

// DelinearizeVec delinearizes each channel.
func DelinearizeVec(v Vec3) Vec3 {
	return Vec3{Delinearize(v[0]), Delinearize(v[1]), Delinearize(v[2])}
}

// labEpsilon and labKappa are CIE constants in exact rational form.
const (
	labEpsilon = 216.0 / 24389
	labKappa   = 24389.0 / 27
)

// XYZToLab converts XYZ to CIE L*a*b* relative to white.
func XYZToLab(xyz, white Vec3) Vec3 {
	f := func(t float64) float64 {
		if t > labEpsilon {
			return math.Cbrt(t)
		}
		return (labKappa*t + 16) / 116
	}
	fx := f(xyz[0] / white[0])
	fy := f(xyz[1] / white[1])
	fz := f(xyz[2] / white[2])
	return Vec3{116*fy - 16, 500 * (fx - fy), 200 * (fy - fz)}
}

// LabToXYZ converts CIE L*a*b* relative to white back to XYZ.
func LabToXYZ(lab, white Vec3) Vec3 {
	fy := (lab[0] + 16) / 116
	fx := lab[1]/500 + fy
	fz := fy - lab[2]/200
	var y float64
	if lab[0] > labKappa*labEpsilon {
		y = fy * fy * fy
	} else {
		y = lab[0] / labKappa
	}
	return Vec3{labFInv(fx) * white[0], y * white[1], labFInv(fz) * white[2]}
}

func labFInv(t float64) float64 {
	if t3 := t * t * t; t3 > labEpsilon {
		return t3
	}
	return (116*t - 16) / labKappa
}

// LinearSRGBToOKLab converts linear sRGB to OKLab.
func LinearSRGBToOKLab(rgb Vec3) Vec3 {
	lms := oklabM1.Mul(rgb)
	return oklabM2.Mul(Vec3{math.Cbrt(lms[0]), math.Cbrt(lms[1]), math.Cbrt(lms[2])})
}

// OKLabToLinearSRGB converts OKLab to (possibly out-of-range) linear sRGB.
func OKLabToLinearSRGB(lab Vec3) Vec3 {
	lms := oklabM2Inv.Mul(lab)
	return oklabM1Inv.Mul(Vec3{lms[0] * lms[0] * lms[0], lms[1] * lms[1] * lms[1], lms[2] * lms[2] * lms[2]})
}

// ToPolar converts rectangular (L, a, b) to (L, C, H°). Hue is 0 when chroma is below
// achromatic, because hue is undefined for greys and reflects only float noise.
func ToPolar(v Vec3, achromatic float64) Vec3 {
	c := math.Hypot(v[1], v[2])
	if c < achromatic {
		return Vec3{v[0], c, 0}
	}
	return Vec3{v[0], c, NormalizeHue(math.Atan2(v[2], v[1]) * 180 / math.Pi)}
}

// FromPolar converts (L, C, H°) to rectangular (L, a, b).
func FromPolar(v Vec3) Vec3 {
	s, c := math.Sincos(v[2] * math.Pi / 180)
	return Vec3{v[0], v[1] * c, v[1] * s}
}

// AchromaticLCH and AchromaticOKLCH are chroma thresholds below which hue has no effect.
const (
	AchromaticLCH   = 1e-4
	AchromaticOKLCH = 1e-6
)

// NormalizeHue wraps h into [0, 360). NaN and infinities become 0.
func NormalizeHue(h float64) float64 {
	if h >= 0 && h < 360 {
		return h
	}
	if math.IsNaN(h) || math.IsInf(h, 0) {
		return 0
	}
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	if h >= 360 { // Adding 360 to a tiny negative value rounds to 360.
		h = 0
	}
	return h
}

// Clamp01 clamps v into [0, 1]. NaN becomes 0.
func Clamp01(v float64) float64 {
	if v > 0 {
		if v > 1 {
			return 1
		}
		return v
	}
	return 0 // v is <= 0 or NaN.
}

// InGamutLinear reports whether linear RGB lies inside [0, 1] (with a tiny tolerance
// for float error).
func InGamutLinear(rgb Vec3) bool {
	const tol = 1e-9
	for _, c := range rgb {
		if !(c >= -tol && c <= 1+tol) { // Also rejects NaN.
			return false
		}
	}
	return true
}

func clipLinear(rgb Vec3) Vec3 {
	return Vec3{Clamp01(rgb[0]), Clamp01(rgb[1]), Clamp01(rgb[2])}
}

// Vec3Distance returns the Euclidean distance between a and b.
func Vec3Distance(a, b Vec3) float64 {
	return math.Sqrt((a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1]) + (a[2]-b[2])*(a[2]-b[2]))
}

// GamutMapLinear maps to sRGB by reducing OKLCH chroma at constant lightness/hue,
// allowing clipping within the CSS Color 4 just-noticeable difference.
func GamutMapLinear(rgb Vec3) Vec3 {
	if InGamutLinear(rgb) {
		return clipLinear(rgb)
	}
	if math.IsNaN(rgb[0]) || math.IsNaN(rgb[1]) || math.IsNaN(rgb[2]) {
		return Vec3{}
	}
	return GamutMapOKLCH(ToPolar(LinearSRGBToOKLab(rgb), AchromaticOKLCH))
}

// GamutMapOKLCH maps an OKLCH color into the sRGB gamut (CSS Color 4 algorithm) and
// returns in-gamut linear sRGB.
func GamutMapOKLCH(lch Vec3) Vec3 {
	const (
		jnd = 0.02
		eps = 0.0001
	)
	switch {
	case math.IsNaN(lch[0]) || math.IsNaN(lch[1]) || math.IsNaN(lch[2]):
		return Vec3{}
	case lch[0] >= 1:
		return Vec3{1, 1, 1}
	case lch[0] <= 0:
		return Vec3{}
	}
	origin := OKLabToLinearSRGB(FromPolar(lch))
	if InGamutLinear(origin) {
		return clipLinear(origin)
	}
	toLab := LinearSRGBToOKLab
	current := lch
	clipped := clipLinear(origin)
	if Vec3Distance(toLab(clipped), FromPolar(current)) < jnd {
		return clipped
	}
	lo, hi := 0.0, lch[1]
	loInGamut := true
	for hi-lo > eps {
		current[1] = (lo + hi) / 2
		rgb := OKLabToLinearSRGB(FromPolar(current))
		if loInGamut && InGamutLinear(rgb) {
			lo = current[1]
			continue
		}
		clipped = clipLinear(rgb)
		e := Vec3Distance(toLab(clipped), FromPolar(current))
		if e < jnd {
			if jnd-e < eps {
				return clipped
			}
			loInGamut = false
			lo = current[1]
		} else {
			hi = current[1]
		}
	}
	current[1] = lo
	return clipLinear(OKLabToLinearSRGB(FromPolar(current)))
}
