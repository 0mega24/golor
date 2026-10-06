# Changelog

All notable changes to this project will be documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- `gradient`: HSV and linear-light RGB (including multi-stop variants), configurable `Between`/`Generate` with positioned stops and hard transitions, shorter/longer/increasing/decreasing hue routes, easing and logarithmic curves, and LCH/OKLCH chroma callbacks and minimum chroma. Existing gradient defaults are unchanged.
- `convert`: OKLab and OKLCH (`ToOKLab`, `FromOKLab`, `ToOKLCH`, `FromOKLCH`, `DeltaEOK`), Display P3 (`ToP3`, `FromP3`), and CIE XYZ D65 (`ToXYZ`, `FromXYZ`).
- `convert`: gamut checks (`InGamut` on `LAB`, `LCH`, `OKLab`, `OKLCH`; `P3.InSRGBGamut`) and CSS Color 4 gamut mapping (`FromLABMapped`, `FromLCHMapped`, `FromOKLabMapped`, `FromOKLCHMapped`, `FromP3Mapped`), which reduce OKLCH chroma instead of clipping individual channels.
- `convert`: `IsAchromatic` on `HSL`, `HSV`, `LCH`, and `OKLCH`.
- `gradient`: `OKLab`, `OKLCH`, `MultiStopHSL`, `MultiStopLAB`, `MultiStopOKLab`, `MultiStopOKLCH`.
- `adjust`: perceptual `LightenOKLCH`, `DarkenOKLCH`, `SaturateOKLCH`, `DesaturateOKLCH`, `ShiftHueOKLCH`; corresponding `transform.Chain` methods and `Chain.Clone`.
- `harmony`: `ComplementaryOKLCH`, `TriadicOKLCH`, `AnalogousOKLCH`, `SplitComplementaryOKLCH`, `TetradicOKLCH`, `ExpandOKLCH`.
- `contrast`: `APCA` (APCA-W3 0.0.98G-4g) and `EnforceContrastOK`, which reports whether the target ratio was reached.
- `deltae`: `DeltaEOK`.
- CLI: `golor convert --to oklab|oklch|p3|xyz`; `golor contrast` reports APCA Lc and whether `--enforce` succeeded.
- Benchmarks in every package, fuzz tests for CSS and hex parsing, conversion fidelity property tests, and `make bench`, `make fuzz`, and `make test-full` (exhaustive 16.7M-color sweeps).
- `ParseCSS` supports the full CSS Color 4 syntax: `#rgb`/`#rgba` shorthand, legacy and modern `rgb()`/`hsl()` (aliases, `none`, clamping of out-of-range values), hue units (`deg`, `rad`, `grad`, `turn`) with wrapping, `hwb()`, `lab()`/`lch()` (D50, as CSS defines them), `oklab()`, `oklch()`, and `color()` with `srgb`, `srgb-linear`, `display-p3`, `xyz`, `xyz-d65`, `xyz-d50`. Colors outside sRGB are gamut-mapped.
- `Color.CSSOKLCHString`.
- `palette.WithMinAlpha` and `palette.DefaultMinAlpha`.

