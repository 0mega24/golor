// Package ansi returns ANSI escape sequences without writing to a terminal.
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

// Nearest256 finds the nearest xterm-256 cube/grey entry by RGB distance.
// Theme-dependent indexes 0-15 are excluded.
func Nearest256(c golor.Color) int {
	r, g, b := int(c.R8()), int(c.G8()), int(c.B8())
	ri, gi, bi := nearestCube(r), nearestCube(g), nearestCube(b)
	cubeIndex := 16 + 36*ri + 6*gi + bi
	cubeDist := square(r-colorCube[ri]) + square(g-colorCube[gi]) + square(b-colorCube[bi])

	// The closest grey level to (r,g,b) is nearest to their mean.
	k := min(max((r+g+b-24+15)/30, 0), 23)
	grayIndex, grayDist := 0, 1<<62
	for _, kk := range []int{k - 1, k, k + 1} {
		if kk < 0 || kk > 23 {
			continue
		}
		v := 8 + 10*kk
		if d := square(r-v) + square(g-v) + square(b-v); d < grayDist {
			grayIndex, grayDist = 232+kk, d
		}
	}
	if grayDist < cubeDist {
		return grayIndex
	}
	return cubeIndex
}

// nearestCube returns the index of the cube level closest to v, breaking ties with the lower level.
func nearestCube(v int) int {
	best, bestDist := 0, 1<<62
	for i, level := range colorCube {
		if d := square(v - level); d < bestDist {
			best, bestDist = i, d
		}
	}
	return best
}

func square(v int) int {
	return v * v
}

var colorCube = [6]int{0, 95, 135, 175, 215, 255}
