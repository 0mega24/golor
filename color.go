// Package golor provides an sRGB Color type with constructors, accessors, and hex parsing.
package golor

import (
	"encoding/json"
	"fmt"
	stdcolor "image/color"
	"math"
	"strconv"
	"strings"
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
// the range are clamped.
func RGBf(r, g, b float64) Color {
	return RGBAf(r, g, b, 1)
}

// RGBAf creates a Color from normalized float64 values (0.0-1.0). Values outside
// the range are clamped.
func RGBAf(r, g, b, a float64) Color {
	return Color{clamp01(r), clamp01(g), clamp01(b), clamp01(a)}
}

// Hex parses a hex string (#rrggbb, rrggbb, #rrggbbaa, or rrggbbaa) into a Color.
func Hex(s string) (Color, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 && len(s) != 8 {
		return Color{}, fmt.Errorf("golor: invalid hex color %q", s)
	}
	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("golor: invalid hex color %q", s)
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("golor: invalid hex color %q", s)
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("golor: invalid hex color %q", s)
	}
	a := uint64(255)
	if len(s) == 8 {
		a, err = strconv.ParseUint(s[6:8], 16, 8)
		if err != nil {
			return Color{}, fmt.Errorf("golor: invalid hex color %q", s)
		}
	}
	return RGBA(uint8(r), uint8(g), uint8(b), uint8(a)), nil
}

// String returns the color as a lowercase #rrggbb hex string for opaque colors
// and #rrggbbaa for colors with alpha.
func (c Color) String() string {
	if c.A8() < 255 {
		return fmt.Sprintf("#%02x%02x%02x%02x", c.R8(), c.G8(), c.B8(), c.A8())
	}
	return fmt.Sprintf("#%02x%02x%02x", c.R8(), c.G8(), c.B8())
}

// RGBA returns c as alpha-premultiplied 16-bit channels, satisfying image/color.Color.
func (c Color) RGBA() (r, g, b, a uint32) {
	alpha := uint32(math.Round(clamp01(c.A) * 0xffff))
	return premul16(c.R, alpha), premul16(c.G, alpha), premul16(c.B, alpha), alpha
}

// FromStdColor converts any standard library color.Color into a straight-alpha Color.
func FromStdColor(c stdcolor.Color) Color {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return RGBAf(0, 0, 0, 0)
	}
	return RGBAf(float64(r)/float64(a), float64(g)/float64(a), float64(b)/float64(a), float64(a)/0xffff)
}

// MarshalJSON serializes c as 8-bit channel values.
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
func (c Color) R8() uint8 { return uint8(math.Round(clamp01(c.R) * 255)) }

// G8 returns the green channel as an 8-bit value (0-255).
func (c Color) G8() uint8 { return uint8(math.Round(clamp01(c.G) * 255)) }

// B8 returns the blue channel as an 8-bit value (0-255).
func (c Color) B8() uint8 { return uint8(math.Round(clamp01(c.B) * 255)) }

// A8 returns the alpha channel as an 8-bit value (0-255).
func (c Color) A8() uint8 { return uint8(math.Round(clamp01(c.A) * 255)) }

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func premul16(v float64, a uint32) uint32 {
	return uint32(math.Round(clamp01(v) * float64(a)))
}

func valid8(v int) bool {
	return v >= 0 && v <= 255
}