### Fixed
- LAB/LCH greys had a spurious 158° hue because the XYZ matrix's Y row summed to 1.0000001, causing LCH gradients from white or grey to pass through teal. sRGB↔XYZ now uses the exact CSS Color 4 matrices, and greys have exactly zero chroma.
- Off-grid LAB/LCH round trips lost up to 1.7e-6 per channel; all conversions now round-trip within 1e-12 (OKLab within 1e-9), and all 16,777,216 8-bit colors round-trip losslessly through every space.
- `FromHSL`/`FromHSV` returned incorrect colors for hues outside [0, 360) (for example `FromHSL(H=720)` was black); hue is now wrapped and S/L/V are clamped.
- Conversions accept struct-literal colors with out-of-range or NaN channels without producing out-of-range output (for example, S > 1, negative K, or luminance > 1).
- Gradients interpolate premultiplied alpha (fading to transparent no longer greys out), use the other endpoint's hue for a grey endpoint, and gamut-map LAB/LCH/OKLab/OKLCH results instead of clipping. Multi-stop gradients convert each stop once instead of per sample.
- `adjust.Saturate`/`SetSaturation` no longer tint greys red.
- `deltae.EnsureDistinct` re-checks every pair until all are distinct, moves colors by the smallest lightness/hue/chroma change (gamut-mapped), and no longer collapses identical inputs toward black.
- `ansi.Nearest256` returns only fixed palette indexes (16-255); indexes 0-15 are theme-defined. It computes the result directly instead of scanning 256 entries.
- `blend`: modes composite with W3C backdrop-alpha handling (`Multiply(transparent, red)` is red, not black); `Hue`, `Saturation`, `Color`, and `Luminosity` use the W3C non-separable formulas (as used by CSS `mix-blend-mode`) instead of HSL; `Mix` interpolates premultiplied alpha.
- `colorblind.Simulate` uses the published Viénot 1999 matrices for protanopia/deuteranopia and Brettel 1997 for tritanopia (via DaltonLens); deuteranopia previously mapped pure blue to lavender. Unknown deficiencies return the input instead of black. `AccessiblePalette` is about 75× faster for 128 colors.
- `contrast.EnforceContrast` tries both black and white; it previously chose white for any background with luminance below 0.5, failing on mid tones. It returns a color that meets the threshold after 8-bit rounding and composites translucent backgrounds before choosing a direction. It is about 6× faster.
- `CSSHSLString` output could not be parsed back (`hsl(360 …)` was emitted) and lost precision for a third of colors; `CSSRGBString`, `CSSHSLString`, `CSSOKLCHString`, and `String` now round-trip through `ParseCSS` to the same 8-bit color for all 16.7M colors and every alpha.
- `ParseCSS` rejected valid CSS (`#fff`, `hsl(360 …)`, `rgb(255, 0, 0, 0.5)`, `rgb(300 0 0)`) and accepted invalid CSS (`rgb(255,0,0,)`).
- `RGBf`/`RGBAf` and channel accessors map NaN to 0 instead of retaining NaN.
- `NearestNamed` uses perceptual OKLab distance instead of RGB distance and no longer sorts the name table on every call (about 50× faster).
- `ImagePixels`, `PaintPixels`, and `Fill` read and write `*image.NRGBA`/`*image.RGBA` pixel buffers directly (one allocation instead of one per pixel); `FromStdColor` is lossless for straight-alpha inputs.
- `Hex` is about 3.5× faster and `String` about 5× faster.
- `palette.Extract` skips fully transparent pixels by default (previously, icons reported `#00000000` as their dominant color); `WithMinAlpha(0)` restores the old behavior.
- `palette` is 10-50× faster with minimal allocations, using cached swatch keys, radix-sorted pixel counting, and direct `Pix` access for NRGBA/RGBA images. Sorting previously recomputed keys on each comparison and took 73% of the time. 512×512 images take 13-43 ms instead of 180 ms to 2.1 s.
- KMeans clusters in perceptual OKLab, bins large inputs into a 5-bit weighted-mean histogram, seeds k-means++ in O(U·k), assigns with Hamerly bounds in parallel (output is identical for any GOMAXPROCS), and stops when centroids converge. 512×512 noise at n=64 takes 43 ms instead of 2.1 s.
- Octree no longer slows down sharply when n is just above a level boundary (n=520 on 256² noise: 79 ms instead of 14.6 s) and returns exactly n colors when the image has at least n distinct colors.

### Changed
- LAB values differ in the fourth decimal place because the CSS Color 4 sRGB matrices replace the 7-digit Lindbloom matrices.
- `transform.Chain` methods mutate and return the receiver, as before. This behavior is documented, and `Clone` supports branching.

## [2.0.0]

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

[Unreleased]: https://github.com/0mega24/golor/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/0mega24/golor/releases/tag/v2.0.0
