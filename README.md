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
| `convert` | HSL, HSV, CMYK, XYZ, L\*a\*b\*, LCH, OKLab, OKLCH, Display P3; gamut checks and CSS Color 4 gamut mapping |
| `adjust` | Lighten, darken, saturate, hue shift, warm/cool, tint, shade |
| `blend` | Blend modes: multiply, screen, overlay, soft light, and more |
| `contrast` | WCAG 2.1 contrast ratio, APCA, and enforcement |
| `deltae` | Perceptual color difference (dE76, CIEDE2000, OKLab) |
| `colorblind` | Deficiency simulation and accessible palette generation |
| `harmony` | Complementary, triadic, analogous, tetradic, split schemes |
| `gradient` | RGB/linear-light RGB, HSL/HSV, LAB/LCH, OKLab/OKLCH; hue routes, easing, chroma controls, and positioned stops |
| `palette` | Dominant color extraction from images |
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

`ParseCSS` accepts CSS hex forms, named colors, `transparent`, `rgb()`/`rgba()`, `hsl()`/`hsla()`, `hwb()`, `lab()`, `lch()`, `oklab()`, and `oklch()`. It also accepts `color()` with `srgb`, `srgb-linear`, `display-p3`, `xyz`, `xyz-d65`, or `xyz-d50`; context-dependent colors return errors.
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

// OKLCH variants produce similar visual changes across hues.
// Results are gamut-mapped instead of clipped.
lighter  = adjust.LightenOKLCH(c, 0.1)
vibrant  = adjust.SaturateOKLCH(c, 0.1)
shifted  = adjust.ShiftHueOKLCH(c, 30)
```

### Fluent chaining

```go
import "github.com/0mega24/golor/v2/transform"

result := transform.From(c).
    Lighten(0.1).
    Saturate(0.2).
    Warm(0.15).
    Color()

// Chain methods modify the chain in place; use Clone to branch.
base  := transform.From(c)
light := base.Clone().LightenOKLCH(0.2).Color()
dark  := base.Clone().DarkenOKLCH(0.2).Color()
```

### Color space conversions

```go
import "github.com/0mega24/golor/v2/convert"

hsl   := convert.ToHSL(c)   // HSL{H:16, S:1.0, L:0.6}
hsv   := convert.ToHSV(c)
lab   := convert.ToLAB(c)     // CIE L*a*b*, D65.
lch   := convert.ToLCH(c)
oklch := convert.ToOKLCH(c)   // Perceptual, as used by CSS oklch().
p3    := convert.ToP3(c)      // Display P3.
cmyk  := convert.ToCMYK(c)

back := convert.FromLCH(lch) // round-trips cleanly
back  = convert.FromCMYK(cmyk)

// Check colors outside sRGB, then map by reducing chroma (CSS Color 4)
// to preserve hue and lightness instead of clipping each channel.
vivid := convert.OKLCH{L: 0.7, C: 0.35, H: 150}
vivid.InGamut()                      // Returns false.
c = convert.FromOKLCHMapped(vivid)   // Nearest in-gamut color at the same L and hue.
```

Every 8-bit sRGB color round-trips losslessly through every space, as checked exhaustively by `make test-full`. Arbitrary float inputs round-trip within 1e-12 (OKLab within 1e-9).

### Blend modes

```go
import "github.com/0mega24/golor/v2/blend"

out := blend.Multiply(base, layer)
out  = blend.Screen(base, layer)
out  = blend.Overlay(base, layer)
out  = blend.SoftLight(base, layer)
out  = blend.Mix(base, layer, 0.5) // Premultiplied linear interpolation.
```

Blend modes follow the W3C Compositing and Blending specification used by CSS `mix-blend-mode`, including backdrop alpha and the non-separable Hue, Saturation, Color, and Luminosity modes.

### Contrast and accessibility

```go
import "github.com/0mega24/golor/v2/contrast"

r := contrast.Ratio(fg, bg)         // e.g. 4.87
ok := contrast.Meets(fg, bg, 4.5)   // WCAG AA

// Adjust fg toward black or white, whichever is closer, to meet
// the threshold against bg; ok is false if neither direction can reach it.
fg, ok = contrast.EnforceContrastOK(fg, bg, 4.5)

// APCA lightness contrast (Lc), as proposed for WCAG 3.
lc := contrast.APCA(text, bg) // For example, 63.1 for #888 on #fff.
```

### Perceptual color difference

```go
import "github.com/0mega24/golor/v2/deltae"

d76   := deltae.DeltaE76(c1, c2)    // fast, CIE 1976
d2000 := deltae.DeltaE2000(c1, c2)  // accurate, CIEDE2000
dok   := deltae.DeltaEOK(c1, c2)    // OKLab distance; ~0.02 is just noticeable.

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

`AccessiblePalette` uses a bounded adaptive hue-shift search. It returns the best effort for hard palettes rather than hanging or returning an error.

### Palette extraction

