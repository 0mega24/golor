# Changelog

All notable changes to this project will be documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- v2 module path: `github.com/0mega24/golor/v2`
- Alpha-aware `Color` model with `A`, `RGBA`, `RGBAf`, `A8`, 8-digit hex parsing, and `#rrggbbaa` formatting for non-opaque colors
- Standard library image interop: `Color.RGBA`, `FromStdColor`, `ImagePixels`, `NewNRGBA`, `PaintPixels`, `Fill`, and `Color.NRGBA`
- CSS parsing and emission: `ParseCSS`, `ParseRGB`, `ParseHSL`, `Named`, `CSSRGBString`, and `CSSHSLString`
- JSON serialization for `Color` as 8-bit `{r,g,b,a}` objects
- `convert.CMYK`, `convert.ToCMYK`, and `convert.FromCMYK`
- Core `Color` type with `RGB`, `RGBf`, and `Hex` constructors
- `convert` package: HSL, HSV, L\*a\*b\*, and LCH color space conversions
- `adjust` package: lighten, darken, saturate, desaturate, hue shift, warm, cool, tint, shade, value
- `blend` package: mix, multiply, screen, overlay, hard light, soft light, difference, luminosity
- `contrast` package: WCAG 2.1 luminance, contrast ratio, and enforcement via binary bisection
- `deltae` package: ΔE76 and CIEDE2000 color difference, EnsureDistinct palette utility
- `colorblind` package: deficiency simulation (deuteranopia, protanopia, tritanopia) and AccessiblePalette
- `harmony` package: complementary, triadic, analogous, tetradic, split-complementary, and lightness-expand generators
- `gradient` package: interpolation in RGB, HSL, L\*a\*b\*, LCH, and multi-stop variants
- `palette` package: dominant image palette extraction with median-cut, k-means, and octree algorithms
- `transform` package: fluent `Chain` builder for composing transformations
- CI pipeline with `golangci-lint`, `gofumpt`, `go vet`, and race-enabled tests

### Changed
- `Color` is now a four-field straight-alpha struct; code using struct literals must set `A` explicitly or switch to constructors.
- Internal imports moved from `github.com/0mega24/golor` to `github.com/0mega24/golor/v2`.
- Existing color transformations preserve alpha unless the function explicitly interpolates or composites it.
- Gradient functions interpolate alpha linearly.
- Contrast functions resolve alpha by compositing foreground over background over white before calculating WCAG ratios.
- Delta E functions composite non-opaque colors over white before Lab comparison.
- `colorblind.AccessiblePalette` now uses a bounded adaptive hue-shift search and returns best-effort results for pathological palettes.

[Unreleased]: https://github.com/0mega24/golor/v2/compare/main...HEAD
