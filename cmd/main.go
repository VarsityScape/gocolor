// Copyright © 2022 VarsityScape. All rights reserved.
//
// Use of this source code is governed by a MIT-style

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/varsityscape/gocolor/web"
)

var (
	port   = flag.String("port", "8080", "port to listen on")
	noHTML = flag.Bool("no-html", false, "disable html output")

	dCount = flag.Int("count", 10, "number of shades and tints to generate")
	dColor = flag.String("color", "#fff", "color to display")
)

func init() {
	flag.Parse()
}

func main() {
	if *noHTML {
		// print the color
		// print the shades
		// print the tints
		return
	}
	// run a webserver on port *port with path color
	// /color should return the color

	r := mux.NewRouter()
	handler := web.New(*dColor, *dCount)
	r.HandleFunc("/", handler.HTMLColor)
	r.HandleFunc("/{color}", handler.HTMLColor)
	r.PathPrefix("/static/").HandlerFunc(handler.FS)

	log.Println("Listening on port", *port)

	http.ListenAndServe(":"+*port, r)
}
