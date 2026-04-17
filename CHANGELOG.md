# Changelog

All notable changes to this project will be documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- Core `Color` type with `RGB`, `RGBf`, and `Hex` constructors
- `convert` package: HSL, HSV, L\*a\*b\*, and LCH color space conversions
- `adjust` package: lighten, darken, saturate, desaturate, hue shift, warm, cool, tint, shade, value
- `blend` package: mix, multiply, screen, overlay, hard light, soft light, difference, luminosity
- `contrast` package: WCAG 2.1 luminance, contrast ratio, and enforcement via binary bisection
- `deltae` package: ΔE76 and CIEDE2000 color difference, EnsureDistinct palette utility
- `colorblind` package: deficiency simulation (deuteranopia, protanopia, tritanopia) and AccessiblePalette
- `harmony` package: complementary, triadic, analogous, tetradic, split-complementary, and lightness-expand generators
- `gradient` package: interpolation in RGB, HSL, L\*a\*b\*, LCH, and multi-stop variants
- `transform` package: fluent `Chain` builder for composing transformations
- CI pipeline with `golangci-lint`, `gofumpt`, `go vet`, and race-enabled tests

[Unreleased]: https://github.com/0mega24/golor/compare/main...HEAD
