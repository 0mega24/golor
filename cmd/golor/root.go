package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/ansi"
	"github.com/0mega24/golor/v2/contrast"
	"github.com/0mega24/golor/v2/convert"
)

type cliOptions struct {
	json bool
}

func newRootCommand(stdout, stderr io.Writer) *cobra.Command {
	opts := &cliOptions{}
	root := &cobra.Command{
		Use:           "golor",
		Short:         "Color conversion, contrast, and terminal preview tooling",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().BoolVar(&opts.json, "json", false, "write JSON output")
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.AddCommand(newConvertCommand(stdout, opts))
	root.AddCommand(newContrastCommand(stdout, opts))
	root.AddCommand(newPreviewCommand(stdout, opts))
	root.AddCommand(newPaletteCommand(stdout, opts))
	return root
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func parseColor(input, from string) (golor.Color, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return golor.Color{}, fmt.Errorf("empty color")
	}
	switch strings.ToLower(from) {
	case "", "auto":
		if c, err := golor.Hex(input); err == nil {
			return c, nil
		}
		return golor.ParseCSS(input)
	case "hex":
		return golor.Hex(input)
	case "rgb", "rgba":
		return golor.ParseRGB(input)
	case "hsl", "hsla":
		return golor.ParseHSL(input)
	case "named", "name":
		return golor.Named(input)
	case "css":
		return golor.ParseCSS(input)
	default:
		return golor.Color{}, fmt.Errorf("unknown input format %q", from)
	}
}

func readColorArg(args []string, stdin io.Reader) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", err
	}
	input := strings.TrimSpace(string(data))
	if input == "" {
		return "", fmt.Errorf("missing color argument")
	}
	return input, nil
}

func newConvertCommand(stdout io.Writer, opts *cliOptions) *cobra.Command {
	var from, to string
	cmd := &cobra.Command{
		Use:   "convert [color]",
		Short: "Convert a color to another format",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := readColorArg(args, cmd.InOrStdin())
			if err != nil {
				return err
			}
			c, err := parseColor(input, from)
			if err != nil {
				return err
			}
			return printConverted(stdout, opts.json, c, to)
		},
	}
	cmd.Flags().StringVar(&from, "from", "auto", "input format: auto, hex, rgb, hsl, named, css")
	cmd.Flags().StringVar(&to, "to", "hex", "output format: hex, rgb, hsl, hsv, lab, lch, cmyk")
	return cmd
}

func printConverted(w io.Writer, jsonOut bool, c golor.Color, to string) error {
	switch strings.ToLower(to) {
	case "hex":
		if jsonOut {
			return writeJSON(w, map[string]string{"hex": c.String()})
		}
		_, err := fmt.Fprintln(w, c.String())
		return err
	case "rgb", "rgba":
		if jsonOut {
			return writeJSON(w, c)
		}
		_, err := fmt.Fprintf(w, "RGB{R:%d G:%d B:%d A:%d}\n", c.R8(), c.G8(), c.B8(), c.A8())
		return err
	case "hsl":
		hsl := convert.ToHSL(c)
		if jsonOut {
			return writeJSON(w, map[string]float64{"h": hsl.H, "s": hsl.S, "l": hsl.L})
		}
		_, err := fmt.Fprintf(w, "HSL{H:%.2f S:%.4f L:%.4f}\n", hsl.H, hsl.S, hsl.L)
		return err
	case "hsv":
		hsv := convert.ToHSV(c)
		if jsonOut {
			return writeJSON(w, map[string]float64{"h": hsv.H, "s": hsv.S, "v": hsv.V})
		}
		_, err := fmt.Fprintf(w, "HSV{H:%.2f S:%.4f V:%.4f}\n", hsv.H, hsv.S, hsv.V)
		return err
	case "lab":
		lab := convert.ToLAB(c)
		if jsonOut {
			return writeJSON(w, map[string]float64{"l": lab.L, "a": lab.A, "b": lab.B})
		}
		_, err := fmt.Fprintf(w, "LAB{L:%.4f A:%.4f B:%.4f}\n", lab.L, lab.A, lab.B)
		return err
	case "lch":
		lch := convert.ToLCH(c)
		if jsonOut {
			return writeJSON(w, map[string]float64{"l": lch.L, "c": lch.C, "h": lch.H})
		}
		_, err := fmt.Fprintf(w, "LCH{L:%.4f C:%.4f H:%.2f}\n", lch.L, lch.C, lch.H)
		return err
	case "cmyk":
		cmyk := convert.ToCMYK(c)
		if jsonOut {
			return writeJSON(w, map[string]float64{"c": cmyk.C, "m": cmyk.M, "y": cmyk.Y, "k": cmyk.K})
		}
		_, err := fmt.Fprintf(w, "CMYK{C:%.4f M:%.4f Y:%.4f K:%.4f}\n", cmyk.C, cmyk.M, cmyk.Y, cmyk.K)
		return err
	default:
		return fmt.Errorf("unknown output format %q", to)
	}
}

