// Copyright © 2022 VarsityScape. All rights reserved.
//
// Use of this source code is governed by a MIT-style

package web

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/varsityscape/gocolor"
)

var (
	//go:embed template/* static/*
	// f is the filesystem containing the static files
	f embed.FS
	// tmpl is the template used to render the HTML
	tmpl = template.New("index.gohtml")
)

func init() {
	funcMap := template.FuncMap{
		"percent": func(i, total int) int {
			return i * 100 / total
		},
	}
	_, err := tmpl.Funcs(funcMap).ParseFS(f, "template/index.gohtml")
	if err != nil {
		panic(err)
	}

}

type handler struct {
	dcolor string
	dcount int
}

func NewHandler(defaultColor string, defaultCount int) *handler {
	return &handler{dcolor: defaultColor, dcount: defaultCount}
}

type Data struct {
	Count int
	Color gocolor.Color
}

func (h *handler) HTMLColor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	colorStr := orDefault(vars["color"], h.dcolor)
	clr, err := gocolor.Hex(colorStr)
	if err != nil {
		http.Error(w, "Invalid color", http.StatusBadRequest)
		return
	}
	count := h.dcount
	qc := r.URL.Query().Get("count")
	if qc != "" {
		c, err := strconv.Atoi(qc)
		if err == nil {
			count = c
		}
	}

	tmpl.Execute(w, Data{Count: count, Color: clr})
}

func (h *handler) FS(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(f)).ServeHTTP(w, r)
}

func orDefault(s string, def string) string {
	if s == "" {
		return def
	}
	return s
}
