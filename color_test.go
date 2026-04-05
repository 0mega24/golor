package golor_test

import (
	"testing"

	"github.com/0mega24/golor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRGB(t *testing.T) {
	c := golor.RGB(255, 0, 128)
	assert.Equal(t, uint8(255), c.R8())
	assert.Equal(t, uint8(0), c.G8())
	assert.Equal(t, uint8(128), c.B8())
}

func TestRGBf(t *testing.T) {
	c := golor.RGBf(1.5, -0.1, 0.5)
	assert.Equal(t, 1.0, c.R)
	assert.Equal(t, 0.0, c.G)
	assert.Equal(t, 0.5, c.B)
}

func TestHex(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
		wantStr string
	}{
		{"#ff0080", false, "#ff0080"},
		{"ff0080", false, "#ff0080"},
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
