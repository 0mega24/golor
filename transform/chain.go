// Package transform provides a fluent Builder for chaining color transformations.
package transform

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/adjust"
)

// Chain holds a color and allows fluent multi-step transformations.
type Chain struct {
	c golor.Color
}

// From creates a new Chain starting from c.
func From(c golor.Color) *Chain {
	return &Chain{c: c}
}

// Color unwraps the current color from the chain.
func (ch *Chain) Color() golor.Color {
	return ch.c
}

// Lighten increases HSL lightness by amount.
func (ch *Chain) Lighten(amount float64) *Chain {
	ch.c = adjust.Lighten(ch.c, amount)
	return ch
}

// Darken decreases HSL lightness by amount.
func (ch *Chain) Darken(amount float64) *Chain {
	ch.c = adjust.Darken(ch.c, amount)
	return ch
}

// Saturate increases HSL saturation by amount.
func (ch *Chain) Saturate(amount float64) *Chain {
	ch.c = adjust.Saturate(ch.c, amount)
	return ch
}

// Desaturate decreases HSL saturation by amount.
func (ch *Chain) Desaturate(amount float64) *Chain {
	ch.c = adjust.Desaturate(ch.c, amount)
	return ch
}

// ShiftHue rotates the hue by degrees.
func (ch *Chain) ShiftHue(degrees float64) *Chain {
	ch.c = adjust.ShiftHue(ch.c, degrees)
	return ch
}

// Warm shifts the hue toward orange (~30°) by amount.
func (ch *Chain) Warm(amount float64) *Chain {
	ch.c = adjust.Warm(ch.c, amount)
	return ch
}

// Cool shifts the hue toward blue (~210°) by amount.
func (ch *Chain) Cool(amount float64) *Chain {
	ch.c = adjust.Cool(ch.c, amount)
	return ch
}

// Tint mixes toward white by amount.
func (ch *Chain) Tint(amount float64) *Chain {
	ch.c = adjust.Tint(ch.c, amount)
	return ch
}

// Shade mixes toward black by amount.
func (ch *Chain) Shade(amount float64) *Chain {
	ch.c = adjust.Shade(ch.c, amount)
	return ch
}

// SetLightness sets HSL lightness to l.
func (ch *Chain) SetLightness(l float64) *Chain {
	ch.c = adjust.SetLightness(ch.c, l)
	return ch
}

// SetSaturation sets HSL saturation to s.
func (ch *Chain) SetSaturation(s float64) *Chain {
	ch.c = adjust.SetSaturation(ch.c, s)
	return ch
}

// SetHue sets the hue to h degrees.
func (ch *Chain) SetHue(h float64) *Chain {
	ch.c = adjust.SetHue(ch.c, h)
	return ch
}

// SetValue sets HSV value to v.
func (ch *Chain) SetValue(v float64) *Chain {
	ch.c = adjust.SetValue(ch.c, v)
	return ch
}

// Mix applies a blend mode function with other as the layer color.
// The current color is used as the base.
func (ch *Chain) Mix(other golor.Color, mode func(golor.Color, golor.Color) golor.Color) *Chain {
	ch.c = mode(ch.c, other)
	return ch
}
