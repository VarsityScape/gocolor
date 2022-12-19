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
		name    string
		c       Color
		noAlpha []bool
		want    string
	}{
		{name: "white", c: RGB(255, 255, 255), want: "#ffffff"},
		{name: "black", c: RGB(0, 0, 0), want: "#000000"},
		{name: "red", c: Color{R: 255, G: 0, B: 0, A: 255}, want: "#ff0000"},
		{name: "with alpha", c: Color{R: 255, G: 255, B: 255, A: 128}, want: "#ffffff80"},
		{name: "with alpha 2", c: RGBA(255, 255, 255, 128), want: "#ffffff80", noAlpha: []bool{false}},
		{name: "don't want alpha", c: Color{R: 255, G: 255, B: 255, A: 128}, want: "#ffffff", noAlpha: []bool{true}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.c.Hex(tc.noAlpha...))
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
		{"invalid color src", Color{R: 300}, RGB(0, 0, 255), Color{}, true},
		{"invalid color param", White, Color{R: 300}, Color{}, true},
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
	tcs := []struct {
		name    string
		color   Color
		percent float64
		exp     Color
		err     bool
	}{
		{"invalid percent", White, 1.1, Color{}, true},
		{"invalid color", Color{R: 300}, 0.5, Color{}, true},
		{"white", White, 1, RGB(255, 255, 255), false},
		{"red", RGB(255, 0, 0), 0.5, RGB(255, 128, 128), false},
	}
	for _, tc := range tcs {
		nc, err := tc.color.Tint(tc.percent)
		if err != nil {
			assert.True(t, tc.err)
		} else {
			assert.False(t, tc.err)
			assert.Equal(t, tc.exp, nc)
		}
	}
}
func TestColorShade(t *testing.T) {
	tcs := []struct {
		name    string
		color   Color
		percent float64
		exp     Color
		err     bool
	}{
		{"invalid percent", White, 1.1, Color{}, true},
		{"invalid color", Color{R: 300}, 0.5, Color{}, true},
		{"white", White, 1, Black, false},
		{"red", RGB(255, 0, 0), 0.5, RGB(127, 0, 0), false},
		{"hexa color", Color{R: 255, G: 0, B: 0, A: 128}, 0.5, Color{R: 127, G: 0, B: 0, A: 128}, false},
	}
	for _, tc := range tcs {
		nc, err := tc.color.Shade(tc.percent)
		if err != nil {
			assert.True(t, tc.err)
		} else {
			assert.False(t, tc.err)
			assert.Equal(t, tc.exp, nc)
		}
	}
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

func TestColor_ShadesTintsErr(t *testing.T) {
	tcs := []struct {
		name      string
		c         Color
		count     int
		shouldErr bool
	}{
		{"invalid color", Color{R: 300}, 5, true},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.c.Shades(tt.count); (err != nil && !tt.shouldErr) || (err == nil && tt.shouldErr) {
				t.Errorf("Color.Shades() error = %v, shouldErr %v", err, tt.shouldErr)
			}
			if _, err := tt.c.Tints(tt.count); (err != nil && !tt.shouldErr) || (err == nil && tt.shouldErr) {
				t.Errorf("Color.Tints() error = %v, shouldErr %v", err, tt.shouldErr)
			}

		})
	}
}
