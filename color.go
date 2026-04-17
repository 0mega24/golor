// Package golor provides an sRGB Color type with constructors, accessors, and hex parsing.
package golor

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color represents an sRGB color with float64 components normalized to [0.0, 1.0].
type Color struct {
	R, G, B float64
}

// RGB creates a Color from 8-bit RGB values (0–255).
func RGB(r, g, b uint8) Color {
	return Color{float64(r) / 255, float64(g) / 255, float64(b) / 255}
}

// RGBf creates a Color from normalized float64 values (0.0–1.0). Values outside
// the range are clamped.
func RGBf(r, g, b float64) Color {
	return Color{clamp01(r), clamp01(g), clamp01(b)}
}

// Hex parses a hex string (#rrggbb or rrggbb) into a Color.
func Hex(s string) (Color, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
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
	return RGB(uint8(r), uint8(g), uint8(b)), nil
}

// String returns the color as a lowercase #rrggbb hex string.
func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R8(), c.G8(), c.B8())
}

// R8 returns the red channel as an 8-bit value (0–255).
func (c Color) R8() uint8 { return uint8(math.Round(c.R * 255)) }

// G8 returns the green channel as an 8-bit value (0–255).
func (c Color) G8() uint8 { return uint8(math.Round(c.G * 255)) }

// B8 returns the blue channel as an 8-bit value (0–255).
func (c Color) B8() uint8 { return uint8(math.Round(c.B * 255)) }

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