func newContrastCommand(stdout io.Writer, opts *cliOptions) *cobra.Command {
	var enforce float64
	cmd := &cobra.Command{
		Use:   "contrast <foreground> <background>",
		Short: "Calculate WCAG contrast between two colors",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fg, err := parseColor(args[0], "auto")
			if err != nil {
				return err
			}
			bg, err := parseColor(args[1], "auto")
			if err != nil {
				return err
			}
			ratio := contrast.Ratio(fg, bg)
			result := contrastResult{
				Ratio: math.Round(ratio*100) / 100,
				AA:    ratio >= 4.5,
				AAA:   ratio >= 7.0,
			}
			if enforce > 0 {
				adjusted := contrast.EnforceContrast(fg, bg, enforce)
				result.Adjusted = &adjusted
			}
			if opts.json {
				return writeJSON(stdout, result)
			}
			if result.Adjusted != nil {
				_, err = fmt.Fprintf(stdout, "ratio %.2f AA:%t AAA:%t adjusted:%s\n", result.Ratio, result.AA, result.AAA, result.Adjusted.String())
				return err
			}
			_, err = fmt.Fprintf(stdout, "ratio %.2f AA:%t AAA:%t\n", result.Ratio, result.AA, result.AAA)
			return err
		},
	}
	cmd.Flags().Float64Var(&enforce, "enforce", 0, "print an adjusted foreground color meeting this ratio")
	return cmd
}

type contrastResult struct {
	Ratio    float64      `json:"ratio"`
	AA       bool         `json:"aa"`
	AAA      bool         `json:"aaa"`
	Adjusted *golor.Color `json:"adjusted,omitempty"`
}

func newPreviewCommand(stdout io.Writer, opts *cliOptions) *cobra.Command {
	var mode string
	cmd := &cobra.Command{
		Use:   "preview <color> [color...]",
		Short: "Preview terminal color swatches",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			colors := make([]golor.Color, 0, len(args))
			for _, arg := range args {
				c, err := parseColor(arg, "auto")
				if err != nil {
					return err
				}
				colors = append(colors, c)
			}
			colorMode, err := parseColorMode(mode)
			if err != nil {
				return err
			}
			return printSwatches(stdout, opts.json, colors, colorMode)
		},
	}
	cmd.Flags().StringVar(&mode, "color-mode", "auto", "color mode: auto, truecolor, 256")
	return cmd
}

func parseColorMode(mode string) (ansi.ColorMode, error) {
	switch strings.ToLower(mode) {
	case "", "auto":
		return ansi.DetectColorMode(), nil
	case "truecolor", "24bit":
		return ansi.Truecolor, nil
	case "256", "256color":
		return ansi.Color256, nil
	default:
		return ansi.AutoColor, fmt.Errorf("unknown color mode %q", mode)
	}
}

func printSwatches(w io.Writer, jsonOut bool, colors []golor.Color, mode ansi.ColorMode) error {
	if jsonOut {
		items := make([]swatchResult, len(colors))
		for i, c := range colors {
			items[i] = swatchResult{Hex: c.String(), Color: c, Index256: ansi.Nearest256(c)}
		}
		return writeJSON(w, items)
	}
	for _, c := range colors {
		bg := ansi.Background(c)
		if mode == ansi.Color256 {
			bg = ansi.Background256(c)
		}
		if _, err := fmt.Fprintf(w, "%s  %s %s\n", bg, ansi.Reset, c.String()); err != nil {
			return err
		}
	}
	return nil
}

type swatchResult struct {
	Hex      string      `json:"hex"`
	Color    golor.Color `json:"color"`
	Index256 int         `json:"index256"`
}

func newPaletteCommand(stdout io.Writer, opts *cliOptions) *cobra.Command {
	var n int
	var algorithm, mode string
	cmd := &cobra.Command{
		Use:   "palette <image>",
		Short: "Extract a simple palette from an image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if n <= 0 {
				return fmt.Errorf("n must be greater than 0")
			}
			if algorithm != "frequency" && algorithm != "default" {
				return fmt.Errorf("unknown palette algorithm %q", algorithm)
			}
			colorMode, err := parseColorMode(mode)
			if err != nil {
				return err
			}
			colors, err := extractPalette(args[0], n)
			if err != nil {
				return err
			}
			if opts.json {
				return writeJSON(stdout, map[string]any{
					"algorithm": algorithm,
					"colors":    colors,
				})
			}
			return printSwatches(stdout, false, colors, colorMode)
		},
	}
	cmd.Flags().IntVarP(&n, "n", "n", 5, "number of colors to extract")
	cmd.Flags().StringVar(&algorithm, "algorithm", "frequency", "palette algorithm: frequency")
	cmd.Flags().StringVar(&mode, "color-mode", "auto", "color mode: auto, truecolor, 256")
	return cmd
}

func extractPalette(path string, n int) ([]golor.Color, error) {
	file, err := os.Open(path) // #nosec G304 -- CLI intentionally reads the user-supplied image path.
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	counts := map[uint32]int{}
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := golor.FromStdColor(img.At(x, y))
			key := uint32(c.R8())<<24 | uint32(c.G8())<<16 | uint32(c.B8())<<8 | uint32(c.A8())
			counts[key]++
		}
	}
	keys := make([]uint32, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if counts[keys[i]] == counts[keys[j]] {
			return keys[i] < keys[j]
		}
		return counts[keys[i]] > counts[keys[j]]
	})
	if n > len(keys) {
		n = len(keys)
	}
	colors := make([]golor.Color, n)
	for i := 0; i < n; i++ {
		key := keys[i]
		colors[i] = golor.RGBA(uint8(key>>24), uint8(key>>16), uint8(key>>8), uint8(key))
	}
	return colors, nil
}
