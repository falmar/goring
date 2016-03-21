package main

import (
	htemplate "html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HomeHandler controller
func HomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tpl, err := htemplate.ParseFiles("tpl/home.gohtml")
	if err != nil {
		return
	}

	err = tpl.Execute(w, nil)

	if err != nil {
		log.Println(err)
	}
}
