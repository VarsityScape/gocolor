package gocolor

import "fmt"

func ExampleColor_Hex() {
	c := RGB(255, 255, 255)
	fmt.Println(c.Hex())
	// Output: #ffffff
}

func ExampleColor_RGBA() {
	c := White
	r, g, b, a := c.RGBA()
	fmt.Println(r, g, b, a)
	// Output: 255 255 255 255
}

func ExampleColor_Tints() {
	c := Black
	tints, _ := c.Tints(5)
	for _, t := range tints {
		fmt.Print(t.Hex(), " ")
	}
	// Output: #000000 #333333 #666666 #999999 #cccccc
}

func ExampleColor_Shades() {
	c := White
	shades, _ := c.Shades(5)
	for _, s := range shades {
		fmt.Print(s.Hex(), " ")
	}
	// Output: #ffffff #cccccc #999999 #666666 #333333
}
