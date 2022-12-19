package gocolor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRGB(t *testing.T) {
	tcs := []struct {
		r uint32
		g uint32
		b uint32
	}{
		{255, 255, 255},
		{128, 128, 128},
		{0, 0, 0},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%d,%d,%d", tc.r, tc.g, tc.b), func(t *testing.T) {
			c := RGB(tc.r, tc.g, tc.b)
			assert.Equal(t, tc.r, c.R)
			assert.Equal(t, tc.g, c.G)
			assert.Equal(t, tc.b, c.B)
		})
	}
}

func TestRGBA(t *testing.T) {
	tcs := []struct {
		r uint32
		g uint32
		b uint32
		a uint32
	}{
		{255, 255, 255, 255},
		{128, 128, 128, 128},
		{0, 0, 0, 0},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%d,%d,%d,%d", tc.r, tc.g, tc.b, tc.a), func(t *testing.T) {
			c := RGBA(tc.r, tc.g, tc.b, tc.a)
			assert.Equal(t, tc.r, c.R)
			assert.Equal(t, tc.g, c.G)
			assert.Equal(t, tc.b, c.B)
			assert.Equal(t, tc.a, c.A)
		})
	}
}

func TestRGBHex(t *testing.T) {
	tcs := []struct {
		hex       string
		expect    Color
		expectErr bool
	}{
		{"#ffffff", RGB(255, 255, 255), false},
		{"#fff", RGB(255, 255, 255), false},
		{"#000000", RGB(0, 0, 0), false},
		{"#ff000000", RGBA(255, 0, 0, 0), false},
		{"ff000000", RGBA(255, 0, 0, 0), false},
		{"", Color{}, true},
		{"#00", Color{}, true},
		{"#dfdsfadsfdsa", Color{}, true},
	}
	for _, tc := range tcs {
		t.Run(tc.hex, func(t *testing.T) {
			c, err := Hex(tc.hex)
			assert.Equal(t, tc.expectErr, err != nil)
			assert.Equal(t, tc.expect, c)
		})
	}
}

func TestRGBAhex(t *testing.T) {
	tcs := []struct {
		name      string
		hex       string
		expect    Color
		expectErr bool
	}{
		{"white", "#ffffff", RGB(255, 255, 255), true},
		{"white with alpha", "#ffffff80", RGBA(255, 255, 255, 128), false},
		{"invalid color: len", "#ffffff80ff", Color{}, true},
		{"invalid color: chars", "#ffffghff", Color{}, true},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			c, err := ahex(tc.hex)
			if err != nil {
				t.Log(err)
				assert.True(t, tc.expectErr)
				return
			}
			assert.False(t, tc.expectErr)
			assert.Equal(t, tc.expect, c)
		})
	}
}
func TestColorHex(t *testing.T) {
	testCases := []struct {
		name string
		c    Color
		want string
	}{
		{"white", RGB(255, 255, 255), "#ffffff"},
		{"black", RGB(0, 0, 0), "#000000"},
		{"red", Color{R: 255, G: 0, B: 0, A: 255}, "#ff0000"},
		{"with alpha", Color{R: 255, G: 255, B: 255, A: 128}, "#ffffff80"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.c.Hex())
		})
	}
}

func TestColorSubtract(t *testing.T) {
	tcs := []struct {
		name      string
		c1        Color
		c2        Color
		exp       Color
		shouldErr bool
	}{
		{"white minus white", RGB(255, 255, 255), RGB(255, 255, 255), Color{}, false},
		{"red minus blue", RGB(255, 0, 0), RGB(0, 0, 255), Color{}, true},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			c, err := tc.c1.Subtract(tc.c2)
			if err != nil {
				assert.True(t, tc.shouldErr)
			} else {
				assert.False(t, tc.shouldErr)
				assert.Equal(t, tc.exp, c)
			}
		})
	}
}

func TestColorTint(t *testing.T) {
	color := RGB(255, 0, 0)
	nc, err := color.Tint(1)
	assert.NoError(t, err)
	assert.Equal(t, White, nc)
}
func TestColorShade(t *testing.T) {
	color := RGB(255, 0, 0)
	nc, err := color.Shade(1)
	assert.NoError(t, err)
	assert.Equal(t, Black, nc)
	assert.NotEqual(t, color, nc)
}

func TestNorm(t *testing.T) {
	tcs := []struct {
		hex      string
		expected string
	}{
		{"fff", "ffffff"},
		{"ffffff", "ffffff"},
		{"#fff", "#ffffff"},
		{"#ffffff", "#ffffff"},
		{"fffa", "ffffffaa"},
		{"#fff0", "#ffffff00"},
	}
	for _, tc := range tcs {
		t.Run(tc.hex, func(t *testing.T) {
			assert.Equal(t, tc.expected, norm(tc.hex))
		})
	}

}
