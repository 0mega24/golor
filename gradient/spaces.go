package gradient

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/internal/space"
)

func rgbMixer(a, b golor.Color) mixer {
	return rect(a, b,
		func(c golor.Color) space.Vec3 {
			return space.Vec3{space.Clamp01(c.R), space.Clamp01(c.G), space.Clamp01(c.B)}
		},
		func(v space.Vec3) golor.Color { return golor.RGBf(v[0], v[1], v[2]) })
}

func hslMixer(a, b golor.Color) mixer {
	return polar(a, b,
		func(c golor.Color) polarColor {
			h := convert.ToHSL(c)
			return polarColor{space.Vec3{h.L, h.S, h.H}, h.IsAchromatic()}
		},
		func(v space.Vec3) golor.Color { return convert.FromHSL(convert.HSL{H: v[2], S: v[1], L: v[0]}) })
}

func labMixer(a, b golor.Color) mixer {
	return rect(a, b,
		func(c golor.Color) space.Vec3 { l := convert.ToLAB(c); return space.Vec3{l.L, l.A, l.B} },
		func(v space.Vec3) golor.Color { return convert.FromLABMapped(convert.LAB{L: v[0], A: v[1], B: v[2]}) })
}

func lchMixer(a, b golor.Color) mixer {
	return polar(a, b,
		func(c golor.Color) polarColor {
			l := convert.ToLCH(c)
			return polarColor{space.Vec3{l.L, l.C, l.H}, l.IsAchromatic()}
		},
		func(v space.Vec3) golor.Color { return convert.FromLCHMapped(convert.LCH{L: v[0], C: v[1], H: v[2]}) })
}

func oklabMixer(a, b golor.Color) mixer {
	return rect(a, b,
		func(c golor.Color) space.Vec3 { l := convert.ToOKLab(c); return space.Vec3{l.L, l.A, l.B} },
		func(v space.Vec3) golor.Color {
			return convert.FromOKLabMapped(convert.OKLab{L: v[0], A: v[1], B: v[2]})
		})
}

func oklchMixer(a, b golor.Color) mixer {
	return polar(a, b,
		func(c golor.Color) polarColor {
			l := convert.ToOKLCH(c)
			return polarColor{space.Vec3{l.L, l.C, l.H}, l.IsAchromatic()}
		},
		func(v space.Vec3) golor.Color {
			return convert.FromOKLCHMapped(convert.OKLCH{L: v[0], C: v[1], H: v[2]})
		})
}

// RGB returns n evenly spaced colors interpolated in sRGB between a and b.
// Both endpoints are included; n == 1 returns only a. Returns nil if n <= 0.
func RGB(a, b golor.Color, n int) []golor.Color { return sample(n, a, rgbMixer(a, b)) }

// HSL returns n evenly spaced colors interpolated in HSL between a and b.
// Hue takes the shorter arc; a grey endpoint borrows the other endpoint's hue.
func HSL(a, b golor.Color, n int) []golor.Color { return sample(n, a, hslMixer(a, b)) }

// LAB returns n evenly spaced colors interpolated in CIE L*a*b* between a and b.
func LAB(a, b golor.Color, n int) []golor.Color { return sample(n, a, labMixer(a, b)) }

// LCH returns n evenly spaced colors interpolated in CIE LCH between a and b.
// Hue takes the shorter arc; a grey endpoint borrows the other endpoint's hue.
func LCH(a, b golor.Color, n int) []golor.Color { return sample(n, a, lchMixer(a, b)) }

// OKLab returns n evenly spaced colors interpolated in OKLab between a and b. This is
// the CSS default interpolation space and gives perceptually even steps.
func OKLab(a, b golor.Color, n int) []golor.Color { return sample(n, a, oklabMixer(a, b)) }

// OKLCH returns n evenly spaced colors interpolated in OKLCH between a and b.
// Hue takes the shorter arc; a grey endpoint borrows the other endpoint's hue.
func OKLCH(a, b golor.Color, n int) []golor.Color { return sample(n, a, oklchMixer(a, b)) }

// MultiStopRGB interpolates through all stops in sRGB, returning n total colors.
// Both outer endpoints are included. Returns nil if n <= 0 or stops is empty.
func MultiStopRGB(stops []golor.Color, n int) []golor.Color {
	return multiStop(stops, n, rgbMixer)
}

// MultiStopHSL interpolates through all stops in HSL, returning n total colors.
func MultiStopHSL(stops []golor.Color, n int) []golor.Color {
	return multiStop(stops, n, hslMixer)
}

// MultiStopLAB interpolates through all stops in CIE L*a*b*, returning n total colors.
func MultiStopLAB(stops []golor.Color, n int) []golor.Color {
	return multiStop(stops, n, labMixer)
}

// MultiStopLCH interpolates through all stops in CIE LCH, returning n total colors.
// Both outer endpoints are included. Returns nil if n <= 0 or stops is empty.
func MultiStopLCH(stops []golor.Color, n int) []golor.Color {
	return multiStop(stops, n, lchMixer)
}

// MultiStopOKLab interpolates through all stops in OKLab, returning n total colors.
func MultiStopOKLab(stops []golor.Color, n int) []golor.Color {
	return multiStop(stops, n, oklabMixer)
}

// MultiStopOKLCH interpolates through all stops in OKLCH, returning n total colors.
func MultiStopOKLCH(stops []golor.Color, n int) []golor.Color {
	return multiStop(stops, n, oklchMixer)
}
