package golor_test

import (
	"encoding/json"
	"image"
	stdcolor "image/color"
	"testing"

	"github.com/0mega24/golor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRGB(t *testing.T) {
	c := golor.RGB(255, 0, 128)
	assert.Equal(t, uint8(255), c.R8())
	assert.Equal(t, uint8(0), c.G8())
	assert.Equal(t, uint8(128), c.B8())
	assert.Equal(t, uint8(255), c.A8())
}

func TestRGBf(t *testing.T) {
	c := golor.RGBf(1.5, -0.1, 0.5)
	assert.Equal(t, 1.0, c.R)
	assert.Equal(t, 0.0, c.G)
	assert.Equal(t, 0.5, c.B)
	assert.Equal(t, 1.0, c.A)
}

func TestRGBA(t *testing.T) {
	c := golor.RGBA(255, 0, 128, 64)
	assert.Equal(t, uint8(64), c.A8())
	assert.Equal(t, "#ff008040", c.String())
}

func TestRGBAfClamps(t *testing.T) {
	c := golor.RGBAf(1.5, -0.1, 0.5, 2)
	assert.Equal(t, 1.0, c.R)
	assert.Equal(t, 0.0, c.G)
	assert.Equal(t, 0.5, c.B)
	assert.Equal(t, 1.0, c.A)
}

func TestHex(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
		wantStr string
	}{
		{"#ff0080", false, "#ff0080"},
		{"ff0080", false, "#ff0080"},
		{"#ff008080", false, "#ff008080"},
		{"ff008000", false, "#ff008000"},
		{"#gg0000", true, ""},
		{"short", true, ""},
	}
	for _, tc := range cases {
		c, err := golor.Hex(tc.input)
		if tc.wantErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tc.wantStr, c.String())
		}
	}
}

func TestColorString(t *testing.T) {
	assert.Equal(t, "#ffffff", golor.RGB(255, 255, 255).String())
	assert.Equal(t, "#000000", golor.RGB(0, 0, 0).String())
	assert.Equal(t, "#ff000080", golor.RGBAf(1, 0, 0, 0.5).String())
}

func TestStdColorInterop(t *testing.T) {
	c := golor.RGBA(255, 128, 0, 128)
	r, g, b, a := c.RGBA()
	assert.Equal(t, uint32(32896), a)
	assert.Equal(t, a, r)
	assert.Equal(t, uint32(16513), g)
	assert.Equal(t, uint32(0), b)

	back := golor.FromStdColor(stdcolor.NRGBA{R: 255, G: 128, B: 0, A: 128})
	assert.Equal(t, uint8(255), back.R8())
	assert.Equal(t, uint8(128), back.G8())
	assert.Equal(t, uint8(0), back.B8())
	assert.Equal(t, uint8(128), back.A8())
}

func TestImageHelpers(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	pixels := []golor.Color{golor.RGBA(255, 0, 0, 128), golor.RGB(0, 0, 255)}
	golor.PaintPixels(img, pixels)

	got := golor.ImagePixels(img)
	require.Len(t, got, 2)
	assert.Equal(t, uint8(128), got[0].A8())
	assert.Equal(t, uint8(255), got[1].B8())

	out := golor.NewNRGBA(2, 1, pixels)
	assert.Equal(t, uint8(128), out.NRGBAAt(0, 0).A)
}

func TestParseCSS(t *testing.T) {
	c, err := golor.ParseCSS("rgb(255 0 0 / 50%)")
	require.NoError(t, err)
	assert.Equal(t, uint8(255), c.R8())
	assert.Equal(t, uint8(0), c.G8())
	assert.Equal(t, uint8(0), c.B8())
	assert.Equal(t, uint8(128), c.A8())

	c, err = golor.ParseCSS("rgba(255, 0, 0, 0.5)")
	require.NoError(t, err)
	assert.Equal(t, uint8(128), c.A8())

	c, err = golor.ParseCSS("hsl(9 100% 60% / 80%)")
	require.NoError(t, err)
	assert.Equal(t, uint8(204), c.A8())

	c, err = golor.Named("tomato")
	require.NoError(t, err)
	assert.Equal(t, golor.RGB(255, 99, 71), c)

	_, err = golor.ParseCSS("hsl(999deg)")
	require.Error(t, err)
}

func TestNamedCSSColors(t *testing.T) {
	c, err := golor.Named("RebeccaPurple")
	require.NoError(t, err)
	assert.Equal(t, golor.RGB(102, 51, 153), c)

	c, err = golor.Named("transparent")
	require.NoError(t, err)
	assert.Equal(t, golor.RGBA(0, 0, 0, 0), c)
}

func TestCSSStringMethods(t *testing.T) {
	assert.Equal(t, "rgb(255 0 0 / 0.502)", golor.RGBA(255, 0, 0, 128).CSSRGBString())
	assert.Equal(t, "hsl(0 100% 50%)", golor.RGB(255, 0, 0).CSSHSLString())
}

func TestColorJSON(t *testing.T) {
	original := golor.Color{R: 1, G: 0.42, B: 0.21, A: 1}
	data, err := json.Marshal(original)
	require.NoError(t, err)
	assert.JSONEq(t, `{"r":255,"g":107,"b":54,"a":255}`, string(data))

	var decoded golor.Color
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, original.R8(), decoded.R8())
	assert.Equal(t, original.G8(), decoded.G8())
	assert.Equal(t, original.B8(), decoded.B8())
	assert.Equal(t, original.A8(), decoded.A8())

	require.Error(t, json.Unmarshal([]byte(`{"r":255,"g":0,"b":0}`), &decoded))
}

func ExampleRGB() {
	c := golor.RGB(255, 128, 0)
	_ = c.String() // "#ff8000"
}

func ExampleHex() {
	c, err := golor.Hex("#ff8000")
	if err != nil {
		panic(err)
	}
	_ = c.R8() // 255
}
