// Copyright © 2022 VarsityScape. All rights reserved.
//
// Use of this source code is governed by a MIT-style

package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/varsityscape/gocolor"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/varsityscape/gocolor/web"
)

var (
	port   = flag.String("port", "8080", "port to listen on")
	noHTML = flag.Bool("no-html", false, "disable html output")
	dCount = flag.Int("count", 10, "number of shades and tints to generate")
	dColor = flag.String("color", "fff", "color to display (without #)")
)

func init() {
	flag.Parse()
}

func main() {
	if *noHTML {
		printColor()
		return
	}
	// run a webserver on port *port with path color
	// /color should return the color

	r := mux.NewRouter()
	handler := web.NewHandler(*dColor, *dCount)
	r.HandleFunc("/", handler.HTMLColor)
	r.HandleFunc("/{color}", handler.HTMLColor)
	r.PathPrefix("/static/").HandlerFunc(handler.FS)

	log.Println("Listening on port", *port)
	log.Printf("Open http://localhost:%s/%s?count=%d to view shades and tints.\n", *port, *dColor, *dCount)

	http.ListenAndServe(":"+*port, r)
}

func printColor() {
	cl, err := gocolor.Hex(*dColor)
	if err != nil {
		log.Fatal(err)
	}
	pad := 9
	fmt.Print("Original color is ")
	color.HEXStyle("ccc", cl.Hex(true)).Printf("%-*s", pad, cl.Hex())
	fmt.Println()
	shades, err := cl.Shades(*dCount)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to generate shades"))
	}
	tints, err := cl.Tints(*dCount)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to generate tints"))
	}
	color.Bold.Println("Shades:")
	for i, c := range shades {
		fmt.Printf("%3d%%: ", percent(i, *dCount))
		color.HEXStyle("ccc", c.Hex(true)).Printf("%-*s", pad, c.Hex())
	}
	color.Bold.Println("\n\nTints:")
	for i, c := range tints {
		fmt.Printf("%3d%%: ", percent(i, *dCount))
		color.HEXStyle("ccc", c.Hex(true)).Printf("%-*s", pad, c.Hex())
	}
	fmt.Println()
}

func percent(i, total int) int {
	return int(math.Round(float64(i) / float64(total) * 100))
}
