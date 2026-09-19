package blend

import "github.com/0mega24/golor/v2"

// VividLight applies the Vivid Light blend mode with Porter-Duff source-over alpha compositing.
func VividLight(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		vividLightChannel(base.R, layer.R),
		vividLightChannel(base.G, layer.G),
		vividLightChannel(base.B, layer.B),
	))
}

func vividLightChannel(base, layer float64) float64 {
	if layer <= 0.5 {
		return colorBurnChannel(base, 2*layer)
	}
	return colorDodgeChannel(base, 2*(layer-0.5))
}
