// Package golor provides an sRGB Color type with constructors, accessors, and hex parsing.
package golor

import (
	"encoding/json"
	"fmt"
	stdcolor "image/color"
	"math"
	"strings"

	"github.com/0mega24/golor/v2/internal/space"
)

// Color represents a straight-alpha sRGB color with float64 components normalized to [0.0, 1.0].
type Color struct {
	R, G, B, A float64
}

// RGB creates an opaque Color from 8-bit RGB values (0-255).
func RGB(r, g, b uint8) Color {
	return RGBA(r, g, b, 255)
}

// RGBA creates a Color from 8-bit RGBA values (0-255).
func RGBA(r, g, b, a uint8) Color {
	return Color{
		R: float64(r) / 255,
		G: float64(g) / 255,
		B: float64(b) / 255,
		A: float64(a) / 255,
	}
}

// RGBf creates an opaque Color from normalized float64 values (0.0-1.0). Values outside
// the range are clamped; NaN becomes 0.
func RGBf(r, g, b float64) Color {
	return RGBAf(r, g, b, 1)
}

// RGBAf creates a Color from normalized float64 values (0.0-1.0). Values outside
// the range are clamped; NaN becomes 0.
func RGBAf(r, g, b, a float64) Color {
	return Color{space.Clamp01(r), space.Clamp01(g), space.Clamp01(b), space.Clamp01(a)}
}

// Hex parses a CSS hex color: #rgb, #rgba, #rrggbb, or #rrggbbaa. The leading # is
// optional, letters are case-insensitive, and surrounding whitespace is ignored.
func Hex(s string) (Color, error) {
	h := strings.TrimSpace(s)
	if len(h) > 0 && h[0] == '#' {
		h = h[1:]
	}
	var ch [4]uint8
	ch[3] = 255
	switch len(h) {
	case 3, 4:
		for i := 0; i < len(h); i++ {
			v, ok := fromHexDigit(h[i])
			if !ok {
				return Color{}, fmt.Errorf("golor: invalid hex color %q", s)
			}
			ch[i] = v * 17
		}
	case 6, 8:
		for i := 0; i < len(h); i += 2 {
			hi, ok1 := fromHexDigit(h[i])
			lo, ok2 := fromHexDigit(h[i+1])
			if !ok1 || !ok2 {
				return Color{}, fmt.Errorf("golor: invalid hex color %q", s)
			}
			ch[i/2] = hi<<4 | lo
		}
	default:
		return Color{}, fmt.Errorf("golor: invalid hex color %q", s)
	}
	return RGBA(ch[0], ch[1], ch[2], ch[3]), nil
}

func fromHexDigit(b byte) (uint8, bool) {
	switch {
	case b >= '0' && b <= '9':
		return b - '0', true
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10, true
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10, true
	default:
		return 0, false
	}
}

const hexDigits = "0123456789abcdef"

// String returns lowercase #rrggbb or #rrggbbaa, rounding channels to 8 bits.
func (c Color) String() string {
	var buf [9]byte
	buf[0] = '#'
	for i, v := range [4]uint8{c.R8(), c.G8(), c.B8(), c.A8()} {
		buf[1+2*i] = hexDigits[v>>4]
		buf[2+2*i] = hexDigits[v&0x0f]
	}
	if buf[7] == 'f' && buf[8] == 'f' {
		return string(buf[:7])
	}
	return string(buf[:])
}

// RGBA returns c as alpha-premultiplied 16-bit channels, satisfying image/color.Color.
func (c Color) RGBA() (r, g, b, a uint32) {
	alpha := uint32(math.Round(space.Clamp01(c.A) * 0xffff))
	return premul16(c.R, alpha), premul16(c.G, alpha), premul16(c.B, alpha), alpha
}

// FromStdColor converts standard colors to straight alpha; NRGBA/NRGBA64/Color are lossless.
// Fully transparent colors become transparent black.
func FromStdColor(c stdcolor.Color) Color {
	switch v := c.(type) {
	case Color:
		if space.Clamp01(v.A) == 0 {
			return Color{}
		}
		return RGBAf(v.R, v.G, v.B, v.A)
	case stdcolor.NRGBA:
		if v.A == 0 {
			return Color{}
		}
		return RGBA(v.R, v.G, v.B, v.A)
	case stdcolor.NRGBA64:
		if v.A == 0 {
			return Color{}
		}
		return RGBAf(float64(v.R)/0xffff, float64(v.G)/0xffff, float64(v.B)/0xffff, float64(v.A)/0xffff)
	}
	r, g, b, a := c.RGBA()
	if a == 0 {
		return RGBAf(0, 0, 0, 0)
	}
	return RGBAf(float64(r)/float64(a), float64(g)/float64(a), float64(b)/float64(a), float64(a)/0xffff)
}

// MarshalJSON returns rounded 8-bit channels, e.g. {"r":255,"g":107,"b":53,"a":255}.
func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		R uint8 `json:"r"`
		G uint8 `json:"g"`
		B uint8 `json:"b"`
		A uint8 `json:"a"`
	}{
		R: c.R8(),
		G: c.G8(),
		B: c.B8(),
		A: c.A8(),
	})
}

// UnmarshalJSON deserializes c from 8-bit channel values.
func (c *Color) UnmarshalJSON(data []byte) error {
	var channels struct {
		R *int `json:"r"`
		G *int `json:"g"`
		B *int `json:"b"`
		A *int `json:"a"`
	}
	if err := json.Unmarshal(data, &channels); err != nil {
		return err
	}
	if channels.R == nil || channels.G == nil || channels.B == nil || channels.A == nil {
		return fmt.Errorf("golor: invalid JSON color channels")
	}
	if !valid8(*channels.R) || !valid8(*channels.G) || !valid8(*channels.B) || !valid8(*channels.A) {
		return fmt.Errorf("golor: invalid JSON color channels")
	}
	*c = RGBA(uint8(*channels.R), uint8(*channels.G), uint8(*channels.B), uint8(*channels.A))
	return nil
}

// R8 returns the red channel as an 8-bit value (0-255).
func (c Color) R8() uint8 { return uint8(math.Round(space.Clamp01(c.R) * 255)) }

// G8 returns the green channel as an 8-bit value (0-255).
func (c Color) G8() uint8 { return uint8(math.Round(space.Clamp01(c.G) * 255)) }

// B8 returns the blue channel as an 8-bit value (0-255).
func (c Color) B8() uint8 { return uint8(math.Round(space.Clamp01(c.B) * 255)) }

// A8 returns the alpha channel as an 8-bit value (0-255).
func (c Color) A8() uint8 { return uint8(math.Round(space.Clamp01(c.A) * 255)) }

func premul16(v float64, a uint32) uint32 {
	return uint32(math.Round(space.Clamp01(v) * float64(a)))
}

func valid8(v int) bool {
	return v >= 0 && v <= 255
}
