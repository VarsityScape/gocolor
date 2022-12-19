package gocolor

import "fmt"

func ExampleColor_Hex() {
	c := RGB(255, 255, 255)
	fmt.Println(c.Hex())
	// Output: #ffffff
}
