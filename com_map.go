package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getMap(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	mapID := r.FormValue("map-id")

	if m, ok := maps[mapID]; ok {
		w.Write(m.getStatus())
	}
}
