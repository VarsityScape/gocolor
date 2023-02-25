// Copyright © 2022 VarsityScape. All rights reserved.
//
// Use of this source code is governed by a MIT-style

package gocolor

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

var (
	White = RGB(255, 255, 255)
	Black = RGB(0, 0, 0)
)

var (
	// #%02x%02x%02x
	ptn3 = addh(rptn(3))
	// #%02x%02x%02x%02x
	ptn4 = addh(rptn(4))
)

var (
	ErrInvalid = func(c Color) error { return fmt.Errorf("invalid color: %v", c) }
)

func RGB(r, g, b uint32) Color {
	return RGBA(r, g, b, 255)
}

func RGBA(r, g, b, a uint32) Color {
	return Color{R: r, G: g, B: b, A: a}
}

func Hex(hex string) (Color, error) {
	hex = addh(norm(hex))
	if len(hex) > 7 {
		return ahex(hex)
	}
	var r, g, b uint32
	_, err := fmt.Sscanf(hex, ptn3, &r, &g, &b)
	if err != nil {
		return Color{}, err
	}
	return Color{R: r, G: g, B: b, A: 255}, nil
}

type Color struct {
	R, G, B, A uint32
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, c.A
}

// Hex returns the hex representation of the color
// Note that it always begins with a hash
// Hex does not return the alpha channel if Color.A != 255
// unless the first value of noAlpha is true
func (c Color) Hex(noAlpha ...bool) string {
	if c.A != 255 && (len(noAlpha) == 0 || !noAlpha[0]) {
		return c.HexA()
	}
	return fmt.Sprintf(ptn3, c.R, c.G, c.B)
}

func (c Color) HexA() string {
	return fmt.Sprintf(ptn4, c.R, c.G, c.B, c.A)
}

func (c Color) Subtract(c2 Color) (Color, error) {
	if !c.IsValid() {
		return Color{}, ErrInvalid(c)
	}
	if !c2.IsValid() {
		return Color{}, ErrInvalid(c2)
	}
	res := Color{
		R: c.R - c2.R,
		G: c.G - c2.G,
		B: c.B - c2.B,
		A: c.A - c2.A,
	}
	if !res.IsValid() {
		return Color{}, ErrInvalid(res)
	}
	return res, nil
}

func (c Color) Tint(percent float64) (Color, error) {
	if percent < 0 || percent > 1 {
		return Color{}, fmt.Errorf("percent must be between 0 and 1")
	}
	if !c.IsValid() {
		return Color{}, ErrInvalid(c)
	}
	r, g, b, a := diff(White, c)
	newc := c
	newc.R += bound(r * percent)
	newc.G += bound(g * percent)
	newc.B += bound(b * percent)
	newc.A += bound(a * percent)
	return newc, nil
}

// Tints return n tints of the color spread evenly
func (c Color) Tints(n int) ([]Color, error) {
	sh := make([]Color, n)
	var err error
	for i := 0; i < n; i++ {
		sh[i], err = c.Tint(float64(i) / float64(n))
		if err != nil {
			return nil, err
		}
	}
	return sh, nil
}

func (c Color) Shade(percent float64) (Color, error) {
	if percent < 0 || percent > 1 {
		return Color{}, fmt.Errorf("percent must be between 0 and 1")
	}
	if !c.IsValid() {
		return Color{}, ErrInvalid(c)
	}
	r, g, b, a := diff(c, Black)
	newc := c
	newc.R -= bound(r * percent)
	newc.G -= bound(g * percent)
	newc.B -= bound(b * percent)
	newc.A -= bound(a * percent)
	return newc, nil
}

func (c Color) Shades(n int) ([]Color, error) {
	sh := make([]Color, n)
	var err error
	for i := 0; i < n; i++ {
		sh[i], err = c.Shade(float64(i) / float64(n))
		if err != nil {
			return nil, err
		}
	}
	return sh, nil
}

func (c Color) IsValid() bool {
	return c.R <= 255 && c.G <= 255 && c.B <= 255 && c.A <= 255
}

// norm normalizes a 3/4 digit hex color to a 6/8 digit hex color
func norm(s string) (result string) {
	s, didrm := Rmh(s)
	defer func() {
		if didrm {
			result = addh(result)
		}
	}()
	if len(s) == 3 || len(s) == 4 {
		var bf bytes.Buffer
		for _, c := range s {
			bf.WriteRune(c)
			bf.WriteRune(c)
		}
		return bf.String()
	} else {
		return s
	}
}

// Rmh removes the hash from a hex color
// it returns the string and a bool indicating if there was a hash
func Rmh(s string) (string, bool) {
	if len(s) > 0 && s[0] == '#' {
		return s[1:], true
	}
	return s, false
}

// rptn stands for repeat pattern
// it returns the format "%02x" repeated n times
func rptn(n int) string {
	return strings.Repeat("%02x", n)
}

// addh adds the hash to a hex color
// it returns the string if it already has a hash
func addh(s string) string {
	if len(s) > 1 && s[0] == '#' {
		return s
	}
	return "#" + s
}

// ahex stands for alpha hex
// it converts hex string that includes an alpha channel
// Example: ahex(#ffffff00) -> Color{R: 255, G: 255, B: 255, A: 0}
// it returns an error if the hex string is invalid or the length is invalid
func ahex(hex string) (Color, error) {
	hex = addh(norm(hex))
	if len(hex) != 9 {
		return Color{}, fmt.Errorf("invalid ahex color length: %s: %d", hex, len(hex))
	}
	var r, g, b, a uint32
	_, err := fmt.Sscanf(hex, ptn4, &r, &g, &b, &a)
	if err != nil {
		return Color{}, err
	}
	return Color{R: r, G: g, B: b, A: a}, nil
}

func diff(c1, c2 Color) (r, g, b, a float64) {
	return float64(c1.R) - float64(c2.R),
		float64(c1.G) - float64(c2.G),
		float64(c1.B) - float64(c2.B),
		float64(c1.A) - float64(c2.A)
}

func bound(v float64) uint32 {
	return uint32(math.Max(0, math.Min(255, math.Round(v))))
}
