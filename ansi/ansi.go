// Package ansi generates ANSI terminal escape codes for golor colors.
//
// It does not write to stdout or stderr; callers compose the returned strings
// into their own terminal, CLI, or TUI output.
package ansi

import (
	"fmt"
	"os"
	"strings"

	"github.com/0mega24/golor/v2"
)

// Reset clears ANSI styling.
const Reset = "\x1b[0m"

// ColorMode describes the terminal color escape style to use.
type ColorMode int

const (
	// AutoColor asks callers to choose a mode from terminal capability detection.
	AutoColor ColorMode = iota
	// Truecolor emits 24-bit ANSI escape codes.
	Truecolor
	// Color256 emits xterm-256 palette ANSI escape codes.
	Color256
)

// Foreground returns a 24-bit truecolor foreground escape code for c.
func Foreground(c golor.Color) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R8(), c.G8(), c.B8())
}

// Background returns a 24-bit truecolor background escape code for c.
func Background(c golor.Color) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R8(), c.G8(), c.B8())
}

// Foreground256 returns an xterm-256 foreground escape code for c.
func Foreground256(c golor.Color) string {
	return fmt.Sprintf("\x1b[38;5;%dm", Nearest256(c))
}

// Background256 returns an xterm-256 background escape code for c.
func Background256(c golor.Color) string {
	return fmt.Sprintf("\x1b[48;5;%dm", Nearest256(c))
}

// SupportsTruecolor reports whether the current terminal environment advertises truecolor support.
func SupportsTruecolor() bool {
	return SupportsTruecolorEnv(os.Getenv("COLORTERM"), os.Getenv("TERM"))
}

// SupportsTruecolorEnv reports whether colorterm/term values advertise truecolor support.
func SupportsTruecolorEnv(colorterm, term string) bool {
	colorterm = strings.ToLower(colorterm)
	term = strings.ToLower(term)
	if strings.Contains(colorterm, "truecolor") || strings.Contains(colorterm, "24bit") {
		return true
	}
	return strings.Contains(term, "truecolor") ||
		strings.Contains(term, "24bit") ||
		strings.Contains(term, "direct") ||
		strings.HasPrefix(term, "wezterm") ||
		strings.HasPrefix(term, "xterm-kitty")
}

// DetectColorMode returns Truecolor when the current terminal advertises 24-bit support,
// otherwise Color256.
func DetectColorMode() ColorMode {
	if SupportsTruecolor() {
		return Truecolor
	}
	return Color256
}

// Nearest256 returns the closest xterm-256 palette index to c using Euclidean RGB distance.
func Nearest256(c golor.Color) int {
	r, g, b := int(c.R8()), int(c.G8()), int(c.B8())
	bestIndex := 0
	bestDistance := 1 << 62
	for i := 0; i < 256; i++ {
		pr, pg, pb := xterm256RGB(i)
		distance := square(r-pr) + square(g-pg) + square(b-pb)
		if distance < bestDistance {
			bestDistance = distance
			bestIndex = i
		}
	}
	return bestIndex
}

func square(v int) int {
	return v * v
}

func xterm256RGB(index int) (r, g, b int) {
	if index < 0 {
		index = 0
	}
	if index > 255 {
		index = 255
	}
	if index < 16 {
		return standard16[index][0], standard16[index][1], standard16[index][2]
	}
	if index < 232 {
		i := index - 16
		return colorCube[i/36], colorCube[(i/6)%6], colorCube[i%6]
	}
	gray := 8 + (index-232)*10
	return gray, gray, gray
}

var colorCube = [6]int{0, 95, 135, 175, 215, 255}

var standard16 = [16][3]int{
	{0, 0, 0},
	{128, 0, 0},
	{0, 128, 0},
	{128, 128, 0},
	{0, 0, 128},
	{128, 0, 128},
	{0, 128, 128},
	{192, 192, 192},
	{128, 128, 128},
	{255, 0, 0},
	{0, 255, 0},
	{255, 255, 0},
	{0, 0, 255},
	{255, 0, 255},
	{0, 255, 255},
	{255, 255, 255},
}
