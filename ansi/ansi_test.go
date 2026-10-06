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
	if got := ansi.Nearest256(red); got != 196 {
		t.Fatalf("Nearest256(red) = %d, want 196 (cube red, not theme-defined 9)", got)
	}
	if got := ansi.Foreground256(red); got != "\x1b[38;5;196m" {
		t.Fatalf("Foreground256(red) = %q", got)
	}
	if got := ansi.Background256(red); got != "\x1b[48;5;196m" {
		t.Fatalf("Background256(red) = %q", got)
	}

	white := golor.RGB(255, 255, 255)
	if got := ansi.Nearest256(white); got != 231 {
		t.Fatalf("Nearest256(white) = %d, want 231", got)
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

// TestNearest256MatchesBruteForce checks Nearest256 against a brute-force search over indexes 16-255.
func TestNearest256MatchesBruteForce(t *testing.T) {
	step := 5
	for r := 0; r < 256; r += step {
		for g := 0; g < 256; g += step {
			for b := 0; b < 256; b += 3 {
				c := golor.RGB(uint8(r), uint8(g), uint8(b))
				best, bestD := 0, 1<<62
				for i := 16; i < 256; i++ {
					pr, pg, pb := paletteRGB(i)
					if d := (r-pr)*(r-pr) + (g-pg)*(g-pg) + (b-pb)*(b-pb); d < bestD {
						best, bestD = i, d
					}
				}
				if got := ansi.Nearest256(c); got != best {
					t.Fatalf("Nearest256(%s) = %d, want %d", c, got, best)
				}
			}
		}
	}
}

func paletteRGB(i int) (int, int, int) {
	cube := [6]int{0, 95, 135, 175, 215, 255}
	if i < 232 {
		i -= 16
		return cube[i/36], cube[(i/6)%6], cube[i%6]
	}
	v := 8 + (i-232)*10
	return v, v, v
}

func BenchmarkNearest256(b *testing.B) {
	c := golor.RGB(200, 100, 50)
	for i := 0; i < b.N; i++ {
		sinkIndex = ansi.Nearest256(c)
	}
}

func BenchmarkForeground(b *testing.B) {
	c := golor.RGB(200, 100, 50)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = ansi.Foreground(c)
	}
}

var (
	sinkIndex  int
	sinkString string
)