```go
import "github.com/0mega24/golor/v2/palette"

colors := palette.Extract(img, 6) // default: median-cut
colors  = palette.Extract(img, 6, palette.WithAlgorithm(palette.Octree))
colors  = palette.Extract(
    img,
    6,
    palette.WithAlgorithm(palette.KMeans),
    palette.WithSeed(42),
)
```

Results are ordered by dominance. Median-cut and octree are fast, deterministic quantizers. K-means clusters in perceptual OKLab with deterministic k-means++ initialization and an explicit seed; its output does not depend on GOMAXPROCS. Fully transparent pixels are ignored; use `palette.WithMinAlpha` to change the threshold. A 512×512 image takes roughly 15-45 ms with any algorithm.

### Color harmonies

```go
import "github.com/0mega24/golor/v2/harmony"

pair    := harmony.Complementary(c)        // [2]Color
triad   := harmony.Triadic(c)             // [3]Color
quad    := harmony.Tetradic(c)            // [4]Color
analog  := harmony.Analogous(c)           // [3]Color
split   := harmony.SplitComplementary(c)  // [3]Color
shades  := harmony.Expand(c, 5)           // []Color, lightness spread

// OKLCH variants keep perceived lightness and chroma constant across the scheme.
triad    = harmony.TriadicOKLCH(c)
shades   = harmony.ExpandOKLCH(c, 5)
```

### Gradients

```go
import "github.com/0mega24/golor/v2/gradient"

steps := gradient.RGB(black, white, 10)   // []Color, RGB interpolation
steps  = gradient.LinearRGB(black, white, 10) // Linear-light interpolation.
steps  = gradient.HSV(red, blue, 10)      // Hue/saturation/value.
steps  = gradient.HSL(red, blue, 10)      // hue takes shortest path
steps  = gradient.OKLab(red, blue, 10)    // Perceptually even (CSS default).
steps  = gradient.OKLCH(red, blue, 10)    // Preserves chroma through the hue arc.
steps  = gradient.LCH(red, blue, 10)
steps  = gradient.LAB(red, blue, 10)

// Multi-stop
stops  := []golor.Color{black, red, yellow, white}
ramp   := gradient.MultiStopOKLCH(stops, 64)
```

Gradients use premultiplied alpha, so fading to transparent does not darken the color. A grey endpoint takes the other endpoint's hue, and out-of-gamut steps are gamut-mapped. Existing functions use linear progression, the shorter hue route, and equally spaced stops; `MultiStopHSV` and `MultiStopLinearRGB` also follow these defaults.

Use `Between` or `Generate` for configurable interpolation:

```go
// A longer hue route with eased progression.
steps, err := gradient.Between(red, blue, 64, gradient.Options{
    Space:  gradient.SpaceOKLCH,
    Hue:    gradient.HueLonger,
    Easing: gradient.EaseInOut,
})
if err != nil {
    return err
}

// An explicit middle color at 30% of the ramp.
ramp, err := gradient.Generate([]gradient.Stop{
    {Color: red, Position: 0},
    {Color: yellow, Position: 0.3},
    {Color: blue, Position: 1},
}, 64, gradient.Options{
    Space:     gradient.SpaceOKLCH,
    Easing:    gradient.Logarithmic(3),
    MinChroma: 0.08,
    Chroma: func(t, c float64) float64 {
        return c + 0.04*4*t*(1-t) // Boost chroma near each segment's middle.
    },
})
if err != nil {
    return err
}
```

`Options.Space` supports all eight spaces (`SpaceRGB` is the zero value). Hue routes are `HueShorter`, `HueLonger`, `HueIncreasing`, and `HueDecreasing`; they affect polar spaces only. Equal hues with `HueLonger` make a full increasing turn. Easing choices include `Linear`, `EaseIn`, `EaseOut`, `EaseInOut`, `Logarithmic(strength)`, or a custom `Curve`. Positive logarithmic strength advances faster near the start; strength in (-1, 0) advances slower. Zero or invalid strength uses linear progression. Easing applies within each segment, keeping stop positions fixed.

Stop positions must be nondecreasing and within [0, 1]. Duplicate positions create hard transitions (the last stop wins at that position); outside the stops' span, the nearest color is held. `Generate` returns exactly `n` samples for nonempty stops and positive `n`, even for a single stop. It samples position 0 when `n == 1` and returns nil for no stops or `n <= 0`. Both configurable functions return errors for invalid options or positions.

`Chroma` and `MinChroma` apply only to LCH/OKLCH between stops, before gamut mapping, using the selected space's chroma units. A callback receives eased segment progress and alpha-aware interpolated chroma. The minimum applies after the callback. Exact stop colors are preserved, but gamut mapping can lower final chroma below the requested minimum. Easing changes progression along a path; use hue routes, chroma controls, or deliberate stops to change the path itself.

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
golor convert "oklch(70% 0.15 150)" --to p3
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
golor palette wallpaper.png -n 6 --algorithm kmeans
golor palette wallpaper.png -n 6 --algorithm octree
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
make test-full # exhaustive 16.7M-color conversion and CSS round-trip sweeps
make bench     # run all benchmarks
make fuzz      # fuzz the CSS and hex parsers (FUZZTIME=30s)
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
