# golor

A Go library for color manipulation, conversion, and accessibility.

```
go get github.com/0mega24/golor/v2
```

---

## Overview

golor provides a single `Color` type (normalized sRGB) and a set of focused sub-packages for everything you might need to do with colors in Go.

| Package | What it does |
|---|---|
| `golor` | Core type, constructors, hex/CSS parsing, JSON, image/color interop |
| `convert` | HSL, HSV, L\*a\*b\*, LCH, CMYK conversions |
| `adjust` | Lighten, darken, saturate, hue shift, warm/cool, tint, shade |
| `blend` | Blend modes: multiply, screen, overlay, soft light, and more |
| `contrast` | WCAG 2.1 contrast ratio and enforcement |
| `deltae` | Perceptual color difference (dE76 and CIEDE2000) |
| `colorblind` | Deficiency simulation and accessible palette generation |
| `harmony` | Complementary, triadic, analogous, tetradic, split schemes |
| `gradient` | Interpolation in RGB, HSL, L\*a\*b\*, and LCH |
| `transform` | Fluent chained transformations |

---

## Usage

### Creating colors

```go
import "github.com/0mega24/golor/v2"

c := golor.RGB(255, 107, 53)       // from 8-bit channels
c := golor.RGBA(255, 107, 53, 128) // with alpha
c := golor.RGBf(1.0, 0.42, 0.21)  // from normalized floats
c, err := golor.Hex("#ff6b35")     // from hex string
c, err := golor.ParseCSS("hsl(9 100% 60% / 80%)")
c, err := golor.Named("tomato")

fmt.Println(c) // #ff6b35
```

Opaque colors print as `#rrggbb`; non-opaque colors print as `#rrggbbaa`.
`Color` marshals to JSON as readable 8-bit channels, e.g. `{"r":255,"g":107,"b":53,"a":255}`.

### Image interop

```go
img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
img.Set(0, 0, golor.RGBA(255, 0, 0, 128)) // Color satisfies image/color.Color

pixels := golor.ImagePixels(img)
out := golor.NewNRGBA(2, 1, pixels)
_ = out
```

### Adjusting colors

```go
import "github.com/0mega24/golor/v2/adjust"

lighter  := adjust.Lighten(c, 0.1)
darker   := adjust.Darken(c, 0.1)
vibrant  := adjust.Saturate(c, 0.2)
muted    := adjust.Desaturate(c, 0.2)
shifted  := adjust.ShiftHue(c, 30)
warmer   := adjust.Warm(c, 0.3)
cooler   := adjust.Cool(c, 0.3)
tinted   := adjust.Tint(c, 0.2)   // toward white
shaded   := adjust.Shade(c, 0.2)  // toward black
```

### Fluent chaining

```go
import "github.com/0mega24/golor/v2/transform"

result := transform.From(c).
    Lighten(0.1).
    Saturate(0.2).
    Warm(0.15).
    Color()
```

### Color space conversions

```go
import "github.com/0mega24/golor/v2/convert"

hsl := convert.ToHSL(c)   // HSL{H:16, S:1.0, L:0.6}
hsv := convert.ToHSV(c)
lab := convert.ToLAB(c)
lch := convert.ToLCH(c)
cmyk := convert.ToCMYK(c)

back := convert.FromLCH(lch) // round-trips cleanly
back  = convert.FromCMYK(cmyk)
```

### Blend modes

```go
import "github.com/0mega24/golor/v2/blend"

out := blend.Multiply(base, layer)
out  = blend.Screen(base, layer)
out  = blend.Overlay(base, layer)
out  = blend.SoftLight(base, layer)
out  = blend.Mix(base, layer, 0.5) // linear interpolation
```

### Contrast and accessibility

```go
import "github.com/0mega24/golor/v2/contrast"

r := contrast.Ratio(fg, bg)         // e.g. 4.87
ok := contrast.Meets(fg, bg, 4.5)   // WCAG AA

// Adjust fg until it meets the threshold against bg
fg = contrast.EnforceContrast(fg, bg, 4.5)
```

### Perceptual color difference

```go
import "github.com/0mega24/golor/v2/deltae"

d76   := deltae.DeltaE76(c1, c2)    // fast, CIE 1976
d2000 := deltae.DeltaE2000(c1, c2)  // accurate, CIEDE2000

// Nudge palette entries apart until each pair differs by at least minDeltaE
palette = deltae.EnsureDistinct(palette, 10.0)
```

### Colorblind simulation

```go
import "github.com/0mega24/golor/v2/colorblind"

sim := colorblind.Simulate(c, colorblind.Deuteranopia)
sim  = colorblind.Simulate(c, colorblind.Protanopia)
sim  = colorblind.Simulate(c, colorblind.Tritanopia)

// Shift palette hues until distinguishable under the given deficiency
safe := colorblind.AccessiblePalette(palette, colorblind.Deuteranopia)
```

### Color harmonies

```go
import "github.com/0mega24/golor/v2/harmony"

pair    := harmony.Complementary(c)        // [2]Color
triad   := harmony.Triadic(c)             // [3]Color
quad    := harmony.Tetradic(c)            // [4]Color
analog  := harmony.Analogous(c)           // [3]Color
split   := harmony.SplitComplementary(c)  // [3]Color
shades  := harmony.Expand(c, 5)           // []Color, lightness spread
```

### Gradients

```go
import "github.com/0mega24/golor/v2/gradient"

steps := gradient.RGB(black, white, 10)   // []Color, RGB interpolation
steps  = gradient.HSL(red, blue, 10)      // hue takes shortest path
steps  = gradient.LCH(red, blue, 10)      // perceptually uniform
steps  = gradient.LAB(red, blue, 10)

// Multi-stop
stops  := []golor.Color{black, red, yellow, white}
ramp   := gradient.MultiStopLCH(stops, 64)
```

---

## CLI

Install the command:

```sh
go install github.com/0mega24/golor/v2/cmd/golor@latest
```

Or build it from this repo:

```sh
make build cli
./bin/golor --help
```

Convert colors:

```sh
golor convert "#ff6b35" --to hsl
golor convert "hsl(9 100% 60% / 80%)" --from css --to cmyk --json
```

Check contrast:

```sh
golor contrast "#ffffff" "#000000"
golor contrast "#ffffff" "#777777" --json
golor contrast "#777777" "#ffffff" --enforce 4.5
```

Preview terminal swatches:

```sh
golor preview "#ff0000" "tomato" "rgb(0 0 255 / 80%)"
golor preview "#ff0000" --color-mode 256
```

Extract a simple image palette:

```sh
golor palette wallpaper.png -n 6
golor palette wallpaper.png -n 6 --json
```

The `ansi` package can also be imported directly by terminal applications:

```go
import "github.com/0mega24/golor/v2/ansi"

swatch := ansi.Background(golor.RGB(255, 0, 0)) + "  " + ansi.Reset
```

---

## Development

```sh
make test      # run tests with race detector
make lint      # run golangci-lint
make fmt       # format with gofumpt
make fix       # auto-fix lint issues where possible
make vet       # run go vet
make build cli # build ./bin/golor
```

Requires [golangci-lint](https://golangci-lint.run) and [gofumpt](https://github.com/mvdan/gofumpt) to be installed for `lint` and `fmt`.

See [CONTRIBUTING.md](CONTRIBUTING.md) for full contribution guidelines.

---

## License

Apache 2.0, see [LICENSE](LICENSE).
