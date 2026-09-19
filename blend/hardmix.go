package blend

import "github.com/0mega24/golor/v2"

// HardMix applies the Hard Mix blend mode with Porter-Duff source-over alpha compositing.
func HardMix(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		hardMixChannel(base.R, layer.R),
		hardMixChannel(base.G, layer.G),
		hardMixChannel(base.B, layer.B),
	))
}

func hardMixChannel(base, layer float64) float64 {
	if vividLightChannel(base, layer) < 0.5 {
		return 0
	}
	return 1
}
