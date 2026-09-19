package ansi_test

import (
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/ansi"
)

func TestTruecolorEscapes(t *testing.T) {
	c := golor.RGB(255, 107, 53)
	if got := ansi.Foreground(c); got != "\x1b[38;2;255;107;53m" {
		t.Fatalf("Foreground() = %q", got)
	}
	if got := ansi.Background(c); got != "\x1b[48;2;255;107;53m" {
		t.Fatalf("Background() = %q", got)
	}
}

func TestNearest256AndEscapes(t *testing.T) {
	red := golor.RGB(255, 0, 0)
	if got := ansi.Nearest256(red); got != 9 {
		t.Fatalf("Nearest256(red) = %d, want 9", got)
	}
	if got := ansi.Foreground256(red); got != "\x1b[38;5;9m" {
		t.Fatalf("Foreground256(red) = %q", got)
	}
	if got := ansi.Background256(red); got != "\x1b[48;5;9m" {
		t.Fatalf("Background256(red) = %q", got)
	}

	white := golor.RGB(255, 255, 255)
	if got := ansi.Nearest256(white); got != 15 {
		t.Fatalf("Nearest256(white) = %d, want 15", got)
	}
	gray := golor.RGB(238, 238, 238)
	if got := ansi.Nearest256(gray); got < 232 {
		t.Fatalf("Nearest256(gray) = %d, want grayscale index", got)
	}
}

func TestSupportsTruecolorEnv(t *testing.T) {
	cases := []struct {
		name      string
		colorterm string
		term      string
		want      bool
	}{
		{"colorterm truecolor", "truecolor", "xterm-256color", true},
		{"colorterm 24bit", "24bit", "xterm-256color", true},
		{"term direct", "", "xterm-direct", true},
		{"kitty", "", "xterm-kitty", true},
		{"plain 256", "", "xterm-256color", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ansi.SupportsTruecolorEnv(tc.colorterm, tc.term); got != tc.want {
				t.Fatalf("SupportsTruecolorEnv() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDetectColorMode(t *testing.T) {
	t.Setenv("COLORTERM", "truecolor")
	t.Setenv("TERM", "xterm-256color")
	if got := ansi.DetectColorMode(); got != ansi.Truecolor {
		t.Fatalf("DetectColorMode() = %v, want Truecolor", got)
	}

	t.Setenv("COLORTERM", "")
	t.Setenv("TERM", "xterm-256color")
	if got := ansi.DetectColorMode(); got != ansi.Color256 {
		t.Fatalf("DetectColorMode() = %v, want Color256", got)
	}
}

func ExampleBackground() {
	red := golor.RGB(255, 0, 0)
	_ = ansi.Background(red) + "  " + ansi.Reset
}
